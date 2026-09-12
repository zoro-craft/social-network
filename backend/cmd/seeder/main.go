package main

import (
	cryptorand "crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"zone/internal/database"

	"golang.org/x/crypto/bcrypt"
)

const seedPassword = "123456789"

type seedUser struct {
	first   string
	last    string
	email   string
	nick    string
	about   string
	dob     string
	private bool
}

var seedUsers = []seedUser{
	{"Alex", "Johnson", "alex.johnson@example.com", "alexj", "Coffee enthusiast and weekend hiker.", "1994-05-12", false},
	{"Maya", "Patel", "maya.patel@example.com", "", "Frontend developer who loves photography.", "1999-11-03", false},
	{"Liam", "O'Connor", "liam.oconnor@example.com", "liamoc", "", "1991-02-28", true},
	{"Sofia", "Rossi", "sofia.rossi@example.com", "", "Travel blogger. Currently exploring Italy.", "1996-07-19", false},
	{"Noah", "Kim", "noah.kim@example.com", "noahk", "Game designer by day, gamer by night.", "1993-09-01", false},
	{"Emma", "Dubois", "emma.dubois@example.com", "", "", "2000-01-15", false},
	{"Lucas", "Santos", "lucas.santos@example.com", "", "Music producer. Ask me about synths.", "1995-12-25", false},
	{"Aisha", "Khan", "aisha.khan@example.com", "", "Med student, cat person.", "2001-04-08", true},
	{"Ethan", "Wright", "ethan.wright@example.com", "ethanw", "Marathon runner and data nerd.", "1988-08-30", false},
	{"Olivia", "Martin", "olivia.martin@example.com", "", "Bakery owner. Sourdough is my love language.", "1990-10-11", false},
	{"Mateo", "Garcia", "mateo.garcia@example.com", "", "Cycling the coast this summer.", "1997-03-22", false},
	{"Hannah", "Cohen", "hannah.cohen@example.com", "", "", "1992-06-17", true},
	{"Daniel", "Novak", "daniel.novak@example.com", "", "Physics grad student.", "1998-09-09", false},
	{"Isabella", "Ferrara", "isabella.ferrara@example.com", "bellaf", "Painter and ceramicist.", "1994-04-02", false},
	{"James", "Chen", "james.chen@example.com", "", "Startup founder. Coffee required.", "1987-02-14", false},
	{"Zara", "Ahmed", "zara.ahmed@example.com", "", "Librarian and bookstagrammer.", "1993-11-28", false},
	{"Benjamin", "Turner", "benjamin.turner@example.com", "", "Amateur astronomer.", "1989-05-05", false},
	{"Chloe", "Nguyen", "chloe.nguyen@example.com", "chloeng", "Dancer, loves street food.", "2002-08-21", false},
	{"Omar", "Haddad", "omar.haddad@example.com", "", "", "1991-12-01", true},
	{"Freya", "Wilson", "freya.wilson@example.com", "", "Gardener and small-town journalist.", "1996-06-06", false},
}

var postTitles = []string{
	"Morning thoughts",
	"Finally finished this project",
	"Road trip planning",
	"New to this city, any tips?",
	"Reading recommendation",
	"Just adopted a puppy",
	"Work from home setup",
	"Favourite hiking trail",
	"Cooking experiments",
	"Concert last night",
	"Movie marathon weekend",
	"Looking for book club members",
}

var postContents = []string{
	"Excited to share a small update today. Things are falling into place nicely!",
	"I have been working on this for weeks and I am really proud of how it turned out.",
	"Anyone know good stops along the coast? Planning a two-week drive next month.",
	"Just moved here for work. Suggest places to eat, weekends spots and everything in between.",
	"I just finished this amazing novel. The character arcs are simply incredible.",
	"Meet Biscuit. She is 9 weeks old and already rules the house.",
	"Here is my new setup. Cable management is not perfect but it is cozy.",
	"Trail of the month: the one through the cedar forest. Breathtaking views at the top.",
	"Attempted sourdough from scratch. It rose! Not sure it is edible yet though.",
	"They played all my favourites. Still can not believe how good the live set was.",
	"No spoilers, but that finale was everything. Highly recommend the whole series.",
	"Small group, big discussions. Looking for people who love to read and debate.",
}

var commentContents = []string{
	"That is awesome, keep it up!",
	"Could not agree more.",
	"Thanks for sharing this!",
	"Ha, I can relate to this one.",
	"How did you manage that?",
	"This made my day.",
	"Count me in!",
	"Great post, saving this for later.",
	"Love this so much. More of this please.",
	"Interesting take, I have a different view though.",
	"Congratulations!",
	"Details, please, we need the full story.",
}

