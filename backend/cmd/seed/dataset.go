package main

// demoPassword is the plaintext password shared by every seeded account.
const demoPassword = "Password123"

// seedUser is a demo account. JoinedHoursAgo and LastSeenMinutesAgo are offsets
// from the seeder's base time, so the dataset always looks freshly active.
type seedUser struct {
	FirstName          string
	LastName           string
	Email              string
	Nickname           string
	AboutMe            string
	DateOfBirth        string
	IsPrivate          bool
	JoinedHoursAgo     int
	LastSeenMinutesAgo int
}

// FullName mirrors the users.full_name generated column.
func (u seedUser) FullName() string { return u.FirstName + " " + u.LastName }

var seedUsers = []seedUser{
	{
		FirstName: "Yasmine", LastName: "Haddad", Email: "yasmine@demo.dev", Nickname: "yasmine",
		AboutMe:     "Backend dev. Go, SQLite, and too much coffee.",
		DateOfBirth: "1998-04-12", JoinedHoursAgo: 24 * 120, LastSeenMinutesAgo: 2,
	},
	{
		FirstName: "Karim", LastName: "Benali", Email: "karim@demo.dev", Nickname: "karim",
		AboutMe:     "Building things for the web. Occasional baker.",
		DateOfBirth: "1995-11-03", JoinedHoursAgo: 24 * 110, LastSeenMinutesAgo: 14,
	},
	{
		FirstName: "Sofia", LastName: "Marchetti", Email: "sofia@demo.dev", Nickname: "sofia",
		AboutMe:     "Photographer. I shoot film and develop it myself.",
		DateOfBirth: "1996-07-25", JoinedHoursAgo: 24 * 95, LastSeenMinutesAgo: 40,
	},
	{
		FirstName: "Lucas", LastName: "Moreau", Email: "lucas@demo.dev", Nickname: "lucas",
		AboutMe:     "Weekend hiker, weekday sysadmin.",
		DateOfBirth: "1993-02-18", JoinedHoursAgo: 24 * 88, LastSeenMinutesAgo: 75,
	},
	{
		FirstName: "Ines", LastName: "Laurent", Email: "ines@demo.dev", Nickname: "ines",
		AboutMe:     "Landscape architect. Private account, low noise.",
		DateOfBirth: "1997-09-30", JoinedHoursAgo: 24 * 70, LastSeenMinutesAgo: 180,
		IsPrivate: true,
	},
	{
		FirstName: "Omar", LastName: "Diallo", Email: "omar@demo.dev", Nickname: "omar",
		AboutMe:     "Distributed systems and long walks.",
		DateOfBirth: "1994-12-08", JoinedHoursAgo: 24 * 64, LastSeenMinutesAgo: 5,
	},
	{
		FirstName: "Chloe", LastName: "Dubois", Email: "chloe@demo.dev", Nickname: "chloe",
		AboutMe:     "Designer. Type nerd.",
		DateOfBirth: "1999-01-22", JoinedHoursAgo: 24 * 55, LastSeenMinutesAgo: 22,
	},
	{
		FirstName: "Hugo", LastName: "Fontaine", Email: "hugo@demo.dev", Nickname: "hugo",
		AboutMe:     "Keep it to myself. Private profile.",
		DateOfBirth: "1992-06-14", JoinedHoursAgo: 24 * 40, LastSeenMinutesAgo: 600,
		IsPrivate: true,
	},
	{
		FirstName: "Lina", LastName: "Rossi", Email: "lina@demo.dev", Nickname: "lina",
		AboutMe:     "Trail runner and map collector.",
		DateOfBirth: "2000-03-09", JoinedHoursAgo: 24 * 33, LastSeenMinutesAgo: 9,
	},
	{
		FirstName: "Adam", LastName: "Kowalski", Email: "adam@demo.dev", Nickname: "adam",
		AboutMe:     "Board games and mechanical keyboards.",
		DateOfBirth: "1991-08-19", JoinedHoursAgo: 24 * 28, LastSeenMinutesAgo: 55,
	},
	{
		FirstName: "Nora", LastName: "Haddad", Email: "nora@demo.dev", Nickname: "nora",
		AboutMe:     "Just joined. Mostly lurking.",
		DateOfBirth: "2001-10-05", JoinedHoursAgo: 24 * 12, LastSeenMinutesAgo: 3,
	},
	{
		FirstName: "Sami", LastName: "Cherif", Email: "sami@demo.dev", Nickname: "sami",
		AboutMe:     "New here. Private until I know you.",
		DateOfBirth: "2002-05-27", JoinedHoursAgo: 24 * 5, LastSeenMinutesAgo: 1440,
		IsPrivate: true,
	},
}

