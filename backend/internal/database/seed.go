package database

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// DemoPassword is the shared password for every seeded demo account.
const DemoPassword = "demo1234"

type demoUser struct {
	Email       string
	FirstName   string
	LastName    string
	Nickname    string
	DateOfBirth string
	AboutMe     string
	Private     bool
}

type demoFollow struct {
	Follower int
	Followee int
	Pending  bool
}

type demoGroup struct {
	Creator     int
	Title       string
	Description string
	Members     []int
}

type demoPost struct {
	Author   int
	Group    int
	Title    string
	Content  string
	Privacy  string
	Viewers  []int
	HoursAgo float64
}

type demoComment struct {
	Post     int
	Author   int
	Content  string
	HoursAgo float64
}

type demoReaction struct {
	Post int
	User int
	Like bool
}

type demoEvent struct {
	Group       int
	Creator     int
	Title       string
	Description string
	DaysAhead   int
	Time        string
	Options     []string
	Votes       []demoVote
}

type demoVote struct {
	User     int
	Response string
}

type demoNotification struct {
	Recipient int
	Actor     int
	Type      string
	RefID     int
	Message   string
	HoursAgo  float64
}

type demoMessage struct {
	From     int
	To       int
	Content  string
	HoursAgo float64
}

type demoGroupMessage struct {
	Group    int
	From     int
	Content  string
	HoursAgo float64
}

var demoUsers = []demoUser{
	{Email: "sara@demo.local", FirstName: "Sara", LastName: "Johnson", Nickname: "sara", DateOfBirth: "1997-04-12", AboutMe: "Frontend developer. Coffee first, code second."},
	{Email: "ahmed@demo.local", FirstName: "Ahmed", LastName: "Ali", Nickname: "ahmed", DateOfBirth: "1995-11-03", AboutMe: "Backend engineer. Go, SQLite and distributed systems."},
	{Email: "nicole@demo.local", FirstName: "Nicole", LastName: "Martin", Nickname: "nicole", DateOfBirth: "1998-06-21", AboutMe: "Movie nerd and amateur cook. Ask me for a recommendation."},
	{Email: "alex@demo.local", FirstName: "Alex", LastName: "Brown", Nickname: "alex", DateOfBirth: "1994-02-17", AboutMe: "Runner, gamer and weekend hiker."},
	{Email: "lina@demo.local", FirstName: "Lina", LastName: "Haddad", Nickname: "lina", DateOfBirth: "1999-09-08", AboutMe: "Illustrator and UI designer. Private but friendly.", Private: true},
	{Email: "omar@demo.local", FirstName: "Omar", LastName: "Farouk", Nickname: "omar", DateOfBirth: "1996-12-30", AboutMe: "Street photographer. Mostly analog, always curious.", Private: true},
}

var demoFollows = []demoFollow{
	{Follower: 0, Followee: 1},
	{Follower: 1, Followee: 0},
	{Follower: 0, Followee: 2},
	{Follower: 2, Followee: 0},
	{Follower: 1, Followee: 3},
	{Follower: 3, Followee: 1},
	{Follower: 2, Followee: 3},
	{Follower: 3, Followee: 2},
	{Follower: 0, Followee: 3},
	{Follower: 3, Followee: 0},
	{Follower: 4, Followee: 1},
	{Follower: 5, Followee: 2},
	{Follower: 0, Followee: 4, Pending: true},
	{Follower: 2, Followee: 5, Pending: true},
}

var demoGroups = []demoGroup{
	{Creator: 0, Title: "Go Developers", Description: "A place for Go fans to share code, tools and meetup plans.", Members: []int{0, 1, 3}},
	{Creator: 2, Title: "Movie Night Club", Description: "Weekly movie picks, reviews and cozy watch parties.", Members: []int{2, 0, 3}},
	{Creator: 3, Title: "Weekend Runners", Description: "Training plans, routes and long runs with good company.", Members: []int{3, 2, 1}},
}

