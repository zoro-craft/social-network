package routing

import (
	"net/http"

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
	"zone/internal/uploads"
	"zone/internal/users"
	"zone/internal/utils"
	"zone/internal/websocket"
	wshandlers "zone/internal/websocket/handlers"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// auth routes
	mux.HandleFunc("/register", auth.RegisterHandler)
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.HandleFunc("/logout", middleware.Auth(auth.LogoutHandler))
	mux.HandleFunc("/me", middleware.Auth(auth.MeHandler))

	// posts routes
	mux.HandleFunc("GET /posts", middleware.Auth(posthandlers.GetPostsAPI))
	mux.HandleFunc("GET /posts/{id}", middleware.Auth(posthandlers.GetPostByIDAPI))
	mux.HandleFunc("POST /posts/create", middleware.Auth(posthandlers.CreatePostAPI))

	// comments routes
	mux.HandleFunc("GET /posts/{id}/comment", middleware.Auth(commenthandlers.GetCommentsAPI))
	mux.HandleFunc("POST /posts/{id}/comment", middleware.Auth(commenthandlers.CreateCommentAPI))
	mux.HandleFunc("POST /posts/react", middleware.Auth(reactionhandlers.ReactToPost))

	// profile and users routes
	mux.HandleFunc("GET /profile/{id}", middleware.Auth(profile.GetProfileAPI))
	mux.HandleFunc("POST /profile/edit", middleware.Auth(profile.EditProfileAPI))
	mux.HandleFunc("GET /users/search", middleware.Auth(users.SearchUsersAPI))

	// followers routes
	mux.HandleFunc("GET /profile/{id}/followers", middleware.Auth(followers.GetFollowersAPI))
	mux.HandleFunc("GET /profile/{id}/following", middleware.Auth(followers.GetFollowingAPI))
	mux.HandleFunc("POST /follow/{id}", middleware.Auth(followers.ToggleFollowAPI))
	mux.HandleFunc("POST /follow/{id}/{action}", middleware.Auth(followers.RespondRequestAPI))

	// media & uploads routes
	mux.HandleFunc("/uploads/posts/", middleware.Auth(uploads.ServeImage))
	mux.HandleFunc("/uploads/comments/", middleware.Auth(uploads.ServeImage))
	mux.HandleFunc("/uploads/avatars/", middleware.Auth(utils.ServeUpload))
	// groups routes
	mux.HandleFunc("POST /groups", middleware.Auth(handlers.CreatGroupeHandler))
	mux.HandleFunc("GET /groups", middleware.Auth(handlers.GetAllGroupsHandler))
	mux.HandleFunc("POST /groups/{group_id}", middleware.Auth(handlers.GetGroupInfo))
	mux.HandleFunc("GET /groups/{group_id}/posts", middleware.Auth(posthandlers.GetPostsByGroupAPI))

	// join request
	mux.HandleFunc("POST /groups/{group_id}/join", middleware.Auth(handlers.Join_Request_handler))
	mux.HandleFunc("POST /join_requests/{request_id}/accept", middleware.Auth(handlers.Accept_Join_Request))
	mux.HandleFunc("POST /join_requests/{request_id}/reject", middleware.Auth(handlers.Reject_Join_Request))

	// invitation
	mux.HandleFunc("POST /groups/{group_id}/invitations", middleware.Auth(handlers.InvitationsHandler))
	mux.HandleFunc("POST /invitations/{invitation_id}/accept", middleware.Auth(handlers.AcceptInvitationResponse))
	mux.HandleFunc("POST /invitations/{invitation_id}/reject", middleware.Auth(handlers.RejectInvitationResponse))
	mux.HandleFunc("GET /groups/{group_id}/users", middleware.Auth(handlers.GetAllUsersNotmemberHandler))
	mux.HandleFunc("GET /groups/{group_id}/members", middleware.Auth(handlers.GetGroupMembersHandler))

	// events routes
	mux.HandleFunc("POST /events", middleware.Auth(eventhandlers.CreateEvent))
	mux.HandleFunc("GET /events", middleware.Auth(eventhandlers.GetEvents))
	mux.HandleFunc("POST /events/vote", middleware.Auth(eventhandlers.VoteEvent))

	// notifications
	mux.HandleFunc("GET /notifications", middleware.Auth(notifhandlers.DefaultHandler.GetNotifications))
	mux.HandleFunc("POST /notifications/{id}/read", middleware.Auth(notifhandlers.DefaultHandler.MarkReadByID))

	// chat routes
	mux.HandleFunc("POST /messages", middleware.Auth(chathandlers.DefaultHandler.SendDMHandler))
	mux.HandleFunc("GET /messages/{user_id}", middleware.Auth(chathandlers.DefaultHandler.DMHistoryHandler))
	mux.HandleFunc("POST /groups/{group_id}/messages", middleware.Auth(chathandlers.DefaultHandler.SendGroupHandler))
	mux.HandleFunc("GET /groups/{group_id}/messages", middleware.Auth(chathandlers.DefaultHandler.GroupHistoryHandler))
	mux.HandleFunc("GET /conversations", middleware.Auth(chathandlers.DefaultHandler.ConversationsHandler))

	// websocket routes
	mux.HandleFunc("GET /ws", middleware.Auth(wshandlers.WSHandler(websocket.DefaultHub)))

	return mux
}