// Follow status values accepted by the followers table CHECK constraint.
const (
	statusAccepted = "accepted"
	statusPending  = "pending"
)

type seedFollow struct {
	From, To int
	Status   string
	HoursAgo int
}

var seedFollows = []seedFollow{
	{1, 2, statusAccepted, 24 * 110},
	{1, 3, statusAccepted, 24 * 100},
	{1, 4, statusAccepted, 24 * 90},
	{1, 5, statusPending, 6}, // private target, waiting on Ines
	{2, 1, statusAccepted, 24 * 109},
	{2, 3, statusAccepted, 24 * 98},
	{2, 6, statusAccepted, 24 * 60},
	{3, 1, statusAccepted, 24 * 99},
	{3, 2, statusAccepted, 24 * 97},
	{3, 7, statusAccepted, 24 * 50},
	{4, 1, statusAccepted, 24 * 89},
	{4, 5, statusPending, 20}, // second pending request for Ines
	{4, 8, statusPending, 30}, // private target, waiting on Hugo
	{4, 9, statusAccepted, 24 * 31},
	{5, 1, statusAccepted, 24 * 69},
	{5, 9, statusAccepted, 24 * 32},
	{6, 1, statusAccepted, 24 * 58},
	{6, 2, statusAccepted, 24 * 59},
	{6, 5, statusAccepted, 24 * 57},
	{6, 9, statusAccepted, 24 * 31},
	{7, 3, statusAccepted, 24 * 49},
	{7, 1, statusAccepted, 24 * 48},
	{8, 3, statusAccepted, 24 * 39},
	{9, 2, statusAccepted, 24 * 28},
	{9, 4, statusAccepted, 24 * 30},
	{9, 5, statusAccepted, 24 * 29},
	{9, 6, statusAccepted, 24 * 30},
	{9, 10, statusAccepted, 24 * 27},
	{10, 3, statusAccepted, 24 * 23},
	{10, 7, statusAccepted, 24 * 24},
	{10, 9, statusAccepted, 24 * 26},
	{10, 1, statusAccepted, 24 * 25},
	{11, 7, statusAccepted, 24 * 11},
	{12, 10, statusPending, 8}, // private user requesting to follow Adam
}

type seedGroupInvite struct {
	Inviter, Invitee int
	Status           string
	HoursAgo         int
}

type seedGroupRequest struct {
	User     int
	Status   string
	HoursAgo int
}

type seedGroup struct {
	Title       string
	Description string
	Creator     int
	Members     []int
	Invitations []seedGroupInvite
	Requests    []seedGroupRequest
}

var seedGroups = []seedGroup{
	{
		Title:       "Go Developers",
		Description: "Daily drivers, code reviews, and arguments about error handling.",
		Creator:     1,
		Members:     []int{1, 2, 3, 6},
		Invitations: []seedGroupInvite{
			{1, 4, "accepted", 24 * 80},
			{1, 7, "accepted", 24 * 45},
			{2, 9, "pending", 10},
		},
		Requests: []seedGroupRequest{
			{11, "accepted", 24 * 11},
			{12, "pending", 4},
		},
	},
	{
		Title:       "Weekend Hikers",
		Description: "Trails within two hours of the city. Post route plans on Thursdays.",
		Creator:     4,
		Members:     []int{4, 5, 9},
		Invitations: []seedGroupInvite{
			{4, 8, "pending", 26},
		},
		Requests: []seedGroupRequest{
			{7, "pending", 15},
			{2, "rejected", 24 * 20},
		},
	},
	{
		Title:       "Photography Club",
		Description: "Darkroom nights, photo walks, and gear talk. Beginners very welcome.",
		Creator:     3,
		Members:     []int{3, 1, 7, 10},
		Invitations: []seedGroupInvite{
			{3, 11, "pending", 2},
		},
		Requests: []seedGroupRequest{
			{6, "pending", 1},
		},
	},
}

