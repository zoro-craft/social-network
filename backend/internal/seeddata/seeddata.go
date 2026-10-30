// Package seeddata inserts the demo dataset shared by the seed command and the
// server's first-run bootstrap.
package seeddata

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	// sqliteTimeLayout is the text format go-sqlite3 reads back out of a
	// DATETIME column, and the same one the app's own inserts produce.
	sqliteTimeLayout = "2006-01-02 15:04:05"

	sessionDuration = 24 * time.Hour
)

// Counts is the number of rows inserted by one successful seed run.
type Counts struct {
	Users         int
	Sessions      int
	Follows       int
	Groups        int
	GroupMembers  int
	Invitations   int
	Requests      int
	Posts         int
	PostViewers   int
	Comments      int
	Reactions     int
	Events        int
	Responses     int
	Messages      int
	Notifications int
	Reels         int
	ReelLikes     int
	ReelViews     int
	ReelSaves     int
	ReelShares    int
	ReelComments  int
}

// DemoUser describes one seeded account. The manual seed command uses it to
// print the credentials table; the server never logs it.
type DemoUser struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Nickname  string
	IsPrivate bool
}

// FullName mirrors the users.full_name generated column.
func (u DemoUser) FullName() string { return u.FirstName + " " + u.LastName }

// Result reports what Run did. Seeded is false when the database already had
// users and wipe was false.
type Result struct {
	Seeded   bool
	Counts   Counts
	Users    []DemoUser
	Sessions map[int]string
}

// wipeOrder lists every application table in reverse dependency order so the
// foreign keys enforced by the _foreign_keys=on DSN stay satisfied.
var wipeOrder = []string{
	"reel_comments", "reel_shares", "reel_saves", "reel_views", "reel_likes", "reels",
	"story_views", "stories",
	"messages", "notifications", "event_responses", "events",
	"post_viewers", "comments", "post_reactions",
	"group_requests", "group_invitations", "group_members", "user_groups",
	"posts", "followers", "user_sessions", "users",
}

// Run seeds the demo dataset into db. When wipe is true it first deletes every
// application table. The whole operation runs in one transaction, so a failure
// leaves the database unchanged.
func Run(db *sql.DB, wipe bool) (*Result, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("seed: begin transaction: %w", err)
	}
	defer tx.Rollback()

	if wipe {
		if err := wipeAll(tx); err != nil {
			return nil, fmt.Errorf("seed: wipe: %w", err)
		}
	}

	var users int
	if err := tx.QueryRow("SELECT COUNT(*) FROM users").Scan(&users); err != nil {
		return nil, fmt.Errorf("seed: count users: %w", err)
	}
	if users > 0 {
		log.Println("Seed: demo data already present, skipping")
		return &Result{Seeded: false}, nil
	}

	if err := copyAssets(); err != nil {
		return nil, fmt.Errorf("seed: copy demo assets: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(DemoPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("seed: hash demo password: %w", err)
	}

	s := &seeder{
		tx:            tx,
		base:          time.Now().Truncate(time.Second),
		invitationIDs: map[[2]int]int{},
		requestIDs:    map[[2]int]int{},
		lastActivity:  map[int]int{},
	}

	s.seedUsers(string(hash))
	s.seedFollows()
	s.seedGroups()
	s.seedPosts()
	s.seedEvents()
	s.seedMessages()
	s.seedReels()
	s.seedLastSeen()
	s.seedNotifications()

	if s.err != nil {
		return nil, s.err
	}

	sessions, err := s.seedSessions()
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("seed: commit: %w", err)
	}

	result := &Result{
		Seeded:   true,
		Counts:   s.counts,
		Users:    demoUsers(),
		Sessions: sessions,
	}
	logSummary(result.Counts)

	return result, nil
}

func wipeAll(tx *sql.Tx) error {
	for _, table := range wipeOrder {
		if _, err := tx.Exec("DELETE FROM " + table); err != nil {
			return fmt.Errorf("clear %s: %w", table, err)
		}
	}
	// Reset the AUTOINCREMENT counters so ids restart at 1 on every run, which
	// keeps the hard-coded references in dataset.go valid.
	if _, err := tx.Exec("DELETE FROM sqlite_sequence"); err != nil {
		return fmt.Errorf("reset sequences: %w", err)
	}
	return nil
}

