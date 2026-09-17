package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"zone/internal/database"
)

const seedPassword = "123456789"

type seedUser struct {
	FirstName string
	LastName  string
	Email     string
	DOB       string
	Nickname  string
	AboutMe   string
	IsPrivate bool
}

var seedUsers = []seedUser{
	{FirstName: "Alice", LastName: "Johnson", Email: "alice@example.com", DOB: "1995-04-12", Nickname: "alice", AboutMe: "Coffee-fueled developer.", IsPrivate: false},
	{FirstName: "Bob", LastName: "Smith", Email: "bob@example.com", DOB: "1990-08-23", Nickname: "bob", AboutMe: "Backend engineer and Go enthusiast.", IsPrivate: false},
	{FirstName: "Carol", LastName: "White", Email: "carol@example.com", DOB: "1998-01-05", Nickname: "carol", AboutMe: "Photographer by day, reader by night.", IsPrivate: true},
	{FirstName: "David", LastName: "Brown", Email: "david@example.com", DOB: "1992-11-30", Nickname: "david", AboutMe: "Just here to learn.", IsPrivate: false},
	{FirstName: "Emma", LastName: "Wilson", Email: "emma@example.com", DOB: "1996-06-19", Nickname: "emma", AboutMe: "Frontend dev, cat lover.", IsPrivate: false},
	{FirstName: "Frank", LastName: "Miller", Email: "frank@example.com", DOB: "1988-03-08", Nickname: "frank", AboutMe: "Privacy matters.", IsPrivate: true},
	{FirstName: "Grace", LastName: "Davis", Email: "grace@example.com", DOB: "1999-09-14", Nickname: "grace", AboutMe: "Student & hobbyist writer.", IsPrivate: false},
	{FirstName: "Henry", LastName: "Taylor", Email: "henry@example.com", DOB: "1993-12-02", Nickname: "henry", AboutMe: "Book club organizer.", IsPrivate: false},
}

type follow struct {
	FollowerID int
	FolloweeID int
	Status     string
}

var seedFollows = []follow{
	{1, 2, "accepted"},
	{2, 1, "accepted"},
	{1, 3, "accepted"},
	{3, 1, "pending"},
	{1, 4, "accepted"},
	{2, 3, "accepted"},
	{4, 5, "accepted"},
	{5, 1, "pending"},
	{5, 6, "pending"},
	{2, 6, "accepted"},
	{7, 3, "pending"},
	{8, 1, "accepted"},
	{6, 8, "accepted"},
}

type seedPost struct {
	UserID   int
	GroupID  *int
	Title    string
	Content  string
	Privacy  string
	ViewerID []int
}

func p(userID int, title, content, privacy string, viewers ...int) seedPost {
	return seedPost{UserID: userID, Title: title, Content: content, Privacy: privacy, ViewerID: viewers}
}

func gp(userID, groupID int, title, content string) seedPost {
	g := groupID
	return seedPost{UserID: userID, GroupID: &g, Title: title, Content: content, Privacy: "public"}
}

var seedPosts = []seedPost{
	p(1, "First post", "Hello world! This is my very first post.", "public"),
	p(1, "Go tips", "Always profile your code before optimizing.", "almost_private"),
	p(1, "Private diary entry", "Today was a good day.", "private", 2),
	p(2, "Backend journey", "Learning Go has been amazing so far.", "public"),
	p(3, "Secret thoughts", "Hidden note visible only to selected readers.", "private", 1),
	p(3, "Photo update", "Sunset photography is my favorite subject.", "public"),
	p(4, "Greetings", "Happy to be part of this community.", "public"),
	p(5, "Almost private note", "Sharing with close friends only.", "almost_private"),
	p(8, "Book recommendation", "Just finished 'The Pragmatic Programmer', a must read.", "public"),
	gp(1, 1, "Group announcement", "Welcome to the Go Developers group!"),
	gp(3, 2, "Theme of the week", "This week's theme is black and white photography."),
}

type seedComment struct {
	PostID  int
	UserID  int
	GroupID *int
	Content string
}

var seedComments = []seedComment{
	{PostID: 1, UserID: 2, Content: "Nice first post!"},
	{PostID: 1, UserID: 4, Content: "Welcome to the platform!"},
	{PostID: 4, UserID: 1, Content: "Great insight, keep it up!"},
	{PostID: 6, UserID: 1, Content: "Love this!"},
	{PostID: 10, UserID: 2, GroupID: ip(1), Content: "Thanks for the update!"},
	{PostID: 11, UserID: 5, GroupID: ip(2), Content: "Great theme, count me in."},
}

type seedReaction struct {
	PostID int
	UserID int
	IsLike bool
}

