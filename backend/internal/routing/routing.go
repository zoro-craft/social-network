package routing

import (
	"net/http"
	"time"

	"zone/internal/auth"
	"zone/internal/middleware"
	"zone/internal/posts"
	"zone/internal/profile"
	"zone/internal/reactions"
	"zone/internal/users"
	"zone/internal/utils"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	limiter := middleware.NewRateLimiter(10, time.Second*3) 

	// Public routes 
	mux.HandleFunc("/register", auth.Register)
	mux.HandleFunc("/login", auth.Login)
	mux.HandleFunc("/me", auth.Me)
	mux.HandleFunc("GET /users/search", users.SearchUsersAPI)
	mux.HandleFunc("GET /posts", posts.GetPostsAPI)
	mux.HandleFunc("GET /profile/{id}", profile.GetProfileAPI)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Protected routes
	mux.HandleFunc("/logout", middleware.Auth(auth.Logout))
	mux.HandleFunc("/upload/avatar", middleware.Auth(utils.UploadAvatar))
	mux.HandleFunc("POST /posts/create", middleware.Auth(posts.CreatePostAPI))
	mux.HandleFunc("POST /profile/toggle-visibility", middleware.Auth(profile.ToggleVisibilityAPI))
	mux.HandleFunc("POST /posts/react",limiter.Middleware(middleware.Auth(reactions.ReactToPost)))
return middleware.CORS(limiter.Middleware(mux.ServeHTTP))
}