package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/internal/helpers/response"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// RateLimit middleware para limitar requisições por IP
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !rl.allow(clientIP) {
			response.Error(c, http.StatusTooManyRequests, "rate_limit.exceeded", "Too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Remove requisições antigas
	requests := rl.requests[key]
	validRequests := make([]time.Time, 0, len(requests))

	for _, reqTime := range requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Verifica se excedeu o limite
	if len(validRequests) >= rl.limit {
		rl.requests[key] = validRequests
		return false
	}

	// Adiciona nova requisição
	validRequests = append(validRequests, now)
	rl.requests[key] = validRequests

	return true
}

// Cleanup remove entradas antigas periodicamente
func (rl *RateLimiter) Cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	for key, requests := range rl.requests {
		validRequests := make([]time.Time, 0, len(requests))

		for _, reqTime := range requests {
			if reqTime.After(windowStart) {
				validRequests = append(validRequests, reqTime)
			}
		}

		if len(validRequests) == 0 {
			delete(rl.requests, key)
		} else {
			rl.requests[key] = validRequests
		}
	}
}

// AuthRateLimiter rate limiter específico para endpoints de auth
func NewAuthRateLimiter() *RateLimiter {
	// 10 tentativas por IP a cada 15 minutos
	return NewRateLimiter(10, 15*time.Minute)
}