// seeder holds the open transaction and the clock the dataset is relative to.
// The first failed statement is sticky so the insert steps read as flat lists.
type seeder struct {
	tx   *sql.Tx
	base time.Time
	err  error

	// ids captured during insert so notifications can point at the real rows.
	groupIDs      []int          // parallel to seedGroups
	invitationIDs map[[2]int]int // {groupID, invitedUserID} -> group_invitations.id
	requestIDs    map[[2]int]int // {groupID, userID}       -> group_requests.id
	eventIDs      []int

	// lastActivity maps a user id to the fewest minutes ago they sent a message.
	lastActivity map[int]int
	counts       Counts
}

func (s *seeder) exec(what, query string, args ...any) {
	if s.err != nil {
		return
	}
	if _, err := s.tx.Exec(query, args...); err != nil {
		s.err = fmt.Errorf("seed: %s: %w", what, err)
	}
}

func (s *seeder) execID(what, query string, args ...any) int64 {
	if s.err != nil {
		return 0
	}
	res, err := s.tx.Exec(query, args...)
	if err != nil {
		s.err = fmt.Errorf("seed: %s: %w", what, err)
		return 0
	}
	id, err := res.LastInsertId()
	if err != nil {
		s.err = fmt.Errorf("seed: %s: last insert id: %w", what, err)
	}
	return id
}

// ts renders a DATETIME column the given number of hours before the base time.
func (s *seeder) ts(hoursAgo int) string {
	return s.base.Add(-time.Duration(hoursAgo) * time.Hour).Format(sqliteTimeLayout)
}

// tsMinutes renders a DATETIME column the given number of minutes before base.
func (s *seeder) tsMinutes(minutesAgo int) string {
	return s.base.Add(-time.Duration(minutesAgo) * time.Minute).Format(sqliteTimeLayout)
}

// millis renders an app-set epoch-millisecond column the given minutes early.
func (s *seeder) millis(minutesAgo int) int64 {
	return s.base.Add(-time.Duration(minutesAgo) * time.Minute).UnixMilli()
}

// seedUsers inserts seedUsers in order so their ids become 1..len(seedUsers).
func (s *seeder) seedUsers(hash string) {
	for _, u := range seedUsers {
		s.exec(
			fmt.Sprintf("insert user %s", u.Email),
			`INSERT INTO users (first_name, last_name, email, password, date_of_birth, nickname, about_me, avatar, is_private, is_public, created_at, last_seen)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?)`,
			u.FirstName, u.LastName, u.Email, hash, u.DateOfBirth, u.Nickname, u.AboutMe,
			nullableAssetURL(u.Avatar), u.IsPrivate, s.ts(u.JoinedHoursAgo), s.tsMinutes(u.LastSeenMinutesAgo),
		)
		s.counts.Users++
	}
}

func (s *seeder) seedFollows() {
	for _, f := range seedFollows {
		s.exec(
			fmt.Sprintf("insert follow %d -> %d", f.From, f.To),
			"INSERT INTO followers (follower_id, followee_id, status, created_at) VALUES (?, ?, ?, ?)",
			f.From, f.To, f.Status, s.ts(f.HoursAgo),
		)
		s.counts.Follows++
	}
}

