package main

import (
	"log"
	"net/http"

	"zone/internal/auth"
	"zone/internal/database"
	"zone/internal/middleware"
	"zone/internal/posts"
	"zone/internal/profile"
	"zone/internal/users"
	"zone/internal/utils"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/register", auth.Register)
	mux.HandleFunc("/login", auth.Login)
	mux.HandleFunc("/me", auth.Me)
	mux.HandleFunc("/logout", auth.Logout)
	mux.HandleFunc("/upload/avatar", utils.UploadAvatar)
	mux.HandleFunc("GET /users/search", users.SearchUsersAPI)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	mux.HandleFunc("GET /posts", posts.GetPostsAPI)
	mux.HandleFunc("POST /posts/create", posts.CreatePostAPI)
	mux.HandleFunc("GET /profile/{id}", profile.GetProfileAPI)
	mux.HandleFunc("POST /profile/toggle-visibility", profile.ToggleVisibilityAPI)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", middleware.CORS(mux))
}
