// Command seed fills the SQLite database with a coherent demo dataset.
//
// It reuses internal/database so it always targets the same schema as the
// server: database.Init applies any pending migrations, and DB_PATH selects
// the file, so the same binary works locally and inside the container.
//
// Usage:
//
//	go run ./cmd/seed              # wipe app data, then seed
//	go run ./cmd/seed -db /tmp/x.db
//	go run ./cmd/seed -wipe=false  # append instead of wiping
package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"zone/internal/database"
)

// sqliteTimeLayout is the text format go-sqlite3 reads back out of a DATETIME
// column, and the same one the app's own inserts produce.
const sqliteTimeLayout = "2006-01-02 15:04:05"

const sessionDuration = 24 * time.Hour

// wipeOrder lists every application table in reverse dependency order so the
// foreign keys enforced by the _foreign_keys=on DSN stay satisfied.
var wipeOrder = []string{
	"messages", "notifications", "event_responses", "events",
	"post_viewers", "comments", "post_reactions",
	"group_requests", "group_invitations", "group_members", "user_groups",
	"posts", "followers", "user_sessions", "users",
}

var reportOrder = []string{
	"users", "user_sessions", "followers", "user_groups", "group_members",
	"group_invitations", "group_requests", "posts", "post_viewers",
	"comments", "post_reactions", "events", "event_responses", "messages",
	"notifications",
}

func main() {
	log.SetFlags(0)

	dbPath := flag.String("db", "", "database file to seed (defaults to DB_PATH, then the local dev database)")
	wipe := flag.Bool("wipe", true, "delete existing application data before seeding")
	flag.Parse()

	if *dbPath != "" {
		if err := os.Setenv("DB_PATH", *dbPath); err != nil {
			log.Fatalf("set DB_PATH: %v", err)
		}
	}

	if err := database.Init(); err != nil {
		log.Fatalf("init database: %v", err)
	}
	defer database.Database.Close()

	if *wipe {
		if err := wipeAll(database.Database); err != nil {
			log.Fatalf("wipe database: %v", err)
		}
		log.Println("wiped existing application data")
	}

	tokens, err := seed()
	if err != nil {
		log.Fatalf("seed: %v", err)
	}

	printReport(tokens)
}

func wipeAll(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

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

	return tx.Commit()
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
}

func (s *seeder) exec(query string, args ...any) {
	if s.err != nil {
		return
	}
	_, s.err = s.tx.Exec(query, args...)
}

func (s *seeder) execID(query string, args ...any) int64 {
	if s.err != nil {
		return 0
	}
	res, err := s.tx.Exec(query, args...)
	if err != nil {
		s.err = fmt.Errorf("%s: %w", query, err)
		return 0
	}
	id, err := res.LastInsertId()
	if err != nil {
		s.err = err
	}
	return id
}