func (s *seeder) seedGroups() {
	for i, g := range seedGroups {
		groupID := int(s.execID(
			fmt.Sprintf("insert group %q", g.Title),
			"INSERT INTO user_groups (creator_id, title, description, image) VALUES (?, ?, ?, ?)",
			g.Creator, g.Title, g.Description, assetURL(g.Image),
		))
		s.groupIDs = append(s.groupIDs, groupID)
		s.counts.Groups++
		groupNumber := i + 1

		for j, member := range g.Members {
			// Stagger joins so member ordering is not uniformly the creator first.
			joinedHoursAgo := 24 * (100 - groupNumber*10 - j)
			s.exec(
				fmt.Sprintf("insert member %d in group %d", member, groupID),
				"INSERT INTO group_members (group_id, user_id, joined_at) VALUES (?, ?, ?)",
				groupID, member, s.ts(joinedHoursAgo),
			)
			s.counts.GroupMembers++
		}

		for _, inv := range g.Invitations {
			id := int(s.execID(
				fmt.Sprintf("insert invitation for user %d in group %d", inv.Invitee, groupID),
				"INSERT INTO group_invitations (group_id, inviter_id, invited_user_id, status, created_at) VALUES (?, ?, ?, ?, ?)",
				groupID, inv.Inviter, inv.Invitee, inv.Status, s.ts(inv.HoursAgo),
			))
			s.invitationIDs[[2]int{groupID, inv.Invitee}] = id
			s.counts.Invitations++
		}

		for _, req := range g.Requests {
			id := int(s.execID(
				fmt.Sprintf("insert request from user %d for group %d", req.User, groupID),
				"INSERT INTO group_requests (group_id, user_id, status, created_at) VALUES (?, ?, ?, ?)",
				groupID, req.User, req.Status, s.ts(req.HoursAgo),
			))
			s.requestIDs[[2]int{groupID, req.User}] = id
			s.counts.Requests++
		}
	}
}

func (s *seeder) seedPosts() {
	for _, p := range seedPosts {
		var groupID any
		if p.Group > 0 {
			groupID = p.Group
		}

		postID := int(s.execID(
			fmt.Sprintf("insert post %q", p.Title),
			"INSERT INTO posts (user_id, group_id, title, content, image, privacy, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			p.Author, groupID, p.Title, p.Content, nullableAssetURL(p.Image), p.Privacy, s.ts(p.HoursAgo),
		))
		s.counts.Posts++

		for _, viewer := range p.Viewers {
			s.exec(
				fmt.Sprintf("insert viewer %d for post %d", viewer, postID),
				"INSERT INTO post_viewers (post_id, user_id) VALUES (?, ?)",
				postID, viewer,
			)
			s.counts.PostViewers++
		}
		for _, r := range p.Reactions {
			isLike := 0
			if r.Like {
				isLike = 1
			}
			s.exec(
				fmt.Sprintf("insert reaction by user %d on post %d", r.User, postID),
				"INSERT INTO post_reactions (user_id, post_id, is_like) VALUES (?, ?, ?)",
				r.User, postID, isLike,
			)
			s.counts.Reactions++
		}
		for _, c := range p.Comments {
			var commentGroup any
			if p.Group > 0 {
				commentGroup = p.Group
			}
			s.exec(
				fmt.Sprintf("insert comment by user %d on post %d", c.Author, postID),
				"INSERT INTO comments (post_id, user_id, group_id, content, image, created_at) VALUES (?, ?, ?, ?, NULL, ?)",
				postID, c.Author, commentGroup, c.Content, s.ts(c.HoursAgo),
			)
			s.counts.Comments++
		}
	}
}

func (s *seeder) seedEvents() {
	for _, e := range seedEvents {
		options := fmt.Sprintf(`["%s","%s"]`, e.Options[0], e.Options[1])
		eventDate := s.base.AddDate(0, 0, e.InDays).Format("2006-01-02")

		eventID := int(s.execID(
			fmt.Sprintf("insert event %q", e.Title),
			"INSERT INTO events (group_id, title, description, event_date, event_time, created_by, options, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			e.Group, e.Title, e.Description, eventDate, e.At, e.Creator, options, s.ts(e.CreatedHoursAgo),
		))
		s.eventIDs = append(s.eventIDs, eventID)
		s.counts.Events++

		for _, v := range e.Votes {
			s.exec(
				fmt.Sprintf("insert response from user %d for event %d", v.User, eventID),
				"INSERT INTO event_responses (event_id, user_id, response) VALUES (?, ?, ?)",
				eventID, v.User, e.Options[v.OptionIdx],
			)
			s.counts.Responses++
		}
	}
}

