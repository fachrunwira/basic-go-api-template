package ratelimiting

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

type ClientLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

func NewClientLimiter(r rate.Limit, b int) *ClientLimiter {
	return &ClientLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

func (cl *ClientLimiter) getLimiter(ip string) *rate.Limiter {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	limiter, exists := cl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(cl.rate, cl.burst)
		cl.limiters[ip] = limiter
	}
	return limiter
}

func (cl *ClientLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			limiter := cl.getLimiter(ip)

			if !limiter.Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"status": false,
					"code":   http.StatusTooManyRequests,
					"errors": "Too many requests. Please try again later.",
				})
			}

			return next(c)
		}
	}
}
