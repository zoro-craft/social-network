package routing

import (
	"net/http"
	"time"

	"zone/internal/auth"
	"zone/internal/comments"
	"zone/internal/followers"
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
	return middleware.CORS(limiter.Middleware(mux.ServeHTTP))
}