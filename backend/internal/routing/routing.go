package routing

import (
	"net/http"
	"time"

	"zone/internal/auth"
	"zone/internal/comments"
	"zone/internal/events"
	"zone/internal/followers"
	"zone/internal/groups/handlers"
	"zone/internal/middleware"
	notifhandlers "zone/internal/notifications/handlers"
	"zone/internal/posts"
	"zone/internal/profile"
	"zone/internal/reactions"
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
	mux.HandleFunc("GET /posts", guarded(limitRead, posts.GetPostsAPI))
	mux.HandleFunc("GET /posts/{id}", guarded(limitRead, posts.GetPostByIDAPI))
	mux.HandleFunc("POST /posts/create", guarded(limitAction, posts.CreatePostAPI))

	// comments routes
	mux.HandleFunc("GET /posts/{id}/comment", guarded(limitRead, comments.GetCommentsAPI))
	mux.HandleFunc("POST /posts/{id}/comment", guarded(limitAction, comments.CreateCommentAPI))
	mux.HandleFunc("POST /posts/react", guarded(limitAction, reactions.ReactToPost))

	// profile and users routes
	mux.HandleFunc("GET /profile/{id}", guarded(limitRead, profile.GetProfileAPI))
	mux.HandleFunc("POST /profile/edit", guarded(limitAction, profile.EditProfileAPI))
	mux.HandleFunc("POST /profile/toggle-visibility", guarded(limitAction, profile.SwitchStatus))
	mux.HandleFunc("GET /users/search", guarded(limitRead, users.SearchUsersAPI))

	// followers routes
	mux.HandleFunc("GET /profile/{id}/followers", guarded(limitRead, followers.GetFollowersAPI))
	mux.HandleFunc("GET /profile/{id}/following", guarded(limitRead, followers.GetFollowingAPI))
	mux.HandleFunc("POST /follow/{id}", guarded(limitAction, followers.ToggleFollowAPI))
	mux.HandleFunc("POST /follow/{id}/unfollow", guarded(limitAction, followers.UnfollowAPI))
	mux.HandleFunc("GET /follow/requests", guarded(limitRead, followers.GetRequestsAPI))
	mux.HandleFunc("POST /follow/{id}/{action}", guarded(limitAction, followers.RespondRequestAPI))

	// media & uploads routes
	mux.HandleFunc("/uploads/posts/", guarded(limitRead, uploads.ServeImage))
	mux.HandleFunc("/uploads/comments/", guarded(limitRead, uploads.ServeImage))
	mux.HandleFunc("/uploads/avatars/", guarded(limitRead, utils.ServeUpload))
	// groups routes
	mux.HandleFunc("POST /groups", guarded(limitAction, handlers.CreatGroupeHandler))
	mux.HandleFunc("GET /groups", guarded(limitRead, handlers.GetAllGroupsHandler))
	mux.HandleFunc("POST /groups/{group_id}", guarded(limitAction, handlers.GetGroupInfo))
	mux.HandleFunc("GET /groups/{group_id}/posts", guarded(limitRead, posts.GetPostsByGroupAPI))

	// join request
	mux.HandleFunc("POST /groups/{group_id}/join", guarded(limitAction, handlers.Join_Request_handler))
	mux.HandleFunc("POST /join_requests/{request_id}/accept", guarded(limitAction, handlers.Accept_Join_Request))
	mux.HandleFunc("POST /join_requests/{request_id}/reject", guarded(limitAction, handlers.Reject_Join_Request))

	// invitation
mux.HandleFunc("POST /groups/{group_id}/invitations", guarded(limitAction, handlers.InvitationsHandler))
	mux.HandleFunc("POST /invitations/{invitation_id}/accept", guarded(limitAction, handlers.AcceptInvitationResponse))
	mux.HandleFunc("POST /invitations/{invitation_id}/reject", guarded(limitAction, handlers.RejectInvitationResponse))
	mux.HandleFunc("GET /groups/{group_id}/users", guarded(limitRead, handlers.GetAllUsersNotmemberHandler))
	mux.HandleFunc("GET /groups/{group_id}/members", guarded(limitRead, handlers.GetGroupMembersHandler))

	// events routes
	mux.HandleFunc("POST /events", guarded(limitAction, events.CreateEvent))
	mux.HandleFunc("GET /events", guarded(limitRead, events.GetEvents))
	mux.HandleFunc("POST /events/vote", guarded(limitAction, events.VoteEvent))

	// notifications
	mux.HandleFunc("GET /notifications", guarded(limitRead, notifhandlers.DefaultHandler.GetNotifications))
	mux.HandleFunc("POST /notifications/{id}/read", guarded(limitAction, notifhandlers.DefaultHandler.MarkReadByID))

	// websocket routes
	mux.HandleFunc("GET /ws", guarded(limitWS, wshandlers.WSHandler(websocket.DefaultHub)))

	return mux
}