var demoPosts = []demoPost{
	{Author: 0, Group: -1, Title: "Hello from the demo feed 👋", Privacy: "public", HoursAgo: 1, Content: "Welcome! This account is seeded with posts, comments, reactions, groups, events and chats so you can explore the whole app right away.\n\nTry following someone, opening a group or sending a message."},
	{Author: 1, Group: -1, Title: "What are you building this weekend?", Privacy: "public", HoursAgo: 2, Content: "I am rewriting a small side project and finally adding realtime updates. Curious what everyone else has on their desk."},
	{Author: 2, Group: -1, Title: "Three underrated sci-fi movies", Privacy: "public", HoursAgo: 3, Content: "Arrival, Ex Machina and Coherence. All three are tense, smart and very rewatchable.\n\nWhat is your most underrated pick?"},
	{Author: 3, Group: -1, Title: "Marathon training: week 6", Privacy: "public", HoursAgo: 4, Content: "Long run done before breakfast. The plan is finally feeling manageable, and the rest days are doing their job."},
	{Author: 0, Group: -1, Title: "Weekend plan: coffee and code", Privacy: "almost_private", HoursAgo: 5, Content: "Keeping this one for my followers: coffee, a small side project and absolutely no meetings on Saturday."},
	{Author: 2, Group: -1, Title: "Small circle movie night", Privacy: "private", Viewers: []int{0, 3}, HoursAgo: 6, Content: "Private post for the movie crew. I have a surprise pick for Friday and snacks are already covered."},
	{Author: 4, Group: -1, Title: "New illustration series", Privacy: "public", HoursAgo: 7, Content: "I just finished a set of warm, grainy illustrations. Sharing a few soon once the colors settle."},
	{Author: 5, Group: -1, Title: "Morning light in the old town", Privacy: "public", HoursAgo: 8, Content: "Walked around with one lens and half a roll of film. The quiet streets before 8am are the best part of the city."},
	{Author: 1, Group: -1, Title: "Looking for Go beta testers", Privacy: "almost_private", HoursAgo: 9, Content: "Followers only: I need a few people to try a tiny CLI before I publish it. It is small, fast and mostly harmless."},
	{Author: 3, Group: -1, Title: "My favorite running playlist", Privacy: "public", HoursAgo: 10, Content: "Synthwave for the first half, something calmer for the cooldown. It makes the long runs feel shorter."},
	{Author: 0, Group: 0, Title: "Welcome to Go Developers", Privacy: "public", HoursAgo: 11, Content: "Introduce yourself and share what you are building. Beginners are very welcome here."},
	{Author: 1, Group: 0, Title: "Favorite Go package?", Privacy: "public", HoursAgo: 12, Content: "Mine is still the standard library net/http. Simple, boring and dependable. What do you reach for first?"},
	{Author: 2, Group: 1, Title: "This week's movie pick", Privacy: "public", HoursAgo: 13, Content: "Arrival is the pick. Slow, beautiful and worth discussing afterwards."},
	{Author: 3, Group: 2, Title: "Saturday route", Privacy: "public", HoursAgo: 14, Content: "Twelve kilometers along the river, easy pace, coffee at the end. Meet at the usual bridge."},
}

var demoComments = []demoComment{
	{Post: 0, Author: 1, Content: "Welcome aboard! The feed looks alive already.", HoursAgo: 0.9},
	{Post: 0, Author: 2, Content: "Love the clean layout. The demo data makes it easy to explore.", HoursAgo: 0.8},
	{Post: 1, Author: 0, Content: "A tiny CLI for tracking reading time. Nothing fancy, but useful.", HoursAgo: 0.7},
	{Post: 2, Author: 0, Content: "Arrival is on my rewatch list. That ending stays with you.", HoursAgo: 0.6},
	{Post: 10, Author: 1, Content: "Glad to be here. I will share a websocket example soon.", HoursAgo: 0.5},
	{Post: 12, Author: 0, Content: "I am in for Friday. I will bring the popcorn.", HoursAgo: 0.4},
	{Post: 13, Author: 1, Content: "See you at the bridge. Slow pace for me this week.", HoursAgo: 0.3},
}

