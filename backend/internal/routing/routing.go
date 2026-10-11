package routing

import (
	"net/http"
	"strings"
	"time"

	"zone/internal/auth"
	chathandlers "zone/internal/chat/handlers"
	commenthandlers "zone/internal/comments/handlers"
	eventhandlers "zone/internal/events/handlers"
	"zone/internal/followers"
	"zone/internal/groups/handlers"
	"zone/internal/middleware"
	notifhandlers "zone/internal/notifications/handlers"
	posthandlers "zone/internal/posts/handlers"
	"zone/internal/profile"
	reactionhandlers "zone/internal/reactions/handlers"
	"zone/internal/reels"
	"zone/internal/stories"
	"zone/internal/uploads"
	"zone/internal/users"
	"zone/internal/utils"
	"zone/internal/websocket"
	wshandlers "zone/internal/websocket/handlers"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Per-endpoint rate limiters: strict on auth, low on mutations, higher on reads,
	// generous for the websocket handshake. All limiters key by client IP.
	limitAuth := middleware.NewRateLimiter(30, time.Minute)
	limitAction := middleware.NewRateLimiter(60, time.Minute)
	limitRead := middleware.NewRateLimiter(120, time.Minute)
	limitWS := middleware.NewRateLimiter(300, time.Minute)
	limitChatSend := middleware.NewRateLimiter(300, time.Minute)

	// guarded applies Auth (handlers read the authenticated user id) plus the IP-keyed limiter.
	guarded := func(limiter *middleware.RateLimiter, h http.HandlerFunc) http.HandlerFunc {
		return middleware.Auth(limiter.Middleware(h))
	}

	// auth routes
	mux.HandleFunc("/register", limitAuth.Middleware(auth.RegisterHandler))
	mux.HandleFunc("/login", limitAuth.Middleware(auth.LoginHandler))
	mux.HandleFunc("/logout", guarded(limitAction, auth.LogoutHandler))
	mux.HandleFunc("/me", guarded(limitRead, auth.MeHandler))

	// posts routes
	mux.HandleFunc("GET /posts", guarded(limitRead, posthandlers.GetPostsAPI))
	mux.HandleFunc("GET /posts/{id}", guarded(limitRead, posthandlers.GetPostByIDAPI))
	mux.HandleFunc("POST /posts/create", guarded(limitAction, posthandlers.CreatePostAPI))

	// comments routes
	mux.HandleFunc("GET /posts/{id}/comment", guarded(limitRead, commenthandlers.GetCommentsAPI))
	mux.HandleFunc("POST /posts/{id}/comment", guarded(limitAction, commenthandlers.CreateCommentAPI))
	mux.HandleFunc("POST /posts/react", guarded(limitAction, reactionhandlers.ReactToPost))

	// profile and users routes
	mux.HandleFunc("GET /profile/{id}", guarded(limitRead, profile.GetProfileAPI))
	mux.HandleFunc("POST /profile/edit", guarded(limitAction, profile.EditProfileAPI))
	mux.HandleFunc("GET /users/search", guarded(limitRead, users.SearchUsersAPI))
	mux.HandleFunc("GET /users/online", guarded(limitRead, users.OnlineUsersAPI))

	// followers routes
	mux.HandleFunc("GET /profile/{id}/followers", guarded(limitRead, followers.GetFollowersAPI))
	mux.HandleFunc("GET /profile/{id}/following", guarded(limitRead, followers.GetFollowingAPI))
	mux.HandleFunc("POST /follow/{id}", guarded(limitAction, followers.ToggleFollowAPI))
	mux.HandleFunc("POST /follow/{id}/{action}", guarded(limitAction, followers.RespondRequestAPI))

	// stories routes
	mux.HandleFunc("GET /stories", guarded(limitRead, stories.GetStoriesAPI))
	mux.HandleFunc("POST /stories/create", guarded(limitAction, stories.CreateStoryAPI))
	mux.HandleFunc("POST /stories/{id}/view", guarded(limitAction, stories.StoryViewAPI))
	mux.HandleFunc("DELETE /stories/{id}", guarded(limitAction, stories.DeleteStoryAPI))

	// reels routes
	mux.HandleFunc("GET /reels", guarded(limitRead, reels.GetReelsAPI))
	mux.HandleFunc("GET /reels/saved", guarded(limitRead, reels.GetSavedReelsAPI))
	mux.HandleFunc("GET /reels/{id}", guarded(limitRead, reels.GetReelAPI))
	mux.HandleFunc("POST /reels/create", guarded(limitAction, reels.CreateReelAPI))
	mux.HandleFunc("DELETE /reels/{id}", guarded(limitAction, reels.DeleteReelAPI))
	mux.HandleFunc("POST /reels/{id}/view", guarded(limitAction, reels.ReelViewAPI))
	mux.HandleFunc("POST /reels/{id}/like", guarded(limitAction, reels.ToggleReelLikeAPI))
	mux.HandleFunc("POST /reels/{id}/save", guarded(limitAction, reels.ToggleReelSaveAPI))
	mux.HandleFunc("POST /reels/{id}/share", guarded(limitAction, reels.ShareReelAPI))
	mux.HandleFunc("GET /reels/{id}/comments", guarded(limitRead, reels.GetReelCommentsAPI))
	mux.HandleFunc("POST /reels/{id}/comments", guarded(limitAction, reels.CreateReelCommentAPI))

	// media & uploads routes
	mux.HandleFunc("/uploads/posts/", guarded(limitRead, uploads.ServeImage))
	mux.HandleFunc("/uploads/stories/", guarded(limitRead, uploads.ServeImage))
	mux.HandleFunc("/uploads/comments/", guarded(limitRead, uploads.ServeImage))
	// Reels are visible to every signed-in user, so they only need the auth guard;
	// ServeUpload serves them with http.ServeFile, which handles Range requests
	// and therefore lets the player seek without buffering the whole video.
	mux.HandleFunc("/uploads/reels/", guarded(limitRead, utils.ServeUpload))
	mux.HandleFunc("/uploads/avatars/", guarded(limitRead, utils.ServeUpload))
	mux.HandleFunc("/uploads/groups/", guarded(limitRead, uploads.ServeImage))
	// groups routes
	mux.HandleFunc("POST /groups", guarded(limitAction, handlers.CreatGroupeHandler))
	mux.HandleFunc("GET /groups", guarded(limitRead, handlers.GetAllGroupsHandler))
	mux.HandleFunc("POST /groups/{group_id}", guarded(limitAction, handlers.GetGroupInfo))
	mux.HandleFunc("GET /groups/{group_id}/posts", guarded(limitRead, posthandlers.GetPostsByGroupAPI))

	// join request
	mux.HandleFunc("POST /groups/{group_id}/join", guarded(limitAction, handlers.Join_Request_handler))
	mux.HandleFunc("DELETE /groups/{group_id}/join", guarded(limitAction, handlers.CancelJoinRequestHandler))
	mux.HandleFunc("POST /join_requests/{request_id}/accept", guarded(limitAction, handlers.Accept_Join_Request))
	mux.HandleFunc("POST /join_requests/{request_id}/reject", guarded(limitAction, handlers.Reject_Join_Request))

	// invitation
	mux.HandleFunc("POST /groups/{group_id}/invitations", guarded(limitAction, handlers.InvitationsHandler))
	mux.HandleFunc("POST /invitations/{invitation_id}/accept", guarded(limitAction, handlers.AcceptInvitationResponse))
	mux.HandleFunc("POST /invitations/{invitation_id}/reject", guarded(limitAction, handlers.RejectInvitationResponse))
	mux.HandleFunc("GET /groups/{group_id}/users", guarded(limitRead, handlers.GetAllUsersNotmemberHandler))
	mux.HandleFunc("GET /groups/{group_id}/members", guarded(limitRead, handlers.GetGroupMembersHandler))

	// events routes
	mux.HandleFunc("POST /events", guarded(limitAction, eventhandlers.CreateEvent))
	mux.HandleFunc("GET /events", guarded(limitRead, eventhandlers.GetEvents))
	mux.HandleFunc("POST /events/vote", guarded(limitAction, eventhandlers.VoteEvent))

	// notifications
	mux.HandleFunc("GET /notifications", guarded(limitRead, notifhandlers.DefaultHandler.GetNotifications))
	mux.HandleFunc("POST /notifications/{id}/read", guarded(limitAction, notifhandlers.DefaultHandler.MarkReadByID))

	// chat routes
	mux.HandleFunc("POST /messages", guarded(limitChatSend, chathandlers.DefaultHandler.SendDMHandler))
	mux.HandleFunc("GET /messages/{user_id}", guarded(limitRead, chathandlers.DefaultHandler.DMHistoryHandler))
	mux.HandleFunc("POST /groups/{group_id}/messages", guarded(limitChatSend, chathandlers.DefaultHandler.SendGroupHandler))
	mux.HandleFunc("GET /groups/{group_id}/messages", guarded(limitRead, chathandlers.DefaultHandler.GroupHistoryHandler))
	mux.HandleFunc("GET /conversations", guarded(limitRead, chathandlers.DefaultHandler.ConversationsHandler))

	// websocket routes
	mux.HandleFunc("GET /ws", guarded(limitWS, wshandlers.WSHandler(websocket.DefaultHub)))

	// Unauthenticated liveness probe: docker-compose gates Caddy on it, and it
	// doubles as the quick "is the stack up?" curl.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// CORS wraps the whole chain so preflights are answered before Auth can reject
	// them, and so every route carries the headers browser clients need.
	return middleware.CORS(tolerateAPIPrefix(mux))
}

// tolerateAPIPrefix lets one router serve both forms of every path:
//
//	/login      what a proxy that strips /api (our Caddyfile) forwards
//	/api/login  what the mobile app and the browser request directly
//
// Without this, running the server on its own (no Caddy in front) would 404
// every client call, since clients always append /api.
func tolerateAPIPrefix(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trimmed, found := strings.CutPrefix(r.URL.Path, "/api/")
		if !found {
			next.ServeHTTP(w, r)
			return
		}

		// Cloning keeps the mutation local to this request.
		r = r.Clone(r.Context())
		r.URL.Path = "/" + trimmed
		if r.URL.RawPath != "" {
			r.URL.RawPath = "/" + strings.TrimPrefix(r.URL.RawPath, "/api/")
		}
		next.ServeHTTP(w, r)
	})
}
