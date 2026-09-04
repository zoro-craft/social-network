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
	mux.HandleFunc("GET /profile/{id}", middleware.OptionalAuth(profile.GetProfileAPI))
	mux.HandleFunc("POST /profile/edit", middleware.Auth(profile.EditProfileAPI))
	mux.HandleFunc("POST /profile/toggle-visibility", middleware.Auth(profile.ToggleVisibilityAPI))
	mux.HandleFunc("GET /users/search", middleware.Auth(users.SearchUsersAPI))
	// followers routes
	mux.HandleFunc("POST /follow/{id}", middleware.Auth(followers.FollowAPI))
	mux.HandleFunc("POST /follow/{id}/unfollow", middleware.Auth(followers.UnfollowAPI))
	mux.HandleFunc("POST /follow/{id}/accept", middleware.Auth(followers.AcceptFollowAPI))
	mux.HandleFunc("POST /follow/{id}/decline", middleware.Auth(followers.DeclineFollowAPI))

	// media & uploads routes
	mux.HandleFunc("/upload/avatar", utils.UploadAvatar)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// global middleware
	return middleware.CORS(limiter.Middleware(mux.ServeHTTP))
}