var demoReactions = []demoReaction{
	{Post: 0, User: 1, Like: true},
	{Post: 0, User: 2, Like: true},
	{Post: 0, User: 3, Like: true},
	{Post: 1, User: 0, Like: true},
	{Post: 1, User: 3, Like: true},
	{Post: 1, User: 2, Like: false},
	{Post: 2, User: 0, Like: true},
	{Post: 2, User: 1, Like: true},
	{Post: 3, User: 0, Like: true},
	{Post: 3, User: 2, Like: true},
	{Post: 4, User: 1, Like: true},
	{Post: 4, User: 2, Like: true},
	{Post: 5, User: 0, Like: true},
	{Post: 5, User: 3, Like: true},
	{Post: 6, User: 0, Like: true},
	{Post: 6, User: 1, Like: true},
	{Post: 6, User: 2, Like: true},
	{Post: 7, User: 0, Like: true},
	{Post: 7, User: 2, Like: true},
	{Post: 7, User: 3, Like: true},
	{Post: 8, User: 0, Like: true},
	{Post: 8, User: 3, Like: true},
	{Post: 9, User: 1, Like: true},
	{Post: 9, User: 2, Like: true},
	{Post: 10, User: 1, Like: true},
	{Post: 10, User: 3, Like: true},
	{Post: 11, User: 0, Like: true},
	{Post: 11, User: 3, Like: true},
	{Post: 12, User: 0, Like: true},
	{Post: 12, User: 3, Like: true},
	{Post: 13, User: 1, Like: true},
	{Post: 13, User: 2, Like: true},
}

var demoEvents = []demoEvent{
	{Group: 0, Creator: 0, Title: "Go meetup: realtime APIs", Description: "Short talks and a friendly discussion about websockets.", DaysAhead: 7, Time: "18:30", Options: []string{"Going", "Not going"}, Votes: []demoVote{{User: 1, Response: "Going"}, {User: 3, Response: "Not going"}}},
	{Group: 1, Creator: 2, Title: "Friday movie night", Description: "Arrival on the big screen, snacks included.", DaysAhead: 3, Time: "20:00", Options: []string{"Going", "Not going"}, Votes: []demoVote{{User: 0, Response: "Going"}, {User: 3, Response: "Going"}}},
	{Group: 2, Creator: 3, Title: "Sunday long run", Description: "Twelve easy kilometers along the river.", DaysAhead: 5, Time: "07:30", Options: []string{"Going", "Not going"}, Votes: []demoVote{{User: 2, Response: "Going"}, {User: 1, Response: "Going"}}},
}

// RefID points at the related invitation, join request, event or actor depending on the notification type.
var demoNotifications = []demoNotification{
	{Recipient: 4, Actor: 0, Type: "follow_request", RefID: 0, Message: "New follow request from Sara Johnson", HoursAgo: 2},
	{Recipient: 5, Actor: 0, Type: "group_invitation", RefID: 0, Message: "Sara Johnson invited you to join Go Developers", HoursAgo: 3},
	{Recipient: 4, Actor: 3, Type: "group_invitation", RefID: 1, Message: "Alex Brown invited you to join Weekend Runners", HoursAgo: 4},
	{Recipient: 2, Actor: 4, Type: "join_request", RefID: 0, Message: "New join request for your group Lina Haddad", HoursAgo: 5},
	{Recipient: 1, Actor: 0, Type: "event_created", RefID: 0, Message: "New event by Sara Johnson in your group: Go meetup: realtime APIs", HoursAgo: 6},
	{Recipient: 3, Actor: 0, Type: "event_created", RefID: 0, Message: "New event by Sara Johnson in your group: Go meetup: realtime APIs", HoursAgo: 6},
	{Recipient: 0, Actor: 2, Type: "event_created", RefID: 1, Message: "New event by Nicole Martin in your group: Friday movie night", HoursAgo: 7},
	{Recipient: 3, Actor: 2, Type: "event_created", RefID: 1, Message: "New event by Nicole Martin in your group: Friday movie night", HoursAgo: 7},
	{Recipient: 2, Actor: 3, Type: "event_created", RefID: 2, Message: "New event by Alex Brown in your group: Sunday long run", HoursAgo: 8},
	{Recipient: 1, Actor: 3, Type: "event_created", RefID: 2, Message: "New event by Alex Brown in your group: Sunday long run", HoursAgo: 8},
}

