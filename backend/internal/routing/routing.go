package routing

import (
	"net/http"
	"time"

	"zone/internal/auth"
	"zone/internal/comments"
	"zone/internal/followers"
	"zone/internal/groups/handlers"
	"zone/internal/middleware"
	"zone/internal/posts"
	"zone/internal/profile"
	"zone/internal/reactions"
	"zone/internal/users"
	"zone/internal/utils"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	limiter := middleware.NewRateLimiter(100, time.Minute)

	mux.HandleFunc("/register", auth.Register)
	mux.HandleFunc("/login", auth.Login)
	mux.HandleFunc("/me", auth.Me)
	mux.HandleFunc("GET /users/search", users.SearchUsersAPI)
	mux.HandleFunc("GET /posts", posts.GetPostsAPI)
	mux.HandleFunc("GET /profile/{id}", profile.GetProfileAPI)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	mux.HandleFunc("/logout", middleware.Auth(auth.Logout))
	mux.HandleFunc("/upload/avatar", utils.UploadAvatar)
	mux.HandleFunc("POST /posts/create", middleware.Auth(posts.CreatePostAPI))
	mux.HandleFunc("POST /profile/toggle-visibility", middleware.Auth(profile.ToggleVisibilityAPI))
	mux.HandleFunc("POST /posts/react", middleware.Auth(reactions.ReactToPost))
	mux.HandleFunc("GET /posts/{id}/comment", comments.GetCommentsAPI)
	mux.HandleFunc("GET /posts/{id}", posts.GetPostByIDAPI)
	mux.HandleFunc("POST /posts/{id}/comment", middleware.Auth(comments.CreateCommentAPI))
	mux.HandleFunc("GET /follow/requests", middleware.Auth(followers.GetPendingRequestsAPI))
	mux.HandleFunc("POST /follow/{id}", middleware.Auth(followers.FollowAPI))
	mux.HandleFunc("POST /follow/{id}/unfollow", middleware.Auth(followers.UnfollowAPI))
	mux.HandleFunc("POST /follow/{id}/accept", middleware.Auth(followers.AcceptFollowAPI))
	mux.HandleFunc("POST /follow/{id}/decline", middleware.Auth(followers.DeclineFollowAPI))
	mux.HandleFunc("GET /follow/{id}/status", followers.GetFollowStatusAPI)
	mux.HandleFunc("GET /profile/{id}/followers", followers.GetFollowersAPI)
	mux.HandleFunc("GET /profile/{id}/following", followers.GetFollowingAPI)

	// Groups
	mux.HandleFunc("GET /groups", middleware.Auth(handlers.BrowseGroupsHandler))
	mux.HandleFunc("POST /groups/create", middleware.Auth(handlers.CreateGroupHandler))
	mux.HandleFunc("GET /groups/{id}", middleware.Auth(handlers.GetGroupHandler))
	mux.HandleFunc("GET /groups/{id}/members", middleware.Auth(handlers.ListMembersHandler))

	// Invitations
	mux.HandleFunc("POST /groups/{id}/invite", middleware.Auth(handlers.InviteUserHandler))
	mux.HandleFunc("GET /invitations", middleware.Auth(handlers.MyInvitationsHandler))
	mux.HandleFunc("POST /invitations/{id}/{action}", middleware.Auth(handlers.RespondInvitationHandler))

	// Join requests
	mux.HandleFunc("POST /groups/{id}/request", middleware.Auth(handlers.RequestToJoinHandler))
	mux.HandleFunc("GET /groups/{id}/requests", middleware.Auth(handlers.GroupRequestsHandler))
	mux.HandleFunc("POST /groups/{id}/requests/{requestID}/{action}", middleware.Auth(handlers.RespondJoinRequestHandler))

	// Group posts & comments
	mux.HandleFunc("GET /groups/{id}/posts", middleware.Auth(handlers.GroupPostsHandler))
	mux.HandleFunc("POST /groups/{id}/posts", middleware.Auth(handlers.CreateGroupPostHandler))
	mux.HandleFunc("GET /groups/posts/{postID}/comments", middleware.Auth(handlers.GroupCommentsHandler))
	mux.HandleFunc("POST /groups/posts/{postID}/comments", middleware.Auth(handlers.CreateGroupCommentHandler))

	// Events
	mux.HandleFunc("GET /groups/{id}/events", middleware.Auth(handlers.EventsListHandler))
	mux.HandleFunc("POST /groups/{id}/events", middleware.Auth(handlers.CreateEventHandler))
	mux.HandleFunc("POST /events/{eventID}/respond", middleware.Auth(handlers.RespondEventHandler))

	return middleware.CORS(limiter.Middleware(mux.ServeHTTP))
}