type seedReaction struct {
	User int
	Like bool
}

type seedComment struct {
	Author   int
	Content  string
	HoursAgo int
}

type seedPost struct {
	Author    int
	Group     int // 0 means a profile post, otherwise the group id
	Title     string
	Content   string
	Privacy   string
	HoursAgo  int
	Viewers   []int // followers allowed to see a private post
	Reactions []seedReaction
	Comments  []seedComment
}

// Privacy levels accepted by the posts table CHECK constraint.
const (
	privacyPublic        = "public"
	privacyAlmostPrivate = "almost_private"
	privacyPrivate       = "private"
)

var seedPosts = []seedPost{
	{
		Author: 1, Title: "Finally shipped the rewrite",
		Content: "Three weeks of yak shaving and the service is behind a feature flag. Learned a lot about SQLite write contention along the way.",
		Privacy: privacyPublic, HoursAgo: 3,
		Reactions: []seedReaction{{2, true}, {3, true}, {6, true}, {7, false}, {10, true}},
		Comments: []seedComment{
			{2, "Congrats! What did you end up using for the migrations?", 2},
			{1, "golang-migrate. Boring and it works.", 150},
			{6, "Feature flag or no prod yet?", 120},
		},
	},
	{
		Author: 1, Title: "Coffee recommendations",
		Content: "The usual light roasts are fine but I want something that does not taste like burnt cardboard. Please send recommendations.",
		Privacy: privacyAlmostPrivate, HoursAgo: 26,
		Reactions: []seedReaction{{2, true}, {3, true}, {4, true}, {7, true}},
		Comments: []seedComment{
			{3, "Ethiopian naturals. Floral, no roast flavour.", 24},
			{4, "Seconding. Ask for a washed Yirgacheffe.", 22},
		},
	},
	{
		Author: 1, Title: "Reading list for systems people",
		Content: "Collected the papers and posts that actually changed how I think about concurrency. Happy to share the doc.",
		Privacy: privacyPrivate, HoursAgo: 50, Viewers: []int{2, 3, 6},
		Reactions: []seedReaction{{2, true}, {6, true}},
		Comments:  []seedComment{{2, "This is the good one, thanks.", 48}},
	},
	{
		Author: 2, Title: "My dotfiles are finally public",
		Content: "Ten years of accumulated shell aliases, finally pushed to a repo. The tmux config alone is worth it.",
		Privacy: privacyPublic, HoursAgo: 8,
		Reactions: []seedReaction{{1, true}, {3, true}, {7, true}, {9, true}, {10, false}},
		Comments: []seedComment{
			{1, "Stolen the git alias section, thank you.", 7},
			{7, "What terminal font are you using with this?", 6},
		},
	},
	{
		Author: 2, Title: "Baked too many croissants",
		Content: "Made a batch of forty. Brought six to the office, the rest are spoken for. Come by before they go stale.",
		Privacy: privacyAlmostPrivate, HoursAgo: 30,
		Reactions: []seedReaction{{1, true}, {3, true}, {6, true}, {9, true}},
		Comments:  []seedComment{{9, "Is there any left?", 28}},
	},
	{
		Author: 2, Title: "Notes from the systems reading group",
		Content: "We went through the paper on optimistic concurrency control. The failure cases are more interesting than the happy path.",
		Privacy: privacyPublic, HoursAgo: 96,
		Reactions: []seedReaction{{1, true}, {6, true}, {10, true}},
		Comments:  []seedComment{{6, "The retry storm section is real.", 90}},
	},
	{
		Author: 3, Title: "Roll of Portra 400 developed",
		Content: "Shot a roll over the weekend, developed it last night. Scanned half of it, will post the rest once I finish culling.",
		Privacy: privacyPublic, HoursAgo: 14,
		Reactions: []seedReaction{{1, true}, {2, true}, {4, true}, {7, true}, {9, true}, {10, true}, {11, true}},
		Comments: []seedComment{
			{1, "The colours in the third shot are unreal.", 12},
			{7, "Portra is forgiving and I am not.", 11},
			{4, "Where is this?", 9},
			{3, "Old town, early morning. Fog rolled in around 6.", 8},
		},
	},
	{
		Author: 3, Title: "Darkroom night this Thursday",
		Content: "The club darkroom is booked for me on Thursday evening. Four enlargers, plenty of chemistry, bring a roll.",
		Privacy: privacyAlmostPrivate, HoursAgo: 40,
		Reactions: []seedReaction{{1, true}, {7, true}, {10, true}},
		Comments:  []seedComment{{10, "Count me in, I need to finish a roll.", 38}},
	},
	{
		Author: 4, Title: "Trail conditions this month",
		Content: "The north route is still muddy below the second ridge. Waterproof boots are not optional. Rest of the loop is fine.",
		Privacy: privacyPublic, HoursAgo: 19,
		Reactions: []seedReaction{{5, true}, {9, true}, {1, true}, {2, false}},
		Comments: []seedComment{
			{9, "Was there on Sunday, Lucas is right about the mud.", 17},
			{5, "I can run a water point check on Saturday.", 16},
		},
	},
	{
		Author: 4, Title: "Gear check before Saturday",
		Content: "Listing what I am bringing so nobody forgets a stove. Post anything I missed in the comments.",
		Privacy: privacyAlmostPrivate, HoursAgo: 21,
		Reactions: []seedReaction{{9, true}, {1, true}},
		Comments:  []seedComment{{9, "Spare fuel canister, please.", 20}},
	},
	{
		Author: 5, Title: "Private: site survey photos",
		Content: "Sharing the survey renders with the people on the project only. Do not forward these outside the team.",
		Privacy: privacyPrivate, HoursAgo: 62, Viewers: []int{6, 9},
		Reactions: []seedReaction{{9, true}},
		Comments:  []seedComment{{9, "The south elevation is much better than expected.", 60}},
	},
	{
		Author: 6, Title: "Consensus algorithms, briefly",
		Content: "Wrote up the shortest explanation of Paxos I could manage without lying. It is two pages and it took a week.",
		Privacy: privacyPublic, HoursAgo: 54,
		Reactions: []seedReaction{{1, true}, {2, true}, {10, true}},
		Comments:  []seedComment{{1, "Two pages is the right length for this.", 52}},
	},
	{
		Author: 7, Title: "Kerning is not a personality trait",
		Content: "Unpopular opinion incoming. I care about optical sizing, not your tracking. Fight me.",
		Privacy: privacyPublic, HoursAgo: 11,
		Reactions: []seedReaction{{3, true}, {1, true}, {10, false}, {2, false}},
		Comments: []seedComment{
			{3, "This is the most Chloe post ever.", 10},
			{10, "The tracking is fine, relax.", 9},
		},
	},
	{
		Author: 9, Title: "New route recommendation",
		Content: "Found a nine kilometre loop that avoids the crowded ridge. Starts at the same trailhead as the old one.",
		Privacy: privacyAlmostPrivate, HoursAgo: 34,
		Reactions: []seedReaction{{4, true}, {5, true}, {6, true}},
		Comments:  []seedComment{{4, "Adding it to the group post for Saturday.", 33}},
	},
	{
		Author: 10, Title: "Board game night, all welcome",
		Content: "Every other Thursday at mine. I have three games that need at least four players and not enough friends.",
		Privacy: privacyPublic, HoursAgo: 72,
		Reactions: []seedReaction{{7, true}, {3, true}, {1, true}, {9, true}},
		Comments:  []seedComment{{7, "I am in if the date works.", 70}},
	},
	{
		Author: 11, Title: "Hello",
		Content: "Just joined. I mostly lurk, but I will try to post something eventually.",
		Privacy: privacyPublic, HoursAgo: 11,
		Reactions: []seedReaction{{3, true}, {7, true}},
		Comments:  []seedComment{{7, "Welcome aboard.", 10}},
	},
	// Group posts: visible to members of the owning group.
	{
		Author: 1, Group: 1, Title: "Code review thread: error wrapping",
		Content: "Started a doc on when to wrap versus when to add context. I want one rule that is not hand-wavy.",
		Privacy: privacyPublic, HoursAgo: 15,
		Reactions: []seedReaction{{2, true}, {3, true}, {6, true}},
		Comments: []seedComment{
			{2, `There is no single rule, but the "add context at the boundary" heuristic holds up.`, 14},
			{6, "I disagree, but I wrote a doc about it.", 13},
			{1, "Post the link, I will read it this week.", 12},
		},
	},
	{
		Author: 2, Group: 1, Title: "SQLite pragmas worth knowing",
		Content: "Short list: enable WAL, set a busy timeout, and never share a connection across goroutines without a mutex. Learned these the hard way.",
		Privacy: privacyPublic, HoursAgo: 28,
		Reactions: []seedReaction{{1, true}, {6, true}, {3, true}},
		Comments:  []seedComment{{1, "The busy timeout one bit me on the events table.", 27}},
	},
	{
		Author: 6, Group: 1, Title: "Anyone up for a pairing session?",
		Content: "Stuck on a deadlock that only shows up under load. Looking for someone to pair on it, an hour is plenty.",
		Privacy: privacyPublic, HoursAgo: 5,
		Reactions: []seedReaction{{1, true}, {2, true}},
		Comments:  []seedComment{{1, "I am free this afternoon.", 4}},
	},
	{
		Author: 4, Group: 2, Title: "Saturday route plan",
		Content: "Leaving at eight from the trailhead parking. Loop is about fourteen kilometres with six hundred metres of ascent. Water points confirmed.",
		Privacy: privacyPublic, HoursAgo: 18,
		Reactions: []seedReaction{{5, true}, {9, true}},
		Comments: []seedComment{
			{5, "I will bring the second first aid kit.", 17},
			{9, "Count me in.", 16},
		},
	},
	{
		Author: 9, Group: 2, Title: "Alternative route if it rains",
		Content: "If Saturday is wet we should switch to the sheltered loop by the reservoir. Shorter, less interesting, but nobody gets rained on.",
		Privacy: privacyPublic, HoursAgo: 9,
		Reactions: []seedReaction{{4, true}, {5, true}},
		Comments:  []seedComment{{4, "Good plan, I will check the forecast on Friday night.", 8}},
	},
	{
		Author: 3, Group: 3, Title: "Photo walk: old town, Sunday morning",
		Content: "Meeting at the main square at seven. We will work through the alleys and finish by the bridge. Bring a wide lens if you have one.",
		Privacy: privacyPublic, HoursAgo: 20,
		Reactions: []seedReaction{{1, true}, {7, true}, {10, true}},
		Comments:  []seedComment{{7, "Early, but I am in.", 19}},
	},
	{
		Author: 7, Group: 3, Title: "Composition notes from the last walk",
		Content: "Wrote up what worked and what did not, including the shots where I framed badly. Photos included, gently roasted.",
		Privacy: privacyPublic, HoursAgo: 6,
		Reactions: []seedReaction{{3, true}, {1, true}, {10, true}},
		Comments:  []seedComment{{3, "The second-to-last one is the best of the set.", 5}},
	},
}

