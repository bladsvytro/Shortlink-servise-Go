package app

import (
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	mu          sync.Mutex
	ips         map[string]*tokenBucket
	limit       int           // requests per window
	window      time.Duration // time window
	cleanupInterval time.Duration
	logger      *zap.Logger
	stopCleanup chan struct{}
}

type tokenBucket struct {
	tokens     int
	lastRefill time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration, logger *zap.Logger) *RateLimiter {
	rl := &RateLimiter{
		ips:         make(map[string]*tokenBucket),
		limit:       limit,
		window:      window,
		cleanupInterval: window * 2,
		logger:      logger,
		stopCleanup: make(chan struct{}),
	}
	go rl.cleanupOldEntries()
	return rl
}

// Allow checks if the request from ip is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.ips[ip]
	if !exists {
		rl.ips[ip] = &tokenBucket{
			tokens:     rl.limit - 1,
			lastRefill: now,
		}
		return true
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(bucket.lastRefill)
	refillTokens := int(elapsed / rl.window)
	if refillTokens > 0 {
		bucket.tokens += refillTokens * rl.limit
		if bucket.tokens > rl.limit {
			bucket.tokens = rl.limit
		}
		bucket.lastRefill = now
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}
	return false
}

// cleanupOldEntries periodically removes old entries to prevent memory leak
func (rl *RateLimiter) cleanupOldEntries() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for ip, bucket := range rl.ips {
				if now.Sub(bucket.lastRefill) > rl.cleanupInterval {
					delete(rl.ips, ip)
				}
			}
			rl.mu.Unlock()
		case <-rl.stopCleanup:
			return
		}
	}
}

// Stop stops the cleanup goroutine
func (rl *RateLimiter) Stop() {
	close(rl.stopCleanup)
}

// RateLimitMiddleware returns a middleware that limits requests per IP
func (a *Application) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// Create rate limiter if not exists (global for the application)
	if a.rateLimiter == nil {
		a.rateLimiter = NewRateLimiter(100, time.Minute, a.logger) // 100 requests per minute
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		// Extract IP from X-Forwarded-For if behind proxy
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		if !a.rateLimiter.Allow(ip) {
			a.logger.Warn("Rate limit exceeded", zap.String("ip", ip))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	}
}