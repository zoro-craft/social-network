package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
	"zone/internal/utils"
)

type bucket struct {
	count int
	reset time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]bucket
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	if limit < 1 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}

	return &RateLimiter{
		clients: make(map[string]bucket),
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		ip := clientIP(r.RemoteAddr)

		rl.mu.Lock()
		b := rl.clients[ip]
		if now.After(b.reset) {
			b = bucket{reset: now.Add(rl.window)}
		}

		if b.count >= rl.limit {
			rl.mu.Unlock()
			utils.HandleError(w, http.StatusTooManyRequests, "too many requests, please slow down")
			return
		}

		b.count++
		rl.clients[ip] = b
		rl.mu.Unlock()

		next(w, r)
	}
}

func clientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil && host != "" {
		return host
	}
	return remoteAddr
}
