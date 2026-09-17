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
	limiter := middleware.NewRateLimiter(100, time.Minute)

	// auth routes
	mux.HandleFunc("/register", auth.RegisterHandler)
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.HandleFunc("/logout", middleware.Auth(auth.LogoutHandler))
	mux.HandleFunc("/me", middleware.Auth(auth.MeHandler))

	// posts routes
	mux.HandleFunc("GET /posts", middleware.Auth(posts.GetPostsAPI))
	mux.HandleFunc("GET /posts/{id}", middleware.Auth(posts.GetPostByIDAPI))
	mux.HandleFunc("POST /posts/create", middleware.Auth(posts.CreatePostAPI))

	// comments routes
	mux.HandleFunc("GET /posts/{id}/comment", middleware.Auth(comments.GetCommentsAPI))
	mux.HandleFunc("POST /posts/{id}/comment", middleware.Auth(comments.CreateCommentAPI))
	mux.HandleFunc("POST /posts/react", middleware.Auth(reactions.ReactToPost))

	// profile and users routes
	mux.HandleFunc("GET /profile/{id}", middleware.Auth(profile.GetProfileAPI))
	mux.HandleFunc("POST /profile/edit", middleware.Auth(profile.EditProfileAPI))
	mux.HandleFunc("POST /profile/toggle-visibility", middleware.Auth(profile.SwitchStatus))
	mux.HandleFunc("GET /users/search", middleware.Auth(users.SearchUsersAPI))

	// followers routes
	mux.HandleFunc("GET /profile/{id}/followers", middleware.Auth(followers.GetFollowersAPI))
	mux.HandleFunc("GET /profile/{id}/following", middleware.Auth(followers.GetFollowingAPI))
	mux.HandleFunc("POST /follow/{id}", middleware.Auth(followers.ToggleFollowAPI))
	mux.HandleFunc("POST /follow/{id}/unfollow", middleware.Auth(followers.UnfollowAPI))
	mux.HandleFunc("GET /follow/requests", middleware.Auth(followers.GetRequestsAPI))
	mux.HandleFunc("POST /follow/{id}/{action}", middleware.Auth(followers.RespondRequestAPI))

	// media & uploads routes
	mux.HandleFunc("/uploads/posts/", middleware.Auth(uploads.ServeImage))
	mux.HandleFunc("/uploads/comments/", middleware.Auth(uploads.ServeImage))
	mux.HandleFunc("/uploads/avatars/", middleware.Auth(utils.ServeUpload))
	// groups routes
	mux.HandleFunc("POST /groups", middleware.Auth(handlers.CreatGroupeHandler))
	mux.HandleFunc("GET /groups", middleware.Auth(handlers.GetAllGroupsHandler))
	mux.HandleFunc("POST /groups/{group_id}", middleware.Auth(handlers.GetGroupInfo))
	mux.HandleFunc("GET /groups/{group_id}/posts", middleware.Auth(posts.GetPostsByGroupAPI))
	mux.HandleFunc("GET /groups/{group_id}/comments", middleware.Auth(comments.GetCommentsByGroupAPI))

	// join request
	mux.HandleFunc("POST /groups/{group_id}/join", middleware.Auth(handlers.Join_Request_handler))
	mux.HandleFunc("POST /join_requests/{request_id}/accept", middleware.Auth(handlers.Accept_Join_Request))
	mux.HandleFunc("POST /join_requests/{request_id}/reject", middleware.Auth(handlers.Reject_Join_Request))

	// invitation
	mux.HandleFunc("POST /groups/{group_id}/invitations", middleware.Auth(handlers.InvitationsHandler))
	mux.HandleFunc("POST /invitations/{invitation_id}/accept", middleware.Auth(handlers.AcceptInvitationResponse))
	mux.HandleFunc("POST /invitations/{invitation_id}/reject", middleware.Auth(handlers.RejectInvitationResponse))
	mux.HandleFunc("GET /groups/{group_id}/users", middleware.Auth(handlers.GetAllUsersNotmemberHandler))


	// events routes
	mux.HandleFunc("POST /events", middleware.Auth(events.CreateEvent))
	mux.HandleFunc("GET /events", middleware.Auth(events.GetEvents))
	mux.HandleFunc("POST /events/vote", middleware.Auth(events.VoteEvent))

	// notifications
	mux.HandleFunc("GET /notifications", middleware.Auth(handlers.NotifcationHandler))

	// websocket routes
	mux.HandleFunc("GET /ws", middleware.Auth(wshandlers.WSHandler(websocket.DefaultHub)))

	// global middleware
	return (limiter.Middleware(mux.ServeHTTP))
}