type seedVote struct {
	User      int
	OptionIdx int
}

type seedEvent struct {
	Group           int
	Creator         int
	Title           string
	Description     string
	InDays          int
	At              string // "15:04" format, matching the create handler
	Options         [2]string
	CreatedHoursAgo int
	Votes           []seedVote
}

var seedEvents = []seedEvent{
	{
		Group: 1, Creator: 1, Title: "Go meetup: talks + pizza",
		Description: "Three lightning talks and then pizza. Sign up so we order the right amount.",
		InDays:      9, At: "18:30", Options: [2]string{"going", "not going"},
		CreatedHoursAgo: 24 * 4,
		Votes:           []seedVote{{2, 0}, {3, 0}, {6, 1}},
	},
	{
		Group: 2, Creator: 4, Title: "Saturday hike",
		Description: "The fourteen kilometre loop. Bring boots and water. Route plan is pinned in the group.",
		InDays:      3, At: "08:00", Options: [2]string{"going", "not going"},
		CreatedHoursAgo: 24 * 2,
		Votes:           []seedVote{{5, 0}, {9, 0}},
	},
	{
		Group: 3, Creator: 3, Title: "Darkroom night",
		Description: "Four enlargers booked. Bring a roll you have already shot.",
		InDays:      6, At: "19:00", Options: [2]string{"going", "not going"},
		CreatedHoursAgo: 24 * 6,
		Votes:           []seedVote{{1, 0}, {7, 0}, {10, 1}},
	},
}