var seedReactions = []seedReaction{
	{1, 2, true},
	{1, 3, true},
	{2, 2, true},
	{4, 1, true},
	{4, 3, true},
	{6, 1, false},
	{7, 1, true},
	{10, 2, true},
	{10, 4, true},
	{11, 6, true},
}

type seedGroup struct {
	CreatorID   int
	Title       string
	Description string
}

var seedGroups = []seedGroup{
	{CreatorID: 1, Title: "Go Developers", Description: "A community for Go enthusiasts."},
	{CreatorID: 3, Title: "Photography Club", Description: "Share your best shots."},
	{CreatorID: 8, Title: "Book Lovers", Description: "Reading and discussing great books."},
}

var seedGroupMembers = [][]int{
	{1, 1},
	{1, 2},
	{1, 4},
	{2, 3},
	{2, 5},
	{2, 6},
	{3, 8},
	{3, 7},
	{3, 1},
}

type seedInvite struct {
	GroupID int
	Inviter int
	Invitee int
	Status  string
}

var seedInvites = []seedInvite{
	{1, 1, 7, "pending"},
	{2, 3, 2, "accepted"},
	{3, 8, 5, "pending"},
}

type seedRequest struct {
	GroupID int
	UserID  int
	Status  string
}

var seedRequests = []seedRequest{
	{2, 7, "pending"},
	{1, 8, "pending"},
	{3, 6, "accepted"},
}

type seedEvent struct {
	GroupID     int
	Title       string
	Description string
	EventDate   string
	EventTime   string
	CreatedBy   int
	Options     []string
}

var seedEvents = []seedEvent{
	{GroupID: 1, Title: "Go Meetup", Description: "Monthly meetup", EventDate: "2026-10-01", EventTime: "18:00", CreatedBy: 1, Options: []string{"Online", "In-person"}},
	{GroupID: 2, Title: "Photo Walk", Description: "City walk photo session", EventDate: "2026-09-25", EventTime: "09:00", CreatedBy: 3, Options: []string{"Morning", "Afternoon"}},
	{GroupID: 3, Title: "Book Discussion", Description: "Discussing the monthly pick", EventDate: "2026-09-20", EventTime: "20:00", CreatedBy: 8, Options: []string{"Chapter 1-5", "Whole book"}},
}

type seedVote struct {
	EventID  int
	UserID   int
	Response string
}

var seedVotes = []seedVote{
	{1, 1, "In-person"},
	{1, 2, "Online"},
	{1, 4, "In-person"},
	{2, 3, "Morning"},
	{2, 5, "Morning"},
	{2, 6, "Afternoon"},
	{3, 8, "Whole book"},
	{3, 1, "Chapter 1-5"},
	{3, 7, "Whole book"},
	{3, 6, "Chapter 1-5"},
}

type seedNotification struct {
	UserID       int
	Type         string
	InvitationID *int
	RequestID    *int
	Message      string
}

var seedNotifications = []seedNotification{
	{UserID: 7, Type: "group_invtation", InvitationID: ip(1), Message: "You have been invited to join Go Developers."},
	{UserID: 8, Type: "join_request", RequestID: ip(2), Message: "Henry Taylor wants to join Go Developers."},
	{UserID: 1, Type: "group_invtation", InvitationID: ip(3), Message: "You have been invited to join Book Lovers."},
	{UserID: 3, Type: "join_request", RequestID: ip(1), Message: "Grace Davis wants to join Photography Club."},
}

func ip(n int) *int { return &n }

func mustExec(tx *sql.Tx, label, query string, args ...any) {
	if _, err := tx.Exec(query, args...); err != nil {
		log.Fatalf("%s: %v", label, err)
	}
}

