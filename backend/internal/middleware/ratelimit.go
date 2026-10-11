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
	mu        sync.Mutex
	clients   map[string]bucket
	limit     int
	window    time.Duration
	lastSweep time.Time
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
		key := clientKey(r)

		rl.mu.Lock()
		rl.sweepLocked(now)
		b := rl.clients[key]
		if now.After(b.reset) {
			b = bucket{reset: now.Add(rl.window)}
		}

		if b.count >= rl.limit {
			rl.mu.Unlock()
			utils.HandleError(w, http.StatusTooManyRequests, "too many requests, please slow down")
			return
		}

		b.count++
		rl.clients[key] = b
		rl.mu.Unlock()

		next(w, r)
	}
}

// sweepLocked drops the buckets whose window has already closed, at most once per
// window, so a long-running server does not accumulate one entry per client IP it
// has ever seen. Must be called with rl.mu held.
func (rl *RateLimiter) sweepLocked(now time.Time) {
	if now.Sub(rl.lastSweep) < rl.window {
		return
	}
	rl.lastSweep = now

	for key, b := range rl.clients {
		if now.After(b.reset) {
			delete(rl.clients, key)
		}
	}
}

// clientKey buckets by the client IP, so every request from the same IP shares one bucket.
func clientKey(r *http.Request) string {
	return "ip:" + clientIP(r.RemoteAddr)
}

func clientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil && host != "" {
		return host
	}
	return remoteAddr
}