// seedMessage is a private DM. To must have an accepted follow with From in at
// least one direction, which is what ChatRepository.CanSendDM requires.
type seedMessage struct {
	From, To   int
	Content    string
	MinutesAgo int
}

var seedMessages = []seedMessage{
	{1, 2, "Did you see the CI failure on the migration step?", 240},
	{2, 1, "Yeah, it was the dirty flag. I pushed a fix.", 236},
	{1, 2, "Nice. I will pull once the build finishes.", 235},
	{2, 1, "There are croissants left if you are quick.", 40},
	{1, 2, "Already took two. Save me one.", 38},
	{1, 3, "The scans from the roll look incredible.", 200},
	{3, 1, "Thank you! The fog helped.", 198},
	{3, 1, "Any interest in a photo walk next Sunday?", 190},
	{1, 3, "Count me in.", 186},
	{2, 6, "Are you around for pairing later?", 150},
	{6, 2, "Yes, after four. Same place as last time?", 148},
	{2, 6, "Works for me.", 147},
	{3, 7, "The composition notes are great. Roast me properly next time.", 120},
	{7, 3, "That was the point.", 118},
	{4, 9, "Is the new loop worth it?", 95},
	{9, 4, "Very. I will post the GPX tonight.", 93},
	{4, 1, "Thanks for pinning the route plan.", 60},
	{1, 4, "Anytime. Bring the big stove.", 58},
	{7, 10, "Thursday works for me.", 300},
	{10, 7, "Excellent, I will dig out the party game.", 295},
	{1, 7, "Welcome to the club by the way, Chloe.", 500},
	{7, 1, "Thanks Yasmine.", 495},
}

