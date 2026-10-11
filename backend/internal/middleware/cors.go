package middleware

import (
	"net/http"
	"strings"
)

// allowedOrigins are the browser origins allowed to call the API directly. The
// Next.js frontend proxies /api same-origin, but the Expo web build runs on its
// own dev port and has no proxy, so it needs real CORS.
var allowedOrigins = []string{
	"http://localhost:8081", // Expo web dev server
	"http://127.0.0.1:8081",
	"http://localhost:19006", // Expo Go
	"http://localhost:3000",  // Next.js, when it is not proxying
}

// allowedHeaders covers what the clients actually send: JSON bodies on native,
// plus multipart uploads, and Cookie on the web session.
var allowedHeaders = "Origin, Content-Type, Accept, Authorization, X-Requested-With, Cookie"

// CORS answers preflights and tags responses so browser clients (notably Expo
// web on :8081) may talk to the API cross-origin. Credentials are allowed, so a
// wildcard origin is not usable — the request origin is echoed back instead, and
// only when it is on the allowlist.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" && originAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Max-Age", "600")
			// Responses vary by Origin, so caches must not share them across origins.
			w.Header().Add("Vary", "Origin")
		}

		// A preflight carries no credentials and must not run the handler.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func originAllowed(origin string) bool {
	for _, allowed := range allowedOrigins {
		if strings.EqualFold(origin, allowed) {
			return true
		}
	}
	return false
}
