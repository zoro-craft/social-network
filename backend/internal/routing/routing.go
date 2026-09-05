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
	mux.HandleFunc("GET /posts/{id}", posts.GetPostByIDAPI)
	mux.HandleFunc("POST /posts/create", middleware.Auth(posts.CreatePostAPI))

	// comments routes
	mux.HandleFunc("GET /posts/{id}/comment", comments.GetCommentsAPI)
	mux.HandleFunc("POST /posts/{id}/comment", middleware.Auth(comments.CreateCommentAPI))
	mux.HandleFunc("POST /posts/react", middleware.Auth(reactions.ReactToPost))

	// profile and users routes
	mux.HandleFunc("GET /profile/{id}", profile.GetProfileAPI)
	mux.HandleFunc("POST /profile/toggle-visibility", middleware.Auth(profile.ToggleVisibilityAPI))
	mux.HandleFunc("GET /users/search", middleware.Auth(users.SearchUsersAPI))
	// followers routes
	mux.HandleFunc("GET /follow/requests", middleware.Auth(followers.GetPendingRequestsAPI))
	mux.HandleFunc("GET /follow/{id}/status", followers.GetFollowStatusAPI)
	mux.HandleFunc("GET /profile/{id}/followers", followers.GetFollowersAPI)
	mux.HandleFunc("GET /profile/{id}/following", followers.GetFollowingAPI)
	mux.HandleFunc("POST /follow/{id}", middleware.Auth(followers.FollowAPI))
	mux.HandleFunc("POST /follow/{id}/unfollow", middleware.Auth(followers.UnfollowAPI))
	mux.HandleFunc("POST /follow/{id}/accept", middleware.Auth(followers.AcceptFollowAPI))
	mux.HandleFunc("POST /follow/{id}/decline", middleware.Auth(followers.DeclineFollowAPI))

	// media & uploads routes
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// global middleware
	return middleware.CORS(limiter.Middleware(mux.ServeHTTP))
}