var groupTitles = []string{"Weekend Hikers", "Indie Devs", "Book Lovers Club", "Street Food Scouts", "Trail Runners"}
var groupDescs = []string{
	"Group for people who love the outdoors and weekend trail adventures.",
	"Community of indie developers building cool side projects.",
	"Monthly book discussions, friendly debates, and reading challenges.",
	"Reviewing and hunting down the best street food in the city.",
	"Casual trail running group. All paces welcome.",
}

var eventTitles = []string{"Sunrise hike to Eagle Peak", "Monthly code & coffee", "Book discussion night", "Food tour: old town", "Couch to 5k kickoff", "Evening sprint session"}

var groupPosts = make(map[int64][]int64)

func has(vals []int64, v int64) bool {
	for _, x := range vals {
		if x == v {
			return true
		}
	}
	return false
}

func mustExec(tx *sql.Tx, query string, args ...any) sql.Result {
	res, err := tx.Exec(query, args...)
	if err != nil {
		log.Fatalf("query failed: %s\nargs: %v\nerror: %v", query, args, err)
	}
	return res
}

func lastID(res sql.Result) int64 {
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("could not read last insert id: %v", err)
	}
	return id
}

func sessionToken() string {
	b := make([]byte, 32)
	if _, err := cryptorand.Read(b); err != nil {
		log.Fatalf("could not generate session token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func imageMaybe(rng *rand.Rand) any {
	if rng.Intn(3) == 0 {
		return fmt.Sprintf("/uploads/images/post-%d.jpg", rng.Intn(100))
	}
	return nil
}

func oneOf(rng *rand.Rand, options ...string) string {
	return options[rng.Intn(len(options))]
}

func clearTables(tx *sql.Tx) {
	order := []string{
		"notifications",
		"user_sessions",
		"event_responses",
		"events",
		"post_viewers",
		"post_reactions",
		"comments",
		"posts",
		"group_members",
		"group_invitations",
		"group_requests",
		"followers",
		"user_groups",
		"users",
	}
	for _, t := range order {
		mustExec(tx, "DELETE FROM "+t)
	}
	mustExec(tx, "DELETE FROM sqlite_sequence")
}

func seed(wipe bool) {
	if err := database.Init(); err != nil {
		log.Fatalf("could not initialize database: %v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("could not hash password: %v", err)
	}

	tx, err := database.Database.Begin()
	if err != nil {
		log.Fatalf("could not start transaction: %v", err)
	}

	if wipe {
		clearTables(tx)
	}

	rng := rand.New(rand.NewSource(42))
	now := time.Now()

	userIDs := make([]int64, 0, len(seedUsers))
	for _, u := range seedUsers {
		var nick, about, avatar any
		if u.nick != "" {
			nick = u.nick
		}
		if u.about != "" {
			about = u.about
		}
		if rng.Intn(4) > 0 {
			avatar = fmt.Sprintf("/uploads/avatars/avatar-%d.jpg", len(userIDs)+1)
		}
		created := now.Add(-time.Duration((len(userIDs)+3)*7*24) * time.Hour).Format("2006-01-02 15:04:05")
		lastSeen := now.Add(-time.Duration(rng.Intn(48)) * time.Hour).Format("2006-01-02 15:04:05")
		res := mustExec(tx, `INSERT INTO users
			(first_name, last_name, email, password, date_of_birth, nickname, about_me, avatar, is_private, is_public, created_at, last_seen)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			u.first, u.last, u.email, string(hash), u.dob, nick, about, avatar, boolToInt(u.private), boolToInt(!u.private), created, lastSeen)
		userIDs = append(userIDs, lastID(res))
	}

	for _, uid := range userIDs {
		n := 1 + rng.Intn(3)
		for c := 0; c < n; c++ {
			other := userIDs[rng.Intn(len(userIDs))]
			if other == uid {
				continue
			}
			var priv int
			if err := tx.QueryRow("SELECT is_private FROM users WHERE id = ?", other).Scan(&priv); err != nil {
				log.Fatalf("could not read privacy flag: %v", err)
			}
			status := "accepted"
			if priv == 1 {
				status = "pending"
			}
			mustExec(tx, `INSERT OR IGNORE INTO followers (follower_id, followee_id, status, created_at) VALUES (?, ?, ?, ?)`,
				uid, other, status, now.Add(-time.Duration(rng.Intn(60))*24*time.Hour).Format("2006-01-02 15:04:05"))
		}
	}

	groupIDs := make([]int64, 0, len(groupTitles))
	groupMembers := make(map[int64][]int64)
	creators := make(map[int64]int64)
	for gi := range groupTitles {
		creator := userIDs[gi%len(userIDs)]
		res := mustExec(tx, `INSERT INTO user_groups (creator_id, title, description) VALUES (?, ?, ?)`,
			creator, groupTitles[gi], groupDescs[gi])
		gid := lastID(res)
		groupIDs = append(groupIDs, gid)
		creators[gid] = creator
		groupMembers[gid] = []int64{creator}
	}

	for _, gid := range groupIDs {
		memberCount := 3 + rng.Intn(4)
		for c := 0; c < memberCount; c++ {
			u := userIDs[rng.Intn(len(userIDs))]
			if has(groupMembers[gid], u) {
				continue
			}
			mustExec(tx, `INSERT INTO group_members (group_id, user_id, joined_at) VALUES (?, ?, ?)`,
				gid, u, now.Add(-time.Duration(rng.Intn(40))*24*time.Hour).Format("2006-01-02 15:04:05"))
			groupMembers[gid] = append(groupMembers[gid], u)
		}
	}

	for _, gid := range groupIDs {
		groupName := groupTitles[giFor(gid, groupIDs)]

		invites := 0
		for _, u := range userIDs {
			if invites >= 2 {
				break
			}
			if has(groupMembers[gid], u) || u == creators[gid] {
				continue
			}
			res := mustExec(tx, `INSERT INTO group_invitations (group_id, inviter_id, invited_user_id, status, created_at)
				VALUES (?, ?, ?, 'pending', ?)`,
				gid, creators[gid], u, now.Add(-time.Duration(rng.Intn(15))*24*time.Hour).Format("2006-01-02 15:04:05"))
			invID := lastID(res)
			mustExec(tx, `INSERT INTO notifications (user_id, type, invitation_id, request_id, message, created_at)
				VALUES (?, 'group_invtation', ?, NULL, ?, ?)`,
				u, invID, "You are invited to join "+groupName, now.Add(-time.Duration(rng.Intn(15))*24*time.Hour).Format("2006-01-02 15:04:05"))
			invites++
		}

		requests := 0
		for _, u := range userIDs {
			if requests >= 2 {
				break
			}
			if has(groupMembers[gid], u) || u == creators[gid] {
				continue
			}
			res := mustExec(tx, `INSERT INTO group_requests (group_id, user_id, status, created_at)
				VALUES (?, ?, 'pending', ?)`,
				gid, u, now.Add(-time.Duration(rng.Intn(10))*24*time.Hour).Format("2006-01-02 15:04:05"))
			reqID := lastID(res)
			mustExec(tx, `INSERT INTO notifications (user_id, type, invitation_id, request_id, message, created_at)
				VALUES (?, 'join_request', NULL, ?, ?, ?)`,
				creators[gid], reqID, "New join request for your group "+groupName, now.Add(-time.Duration(rng.Intn(10))*24*time.Hour).Format("2006-01-02 15:04:05"))
			requests++
		}

		postCount := 2 + rng.Intn(2)
		for p := 0; p < postCount; p++ {
			member := groupMembers[gid][rng.Intn(len(groupMembers[gid]))]
			res := mustExec(tx, `INSERT INTO posts (user_id, title, content, image, privacy, group_id, created_at)
				VALUES (?, ?, ?, ?, 'public', ?, ?)`,
				member, postTitles[rng.Intn(len(postTitles))], postContents[rng.Intn(len(postContents))], imageMaybe(rng), gid,
				now.Add(-time.Duration(rng.Intn(30))*24*time.Hour).Format("2006-01-02 15:04:05"))
			groupPosts[gid] = append(groupPosts[gid], lastID(res))
		}
	}

	postIDs := make([]int64, 0)
	for _, uid := range userIDs {
		var privacy string
		switch {
		case uid == userIDs[len(userIDs)-1]:
			privacy = "private"
		case rng.Intn(5) == 0:
			privacy = "almost_private"
		default:
			privacy = "public"
		}
		res := mustExec(tx, `INSERT INTO posts (user_id, title, content, image, privacy, group_id, created_at)
			VALUES (?, ?, ?, ?, ?, NULL, ?)`,
			uid, postTitles[rng.Intn(len(postTitles))], postContents[rng.Intn(len(postContents))], imageMaybe(rng), privacy,
			now.Add(-time.Duration(rng.Intn(60))*24*time.Hour).Format("2006-01-02 15:04:05"))
		postIDs = append(postIDs, lastID(res))
	}

	ownerCache := make(map[int64]int64)
	ownerOf := func(postID int64) int64 {
		if o, ok := ownerCache[postID]; ok {
			return o
		}
		var owner int64
		if err := tx.QueryRow("SELECT user_id FROM posts WHERE id = ?", postID).Scan(&owner); err != nil {
			log.Fatalf("could not read post owner: %v", err)
		}
		ownerCache[postID] = owner
		return owner
	}

	for _, pid := range postIDs {
		for _, u := range userIDs {
			if u == ownerOf(pid) {
				continue
			}
			if rng.Intn(4) > 0 {
				continue
			}
			mustExec(tx, `INSERT OR IGNORE INTO post_viewers (post_id, user_id) VALUES (?, ?)`, pid, u)
		}
	}

	allPostIDs := make([]int64, 0, len(postIDs))
	allPostIDs = append(allPostIDs, postIDs...)
	for _, gid := range groupIDs {
		allPostIDs = append(allPostIDs, groupPosts[gid]...)
	}

	for _, pid := range allPostIDs {
		owner := ownerOf(pid)
		reactor1 := userIDs[rng.Intn(len(userIDs))]
		if reactor1 != owner {
			mustExec(tx, `INSERT OR IGNORE INTO post_reactions (user_id, post_id, is_like) VALUES (?, ?, ?)`,
				reactor1, pid, rng.Intn(2))
		}
		reactor2 := userIDs[rng.Intn(len(userIDs))]
		if rng.Intn(2) == 0 && reactor2 != owner && reactor2 != reactor1 {
			mustExec(tx, `INSERT OR IGNORE INTO post_reactions (user_id, post_id, is_like) VALUES (?, ?, ?)`,
				reactor2, pid, rng.Intn(2))
		}
	}

	for _, pid := range allPostIDs {
		owner := ownerOf(pid)
		commentCount := 1 + rng.Intn(4)
		for c := 0; c < commentCount; c++ {
			commenter := userIDs[rng.Intn(len(userIDs))]
			if commenter == owner {
				continue
			}
			var gid any
			_ = tx.QueryRow("SELECT group_id FROM posts WHERE id = ?", pid).Scan(&gid)
			mustExec(tx, `INSERT INTO comments (post_id, user_id, content, image, group_id, created_at)
				VALUES (?, ?, ?, NULL, ?, ?)`,
				pid, commenter, commentContents[rng.Intn(len(commentContents))], gid,
				now.Add(-time.Duration(rng.Intn(20))*24*time.Hour).Format("2006-01-02 15:04:05"))
		}
	}

	eventIDs := make([]int64, 0)
	for _, gid := range groupIDs {
		eventCount := 1 + rng.Intn(2)
		for e := 0; e < eventCount; e++ {
			options := []string{"Going", "Not Going"}
			if rng.Intn(2) == 0 {
				options = []string{"Yes", "No"}
			}
			optionsJSON, err := json.Marshal(options)
			if err != nil {
				log.Fatalf("could not marshal event options: %v", err)
			}
			res := mustExec(tx, `INSERT INTO events (group_id, title, created_by, options, created_at)
				VALUES (?, ?, ?, ?, ?)`, gid, eventTitles[rng.Intn(len(eventTitles))], creators[gid], string(optionsJSON),
				now.Add(-time.Duration(rng.Intn(25))*24*time.Hour).Format("2006-01-02 15:04:05"))
			eventIDs = append(eventIDs, lastID(res))
		}
	}

	for _, eid := range eventIDs {
		var gid int64
		if err := tx.QueryRow("SELECT group_id FROM events WHERE id = ?", eid).Scan(&gid); err != nil {
			log.Fatalf("could not read event group: %v", err)
		}
		members := groupMembers[gid]
		voters := 0
		for c := 0; c < len(members); c++ {
			if rng.Intn(3) == 0 {
				continue
			}
			mustExec(tx, `INSERT OR IGNORE INTO event_responses (event_id, user_id, response) VALUES (?, ?, ?)`,
				eid, members[c], oneOf(rng, "Going", "Not Going", "Maybe", "Yes", "No"))
			voters++
			if voters >= 4+len(members)/2 {
				break
			}
		}
	}

	for i := 0; i < 5 && i < len(userIDs); i++ {
		expires := now.Add(24 * time.Hour).Format("2006-01-02 15:04:05")
		mustExec(tx, `INSERT INTO user_sessions (session_token, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)`,
			sessionToken(), userIDs[i], expires, now.Format("2006-01-02 15:04:05"))
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("could not commit transaction: %v", err)
	}

	fmt.Printf("Seeding complete.\n")
	fmt.Printf("  users:         %d\n", len(userIDs))
	fmt.Printf("  followers:     %d\n", countRows("followers"))
	fmt.Printf("  groups:        %d\n", len(groupIDs))
	fmt.Printf("  posts:         %d\n", len(allPostIDs))
	fmt.Printf("  comments:      %d\n", countRows("comments"))
	fmt.Printf("  events:        %d\n", len(eventIDs))
	fmt.Printf("  notifications: %d\n", countRows("notifications"))
	fmt.Printf("All seeded users have password: %s\n", seedPassword)
}

func main() {
	wipe := flag.Bool("wipe", true, "clear all existing rows before seeding")
	flag.Parse()
	seed(*wipe)
}

func giFor(gid int64, groupIDs []int64) int {
	for i, id := range groupIDs {
		if id == gid {
			return i
		}
	}
	return 0
}

func countRows(table string) int {
	var n int
	_ = database.Database.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&n)
	return n
}