var demoMessages = []demoMessage{
	{From: 0, To: 1, Content: "Hey Ahmed! Did you see the new Go release?", HoursAgo: 3.0},
	{From: 1, To: 0, Content: "Just read the notes 😄 The new tooling looks great.", HoursAgo: 2.9},
	{From: 0, To: 1, Content: "Want to pair on a small demo this weekend?", HoursAgo: 2.8},
	{From: 1, To: 0, Content: "Absolutely. Saturday morning works for me.", HoursAgo: 2.7},
	{From: 1, To: 0, Content: "I will bring coffee ☕", HoursAgo: 2.6},
	{From: 2, To: 0, Content: "Movie night on Friday?", HoursAgo: 6.0},
	{From: 0, To: 2, Content: "Yes! I will bring snacks 🍿", HoursAgo: 5.9},
	{From: 2, To: 0, Content: "Perfect. Alex is joining too.", HoursAgo: 5.8},
	{From: 3, To: 0, Content: "Your marathon post is motivating me.", HoursAgo: 8.0},
	{From: 0, To: 3, Content: "Come run with us on Sunday!", HoursAgo: 7.9},
	{From: 2, To: 3, Content: "Did you finish that series?", HoursAgo: 10.0},
	{From: 3, To: 2, Content: "Two episodes left. No spoilers!", HoursAgo: 9.9},
	{From: 5, To: 2, Content: "I printed that photo you liked.", HoursAgo: 12.0},
	{From: 2, To: 5, Content: "Can't wait to see it!", HoursAgo: 11.9},
}

var demoGroupMessages = []demoGroupMessage{
	{Group: 0, From: 0, Content: "Welcome to Go Developers! Share what you are building.", HoursAgo: 1.5},
	{Group: 0, From: 1, Content: "Working on a tiny websocket service.", HoursAgo: 1.4},
	{Group: 0, From: 3, Content: "Count me in for the next meetup.", HoursAgo: 1.3},
	{Group: 0, From: 0, Content: "I will post the event details today.", HoursAgo: 1.2},
	{Group: 1, From: 2, Content: "This week's pick is Arrival.", HoursAgo: 2.5},
	{Group: 1, From: 0, Content: "Great choice. I will bring popcorn 🍿", HoursAgo: 2.4},
	{Group: 1, From: 3, Content: "I have never seen it. No spoilers!", HoursAgo: 2.3},
	{Group: 2, From: 3, Content: "Sunday route is 12km along the river.", HoursAgo: 4.0},
	{Group: 2, From: 2, Content: "I will be there. Slow pace for me.", HoursAgo: 3.9},
	{Group: 2, From: 1, Content: "Joining too, see you at 7:30.", HoursAgo: 3.8},
}

