package ratelimit

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type clientEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type clientLimiter struct {
	clients         map[string]*clientEntry
	mu              sync.Mutex
	rate            rate.Limit
	burst           int
	cleanupInterval time.Duration
	entryTTL        time.Duration
	abort           context.CancelFunc
	logger          *slog.Logger
}

func New(r rate.Limit, b int, interval, ttl time.Duration, logs *slog.Logger) *clientLimiter {
	cl := &clientLimiter{
		clients:         make(map[string]*clientEntry),
		rate:            r,
		burst:           b,
		cleanupInterval: interval,
		entryTTL:        ttl,
		logger:          logs,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cl.abort = cancel
	go cl.startCleanup(ctx)
	return cl
}

func (cl *clientLimiter) getLimiter(ip string) *rate.Limiter {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	entry, exists := cl.clients[ip]
	if !exists {
		limiter := rate.NewLimiter(cl.rate, cl.burst)
		cl.clients[ip] = &clientEntry{limiter: limiter, lastSeen: time.Now()}

		cl.setMessageLog("new", "client", ip)
		return limiter
	}

	entry.lastSeen = time.Now()
	return entry.limiter
}

func (cl *clientLimiter) cleanup() {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	removed := 0
	total := len(cl.clients)

	for ip, entry := range cl.clients {
		if time.Since(entry.lastSeen) > 5*time.Minute {
			delete(cl.clients, ip)
			removed++
			cl.setMessageLog("removed staled", "client", ip, "(last seen", time.Since(entry.lastSeen), ") ago")
		}
	}

	cl.setMessageLog("cleanup done: ", "removed", removed, " / total", total)
}

func (cl *clientLimiter) Stop() {
	if cl.abort != nil {
		cl.abort()
	}
}

func (cl *clientLimiter) startCleanup(ctx context.Context) {
	ticker := time.NewTicker(cl.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			cl.cleanup()
		case <-ctx.Done():
			return
		}
	}
}

func (cl *clientLimiter) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		limiter := cl.getLimiter(ip)

		if !limiter.Allow() {
			cl.setMessageLog("rate limit exceeded for", "ip", ip)
			ctx.Header("Retry-After", "60")
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"errors": "too many request",
			})
			return
		}

		ctx.Next()
	}
}

func (cl *clientLimiter) setMessageLog(message string, args ...any) {
	if cl.logger != nil {
		cl.logger.Info(message, args...)
	}
}