// seedGroupMessage is a message in a group's shared chat room.
type seedGroupMessage struct {
	Group      int
	From       int
	Content    string
	MinutesAgo int
}

var seedGroupMessages = []seedGroupMessage{
	{1, 1, "Pushed the error wrapping draft, comments welcome.", 200},
	{1, 2, "Reading it now.", 198},
	{1, 2, "The boundary heuristic is the only one I trust.", 190},
	{1, 6, "I wrote the opposite opinion into a doc, will share later.", 150},
	{1, 1, "Please do, I want the disagreement on the record.", 148},
	{1, 3, "Anyone up for lunch after the meetup?", 45},
	{2, 4, "Route plan is up. Forecast still looks dry.", 180},
	{2, 5, "I will bring the first aid kit.", 176},
	{2, 9, "Adding the sheltered loop as a backup.", 130},
	{2, 4, "Good thinking. Decision on Friday night.", 128},
	{3, 3, "Photo walk Sunday, seven at the main square.", 240},
	{3, 7, "Bringing a wide lens, might borrow two spares.", 236},
	{3, 10, "I will be fifteen minutes late, shoot without me.", 120},
	{3, 1, "No problem, we will start with the alleys.", 118},
}

// Notification types emitted by the backend.
const (
	notifFollowRequest = "follow_request"
	notifNewFollower   = "new_follower"
	notifGroupInvite   = "group_invitation"
	notifJoinRequest   = "join_request"
	notifEventCreated  = "event_created"
)

type seedNotification struct {
	User       int
	Actor      int
	Type       string
	RefID      int
	Message    string
	IsRead     bool
	MinutesAgo int
}

// acceptedFollowAlerts are the accepted follows that produced a new_follower
// notification back when they were accepted, so those rows start out read.
var acceptedFollowAlerts = []struct {
	From, To   int
	MinutesAgo int
}{
	{1, 2, 24 * 109},
	{1, 3, 24 * 99},
	{2, 1, 24 * 96},
	{3, 2, 24 * 88},
	{4, 1, 24 * 88},
}