func (s *seeder) seedMessages() {
	for _, m := range seedMessages {
		s.exec(
			fmt.Sprintf("insert private message from user %d", m.From),
			"INSERT INTO messages (sender_id, type, ref_id, content, created_at) VALUES (?, 'private', ?, ?, ?)",
			m.From, m.To, m.Content, s.millis(m.MinutesAgo),
		)
		s.noteActivity(m.From, m.MinutesAgo)
		s.counts.Messages++
	}
	for _, m := range seedGroupMessages {
		s.exec(
			fmt.Sprintf("insert group message from user %d in group %d", m.From, m.Group),
			"INSERT INTO messages (sender_id, type, ref_id, content, created_at) VALUES (?, 'group', ?, ?, ?)",
			m.From, m.Group, m.Content, s.millis(m.MinutesAgo),
		)
		s.noteActivity(m.From, m.MinutesAgo)
		s.counts.Messages++
	}
}

func (s *seeder) seedReels() {
	for _, r := range seedReels {
		createdAt := s.millis(r.HoursAgo * 60)
		reelID := int(s.execID(
			fmt.Sprintf("insert reel by user %d", r.Author),
			"INSERT INTO reels (user_id, video, caption, created_at) VALUES (?, ?, ?, ?)",
			r.Author, assetURL(r.Video), r.Caption, createdAt,
		))
		s.counts.Reels++

		for _, userID := range r.Likes {
			s.exec(
				fmt.Sprintf("insert like by user %d on reel %d", userID, reelID),
				"INSERT INTO reel_likes (reel_id, user_id, created_at) VALUES (?, ?, ?)",
				reelID, userID, createdAt,
			)
			s.counts.ReelLikes++
		}
		for _, userID := range r.Views {
			s.exec(
				fmt.Sprintf("insert view by user %d on reel %d", userID, reelID),
				"INSERT INTO reel_views (reel_id, user_id, viewed_at) VALUES (?, ?, ?)",
				reelID, userID, createdAt,
			)
			s.counts.ReelViews++
		}
		for _, userID := range r.Saves {
			s.exec(
				fmt.Sprintf("insert save by user %d on reel %d", userID, reelID),
				"INSERT INTO reel_saves (reel_id, user_id, created_at) VALUES (?, ?, ?)",
				reelID, userID, createdAt,
			)
			s.counts.ReelSaves++
		}
		for _, userID := range r.Shares {
			s.exec(
				fmt.Sprintf("insert share by user %d on reel %d", userID, reelID),
				"INSERT INTO reel_shares (reel_id, user_id, created_at) VALUES (?, ?, ?)",
				reelID, userID, createdAt,
			)
			s.counts.ReelShares++
		}
		for _, c := range r.Comments {
			s.exec(
				fmt.Sprintf("insert comment by user %d on reel %d", c.Author, reelID),
				"INSERT INTO reel_comments (reel_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
				reelID, c.Author, c.Content, s.millis(c.HoursAgo*60),
			)
			s.counts.ReelComments++
		}
	}
}

// noteActivity keeps the freshest message timestamp seen for a user.
func (s *seeder) noteActivity(userID, minutesAgo int) {
	current, ok := s.lastActivity[userID]
	if !ok || minutesAgo < current {
		s.lastActivity[userID] = minutesAgo
	}
}

// seedLastSeen derives last_seen from the latest seeded message. Users without
// messages keep the offset from dataset.go.
func (s *seeder) seedLastSeen() {
	for i := range seedUsers {
		userID := i + 1
		minutesAgo, ok := s.lastActivity[userID]
		if !ok {
			continue
		}
		s.exec(
			fmt.Sprintf("update last_seen for user %d", userID),
			"UPDATE users SET last_seen = ? WHERE id = ?",
			s.tsMinutes(minutesAgo), userID,
		)
	}
}