func (s *seeder) queryID(query string, args ...any) int {
	if s.err != nil {
		return 0
	}
	var id int
	if err := s.tx.QueryRow(query, args...).Scan(&id); err != nil {
		s.err = fmt.Errorf("%s: %w", query, err)
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

func seed() (map[int]string, error) {
	tx, err := database.Database.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	hash, err := bcrypt.GenerateFromPassword([]byte(demoPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash demo password: %w", err)
	}

	s := &seeder{
		tx:            tx,
		base:          time.Now().Truncate(time.Second),
		invitationIDs: map[[2]int]int{},
		requestIDs:    map[[2]int]int{},
	}

	s.seedUsers(string(hash))
	s.seedFollows()
	s.seedGroups()
	s.seedPosts()
	s.seedEvents()
	s.seedMessages()
	s.seedNotifications()

	if s.err != nil {
		return nil, s.err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.seedSessions()
}

// seedUsers inserts seedUsers in order so their ids become 1..len(seedUsers).
func (s *seeder) seedUsers(hash string) {
	for _, u := range seedUsers {
		s.exec(
			`INSERT INTO users (first_name, last_name, email, password, date_of_birth, nickname, about_me, is_private, is_public, created_at, last_seen)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?)`,
			u.FirstName, u.LastName, u.Email, hash, u.DateOfBirth, u.Nickname, u.AboutMe,
			u.IsPrivate, s.ts(u.JoinedHoursAgo), s.tsMinutes(u.LastSeenMinutesAgo),
		)
	}
}

func (s *seeder) seedFollows() {
	for _, f := range seedFollows {
		s.exec(
			"INSERT INTO followers (follower_id, followee_id, status, created_at) VALUES (?, ?, ?, ?)",
			f.From, f.To, f.Status, s.ts(f.HoursAgo),
		)
	}
}

func (s *seeder) seedGroups() {
	for i, g := range seedGroups {
		groupID := int(s.execID(
			"INSERT INTO user_groups (creator_id, title, description) VALUES (?, ?, ?)",
			g.Creator, g.Title, g.Description,
		))
		s.groupIDs = append(s.groupIDs, groupID)
		groupNumber := i + 1

		for j, member := range g.Members {
			// Stagger joins so member ordering is not uniformly the creator first.
			joinedHoursAgo := 24 * (100 - groupNumber*10 - j)
			s.exec(
				"INSERT INTO group_members (group_id, user_id, joined_at) VALUES (?, ?, ?)",
				groupID, member, s.ts(joinedHoursAgo),
			)
		}

		for _, inv := range g.Invitations {
			id := int(s.execID(
				"INSERT INTO group_invitations (group_id, inviter_id, invited_user_id, status, created_at) VALUES (?, ?, ?, ?, ?)",
				groupID, inv.Inviter, inv.Invitee, inv.Status, s.ts(inv.HoursAgo),
			))
			s.invitationIDs[[2]int{groupID, inv.Invitee}] = id
		}

		for _, req := range g.Requests {
			id := int(s.execID(
				"INSERT INTO group_requests (group_id, user_id, status, created_at) VALUES (?, ?, ?, ?)",
				groupID, req.User, req.Status, s.ts(req.HoursAgo),
			))
			s.requestIDs[[2]int{groupID, req.User}] = id
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
			"INSERT INTO posts (user_id, group_id, title, content, image, privacy, created_at) VALUES (?, ?, ?, ?, NULL, ?, ?)",
			p.Author, groupID, p.Title, p.Content, p.Privacy, s.ts(p.HoursAgo),
		))

		for _, viewer := range p.Viewers {
			s.exec("INSERT INTO post_viewers (post_id, user_id) VALUES (?, ?)", postID, viewer)
		}
		for _, r := range p.Reactions {
			isLike := 0
			if r.Like {
				isLike = 1
			}
			s.exec("INSERT INTO post_reactions (user_id, post_id, is_like) VALUES (?, ?, ?)", r.User, postID, isLike)
		}
		for _, c := range p.Comments {
			var commentGroup any
			if p.Group > 0 {
				commentGroup = p.Group
			}
			s.exec(
				"INSERT INTO comments (post_id, user_id, group_id, content, image, created_at) VALUES (?, ?, ?, ?, NULL, ?)",
				postID, c.Author, commentGroup, c.Content, s.ts(c.HoursAgo),
			)
		}
	}
}

func (s *seeder) seedEvents() {
	for _, e := range seedEvents {
		options := fmt.Sprintf(`["%s","%s"]`, e.Options[0], e.Options[1])
		eventDate := s.base.AddDate(0, 0, e.InDays).Format("2006-01-02")

		eventID := int(s.execID(
			"INSERT INTO events (group_id, title, description, event_date, event_time, created_by, options, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			e.Group, e.Title, e.Description, eventDate, e.At, e.Creator, options, s.ts(e.CreatedHoursAgo),
		))
		s.eventIDs = append(s.eventIDs, eventID)

		for _, v := range e.Votes {
			s.exec(
				"INSERT INTO event_responses (event_id, user_id, response) VALUES (?, ?, ?)",
				eventID, v.User, e.Options[v.OptionIdx],
			)
		}
	}
}

func (s *seeder) seedMessages() {
	for _, m := range seedMessages {
		s.exec(
			"INSERT INTO messages (sender_id, type, ref_id, content, created_at) VALUES (?, 'private', ?, ?, ?)",
			m.From, m.To, m.Content, s.millis(m.MinutesAgo),
		)
	}
	for _, m := range seedGroupMessages {
		s.exec(
			"INSERT INTO messages (sender_id, type, ref_id, content, created_at) VALUES (?, 'group', ?, ?, ?)",
			m.From, m.Group, m.Content, s.millis(m.MinutesAgo),
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
		"INSERT INTO notifications (user_id, actor_id, type, ref_id, message, is_read, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userID, actorID, notifType, refID, message, read, s.millis(minutesAgo),
	)
}

// seedSessions runs after commit so a session row is never left behind by a
// rollback, and returns one ready-to-use token per user for API testing.
func (s *seeder) seedSessions() (map[int]string, error) {
	tokens := make(map[int]string, len(seedUsers))
	for i := range seedUsers {
		token, err := newSessionToken()
		if err != nil {
			return nil, err
		}
		userID := i + 1
		expiresAt := time.Now().Add(sessionDuration).Format(sqliteTimeLayout)
		if _, err := database.Database.Exec(
			"INSERT INTO user_sessions (session_token, user_id, expires_at) VALUES (?, ?, ?)",
			token, userID, expiresAt,
		); err != nil {
			return nil, err
		}
		tokens[userID] = token
	}
	return tokens, nil
}

func newSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate session token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func fullName(userID int) string {
	if userID >= 1 && userID <= len(seedUsers) {
		return seedUsers[userID-1].FullName()
	}
	return "Someone"
}

func printReport(tokens map[int]string) {
	log.Println("\nseeded data:")
	for _, table := range reportOrder {
		var count int
		if err := database.Database.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			log.Printf("  %-20s error: %v", table, err)
			continue
		}
		log.Printf("  %-20s %d", table, count)
	}

	log.Printf("\nlog in as any of these (password: %s):", demoPassword)
	for i, u := range seedUsers {
		kind := "public"
		if u.IsPrivate {
			kind = "private"
		}
		log.Printf("  %-24s %-18s %-10s id=%d", u.Email, u.Nickname, kind, i+1)
	}

	log.Println("\nsession tokens (cookie name: session_token):")
	for i, u := range seedUsers {
		log.Printf("  %-24s %s", u.Nickname, tokens[i+1])
	}
}