func insertUser(tx *sql.Tx, u seedUser, hashed string) (int, error) {
	isPrivate := 0
	if u.IsPrivate {
		isPrivate = 1
	}
	res, err := tx.Exec(
		`INSERT INTO users (email, password, date_of_birth, first_name, last_name, nickname, about_me, avatar, is_private)
		 VALUES (?, ?, ?, ?, ?, ?, ?, NULL, ?)`,
		u.Email, hashed, u.DOB, u.FirstName, u.LastName, u.Nickname, u.AboutMe, isPrivate,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func reset(tx *sql.Tx) error {
	tables := []string{
		"post_viewers",
		"event_responses",
		"events",
		"notifications",
		"group_requests",
		"group_invitations",
		"group_members",
		"user_groups",
		"comments",
		"post_reactions",
		"posts",
		"followers",
		"user_sessions",
		"users",
	}
	for _, t := range tables {
		if _, err := tx.Exec("DELETE FROM " + t); err != nil {
			return fmt.Errorf("clearing %s: %w", t, err)
		}
	}
	if _, err := tx.Exec("DELETE FROM sqlite_sequence WHERE name IN ('users','posts','comments','user_groups','group_members','group_invitations','group_requests','events','event_responses','notifications')"); err != nil {
		return fmt.Errorf("resetting sequences: %w", err)
	}
	return nil
}

func main() {
	if err := database.Init(); err != nil {
		log.Fatal(err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := database.Database.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	if err := reset(tx); err != nil {
		log.Fatal(err)
	}

	for _, u := range seedUsers {
		if _, err := insertUser(tx, u, string(hashed)); err != nil {
			log.Fatal(err)
		}
	}

	for _, f := range seedFollows {
		mustExec(tx, "seed followers",
			"INSERT INTO followers (follower_id, followee_id, status) VALUES (?, ?, ?)",
			f.FollowerID, f.FolloweeID, f.Status)
	}

	postIDs := make([]int, 0, len(seedPosts))
	for _, p := range seedPosts {
		res, err := tx.Exec(
			"INSERT INTO posts (user_id, group_id, title, content, image, privacy) VALUES (?, ?, ?, ?, NULL, ?)",
			p.UserID, p.GroupID, p.Title, p.Content, p.Privacy,
		)
		if err != nil {
			log.Fatalf("seed posts: %v", err)
		}
		id, _ := res.LastInsertId()
		postIDs = append(postIDs, int(id))
	}

	for i, p := range seedPosts {
		for _, v := range p.ViewerID {
			mustExec(tx, "seed post_viewers",
				"INSERT OR IGNORE INTO post_viewers (post_id, user_id) VALUES (?, ?)",
				postIDs[i], v)
		}
	}

	for _, c := range seedComments {
		mustExec(tx, "seed comments",
			"INSERT INTO comments (post_id, user_id, group_id, content, image) VALUES (?, ?, ?, ?, NULL)",
			c.PostID, c.UserID, c.GroupID, c.Content)
	}

	for _, r := range seedReactions {
		isLike := 0
		if r.IsLike {
			isLike = 1
		}
		mustExec(tx, "seed post_reactions",
			"INSERT OR IGNORE INTO post_reactions (user_id, post_id, is_like) VALUES (?, ?, ?)",
			r.UserID, r.PostID, isLike)
	}

	groupIDs := make([]int, 0, len(seedGroups))
	for _, g := range seedGroups {
		res, err := tx.Exec(
			"INSERT INTO user_groups (creator_id, title, description) VALUES (?, ?, ?)",
			g.CreatorID, g.Title, g.Description,
		)
		if err != nil {
			log.Fatalf("seed user_groups: %v", err)
		}
		id, _ := res.LastInsertId()
		groupIDs = append(groupIDs, int(id))
	}

	for _, m := range seedGroupMembers {
		mustExec(tx, "seed group_members",
			"INSERT OR IGNORE INTO group_members (group_id, user_id) VALUES (?, ?)",
			m[0], m[1])
	}

	for _, inv := range seedInvites {
		mustExec(tx, "seed group_invitations",
			"INSERT INTO group_invitations (group_id, inviter_id, invited_user_id, status) VALUES (?, ?, ?, ?)",
			inv.GroupID, inv.Inviter, inv.Invitee, inv.Status)
	}

	for _, req := range seedRequests {
		mustExec(tx, "seed group_requests",
			"INSERT INTO group_requests (group_id, user_id, status) VALUES (?, ?, ?)",
			req.GroupID, req.UserID, req.Status)
	}

	for _, n := range seedNotifications {
		mustExec(tx, "seed notifications",
			"INSERT INTO notifications (user_id, type, invitation_id, request_id, message, is_read) VALUES (?, ?, ?, ?, ?, 0)",
			n.UserID, n.Type, n.InvitationID, n.RequestID, n.Message)
	}

	eventIDs := make([]int, 0, len(seedEvents))
	for _, e := range seedEvents {
		opts, err := json.Marshal(e.Options)
		if err != nil {
			log.Fatalf("seed events: %v", err)
		}
		res, err := tx.Exec(
			"INSERT INTO events (group_id, title, description, event_date, event_time, created_by, options) VALUES (?, ?, ?, ?, ?, ?, ?)",
			e.GroupID, e.Title, e.Description, e.EventDate, e.EventTime, e.CreatedBy, string(opts),
		)
		if err != nil {
			log.Fatalf("seed events: %v", err)
		}
		id, _ := res.LastInsertId()
		eventIDs = append(eventIDs, int(id))
	}

	for _, v := range seedVotes {
		mustExec(tx, "seed event_responses",
			"INSERT OR IGNORE INTO event_responses (event_id, user_id, response) VALUES (?, ?, ?)",
			v.EventID, v.UserID, v.Response)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Seeding complete. %d users, %d posts, %d comments, %d groups, %d events (password for all users: %s)",
		len(seedUsers), len(postIDs), len(seedComments), len(groupIDs), len(eventIDs), seedPassword)
}