// Seed inserts demo data once, on the first startup that has no demo accounts.
func Seed() error {
	tx, err := Database.Begin()
	if err != nil {
		return fmt.Errorf("seed: begin tx: %w", err)
	}
	defer tx.Rollback()

	var existing int
	if err := tx.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", demoUsers[0].Email).Scan(&existing); err != nil {
		return fmt.Errorf("seed: check demo account: %w", err)
	}
	if existing > 0 {
		log.Println("Seed: demo data already present, skipping")
		return nil
	}

	now := time.Now()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(DemoPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("seed: hash demo password: %w", err)
	}

	userIDs := make([]int64, len(demoUsers))
	for i, u := range demoUsers {
		createdAt := now.Add(-time.Duration(i+1) * 24 * time.Hour).Format("2006-01-02 15:04:05")
		res, err := tx.Exec(
			`INSERT INTO users (first_name, last_name, email, password, date_of_birth, nickname, about_me, is_private, created_at, last_seen)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			u.FirstName, u.LastName, u.Email, string(passwordHash), u.DateOfBirth, u.Nickname, u.AboutMe, u.Private, createdAt, createdAt,
		)
		if err != nil {
			return fmt.Errorf("seed: insert user %s: %w", u.Nickname, err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("seed: user id %s: %w", u.Nickname, err)
		}
		userIDs[i] = id
	}

	for _, f := range demoFollows {
		status := "accepted"
		if f.Pending {
			status = "pending"
		}
		if _, err := tx.Exec(
			"INSERT INTO followers (follower_id, followee_id, status) VALUES (?, ?, ?)",
			userIDs[f.Follower], userIDs[f.Followee], status,
		); err != nil {
			return fmt.Errorf("seed: insert follow %d -> %d: %w", f.Follower, f.Followee, err)
		}
	}

	groupIDs := make([]int64, len(demoGroups))
	for i, g := range demoGroups {
		res, err := tx.Exec(
			"INSERT INTO user_groups (creator_id, title, description) VALUES (?, ?, ?)",
			userIDs[g.Creator], g.Title, g.Description,
		)
		if err != nil {
			return fmt.Errorf("seed: insert group %s: %w", g.Title, err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("seed: group id %s: %w", g.Title, err)
		}
		groupIDs[i] = id
		for _, member := range g.Members {
			if _, err := tx.Exec(
				"INSERT INTO group_members (group_id, user_id) VALUES (?, ?)",
				id, userIDs[member],
			); err != nil {
				return fmt.Errorf("seed: add %d to group %s: %w", member, g.Title, err)
			}
		}
	}

	invitationIDs := make([]int64, 2)
	invitations := []struct {
		Group   int
		Inviter int
		Invited int
	}{
		{Group: 0, Inviter: 0, Invited: 5},
		{Group: 2, Inviter: 3, Invited: 4},
	}
	for i, inv := range invitations {
		res, err := tx.Exec(
			"INSERT INTO group_invitations (group_id, inviter_id, invited_user_id) VALUES (?, ?, ?)",
			groupIDs[inv.Group], userIDs[inv.Inviter], userIDs[inv.Invited],
		)
		if err != nil {
			return fmt.Errorf("seed: insert invitation %d: %w", i, err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("seed: invitation id %d: %w", i, err)
		}
		invitationIDs[i] = id
	}

	res, err := tx.Exec(
		"INSERT INTO group_requests (group_id, user_id) VALUES (?, ?)",
		groupIDs[1], userIDs[4],
	)
	if err != nil {
		return fmt.Errorf("seed: insert join request: %w", err)
	}
	requestID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("seed: join request id: %w", err)
	}

	postIDs := make([]int64, len(demoPosts))
	for i, p := range demoPosts {
		var groupID any
		if p.Group >= 0 {
			groupID = groupIDs[p.Group]
		}
		createdAt := now.Add(-time.Duration(p.HoursAgo * float64(time.Hour))).Format("2006-01-02 15:04:05")
		res, err := tx.Exec(
			"INSERT INTO posts (user_id, group_id, title, content, privacy, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			userIDs[p.Author], groupID, p.Title, p.Content, p.Privacy, createdAt,
		)
		if err != nil {
			return fmt.Errorf("seed: insert post %d: %w", i, err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("seed: post id %d: %w", i, err)
		}
		postIDs[i] = id
		for _, viewer := range p.Viewers {
			if _, err := tx.Exec(
				"INSERT INTO post_viewers (post_id, user_id) VALUES (?, ?)",
				id, userIDs[viewer],
			); err != nil {
				return fmt.Errorf("seed: post %d viewer %d: %w", i, viewer, err)
			}
		}
	}

	for i, c := range demoComments {
		var groupID any
		if demoPosts[c.Post].Group >= 0 {
			groupID = groupIDs[demoPosts[c.Post].Group]
		}
		createdAt := now.Add(-time.Duration(c.HoursAgo * float64(time.Hour))).Format("2006-01-02 15:04:05")
		if _, err := tx.Exec(
			"INSERT INTO comments (post_id, group_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)",
			postIDs[c.Post], groupID, userIDs[c.Author], c.Content, createdAt,
		); err != nil {
			return fmt.Errorf("seed: insert comment %d: %w", i, err)
		}
	}

	for _, r := range demoReactions {
		isLike := 0
		if r.Like {
			isLike = 1
		}
		if _, err := tx.Exec(
			"INSERT INTO post_reactions (user_id, post_id, is_like) VALUES (?, ?, ?)",
			userIDs[r.User], postIDs[r.Post], isLike,
		); err != nil {
			return fmt.Errorf("seed: insert reaction post %d user %d: %w", r.Post, r.User, err)
		}
	}

	eventIDs := make([]int64, len(demoEvents))
	for i, e := range demoEvents {
		eventDate := now.AddDate(0, 0, e.DaysAhead).Format("2006-01-02")
		options := fmt.Sprintf(`["%s","%s"]`, e.Options[0], e.Options[1])
		res, err := tx.Exec(
			`INSERT INTO events (group_id, title, description, event_date, event_time, created_by, options)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			groupIDs[e.Group], e.Title, e.Description, eventDate, e.Time, userIDs[e.Creator], options,
		)
		if err != nil {
			return fmt.Errorf("seed: insert event %s: %w", e.Title, err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("seed: event id %s: %w", e.Title, err)
		}
		eventIDs[i] = id
		for _, vote := range e.Votes {
			if _, err := tx.Exec(
				"INSERT INTO event_responses (event_id, user_id, response) VALUES (?, ?, ?)",
				id, userIDs[vote.User], vote.Response,
			); err != nil {
				return fmt.Errorf("seed: event %s vote: %w", e.Title, err)
			}
		}
	}

	for i, n := range demoNotifications {
		refID := int64(n.RefID)
		switch n.Type {
		case "group_invitation":
			refID = invitationIDs[n.RefID]
		case "join_request":
			refID = requestID
		case "event_created":
			refID = eventIDs[n.RefID]
		case "follow_request":
			refID = userIDs[n.Actor]
		}
		createdAt := now.Add(-time.Duration(n.HoursAgo * float64(time.Hour))).UnixMilli()
		if _, err := tx.Exec(
			"INSERT INTO notifications (user_id, actor_id, type, ref_id, message, is_read, created_at) VALUES (?, ?, ?, ?, ?, 0, ?)",
			userIDs[n.Recipient], userIDs[n.Actor], n.Type, refID, n.Message, createdAt,
		); err != nil {
			return fmt.Errorf("seed: insert notification %d: %w", i, err)
		}
	}

	for i, m := range demoMessages {
		createdAt := now.Add(-time.Duration(m.HoursAgo * float64(time.Hour))).UnixMilli()
		if _, err := tx.Exec(
			"INSERT INTO messages (sender_id, type, ref_id, content, created_at) VALUES (?, 'private', ?, ?, ?)",
			userIDs[m.From], userIDs[m.To], m.Content, createdAt,
		); err != nil {
			return fmt.Errorf("seed: insert message %d: %w", i, err)
		}
	}

	for i, m := range demoGroupMessages {
		createdAt := now.Add(-time.Duration(m.HoursAgo * float64(time.Hour))).UnixMilli()
		if _, err := tx.Exec(
			"INSERT INTO messages (sender_id, type, ref_id, content, created_at) VALUES (?, 'group', ?, ?, ?)",
			userIDs[m.From], groupIDs[m.Group], m.Content, createdAt,
		); err != nil {
			return fmt.Errorf("seed: insert group message %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("seed: commit: %w", err)
	}

	log.Printf("Seed: inserted %d users, %d posts, %d comments, %d groups, %d events, %d messages",
		len(demoUsers), len(demoPosts), len(demoComments), len(demoGroups), len(demoEvents), len(demoMessages)+len(demoGroupMessages))
	return nil
}