// seedNotifications rebuilds the notifications the real handlers would have
// written for the rows above, reusing their exact message formats.
func (s *seeder) seedNotifications() {
	// A pending follow notifies the followee, with ref_id pointing at the requester.
	for _, f := range seedFollows {
		if f.Status != statusPending {
			continue
		}
		s.notify(f.To, f.From, notifFollowRequest, f.From,
			"New follow request from "+fullName(f.From), false, 60*f.HoursAgo)
	}

	// A newly accepted follow notifies the followee that someone followed them.
	for _, a := range acceptedFollowAlerts {
		s.notify(a.To, a.From, notifNewFollower, a.From,
			fullName(a.From)+" started following you", true, a.MinutesAgo)
	}

	for i, g := range seedGroups {
		groupID := s.groupIDs[i]
		for _, inv := range g.Invitations {
			invitationID := s.invitationIDs[[2]int{groupID, inv.Invitee}]
			s.notify(inv.Invitee, inv.Inviter, notifGroupInvite, invitationID,
				fullName(inv.Inviter)+" invited you to join "+g.Title, inv.Status != "pending", 60*inv.HoursAgo)
		}
		for _, req := range g.Requests {
			requestID := s.requestIDs[[2]int{groupID, req.User}]
			s.notify(g.Creator, req.User, notifJoinRequest, requestID,
				"New join request for your group "+fullName(req.User), req.Status != "pending", 60*req.HoursAgo)
		}
	}

	// Creating an event fans a notification out to every member but the creator.
	for i, e := range seedEvents {
		message := "New event by " + fullName(e.Creator) + " in your group: " + e.Title
		for _, member := range seedGroups[e.Group-1].Members {
			if member == e.Creator {
				continue
			}
			s.notify(member, e.Creator, notifEventCreated, s.eventIDs[i], message, false, 60*e.CreatedHoursAgo)
		}
	}
}

func (s *seeder) notify(userID, actorID int, notifType string, refID int, message string, isRead bool, minutesAgo int) {
	read := 0
	if isRead {
		read = 1
	}
	s.exec(
		fmt.Sprintf("insert notification for user %d", userID),
		"INSERT INTO notifications (user_id, actor_id, type, ref_id, message, is_read, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userID, actorID, notifType, refID, message, read, s.millis(minutesAgo),
	)
	s.counts.Notifications++
}

// seedSessions creates one ready-to-use token per user for API testing. The
// rows are part of the same transaction as the rest of the seed.
func (s *seeder) seedSessions() (map[int]string, error) {
	tokens := make(map[int]string, len(seedUsers))
	for i := range seedUsers {
		token, err := newSessionToken()
		if err != nil {
			return nil, err
		}
		userID := i + 1
		expiresAt := s.base.Add(sessionDuration).Format(sqliteTimeLayout)
		if _, err := s.tx.Exec(
			"INSERT INTO user_sessions (session_token, user_id, expires_at) VALUES (?, ?, ?)",
			token, userID, expiresAt,
		); err != nil {
			return nil, fmt.Errorf("seed: insert session for user %d: %w", userID, err)
		}
		tokens[userID] = token
	}
	s.counts.Sessions = len(tokens)
	return tokens, nil
}

func newSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("seed: generate session token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func fullName(userID int) string {
	if userID >= 1 && userID <= len(seedUsers) {
		return seedUsers[userID-1].FullName()
	}
	return "Someone"
}

func demoUsers() []DemoUser {
	users := make([]DemoUser, 0, len(seedUsers))
	for i, u := range seedUsers {
		users = append(users, DemoUser{
			ID:        i + 1,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Nickname:  u.Nickname,
			IsPrivate: u.IsPrivate,
		})
	}
	return users
}

func logSummary(c Counts) {
	log.Printf(
		"Seed: inserted %d users, %d follows, %d groups, %d posts, %d comments, %d reactions, %d events, %d messages, %d notifications, %d reels, %d reel likes, %d reel comments",
		c.Users, c.Follows, c.Groups, c.Posts, c.Comments, c.Reactions,
		c.Events, c.Messages, c.Notifications, c.Reels, c.ReelLikes, c.ReelComments,
	)
}
