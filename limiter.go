package main

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const NUM_REQUESTS = 2
const NUM_SECONDS = 20

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
}

var rateLimiter = NewRateLimiter()

func RateLimited(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		ip := getIP(r)
		limiter := rateLimiter.GetLimiter(ip)

		if limiter.Allow() {
			log.Println("Request allowed for IP:", ip)
			next.ServeHTTP(w, r)
		} else {
			log.Println("Request rate limited for IP:", ip)
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	}
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

func (r *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()

	limiter, exists := r.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(NUM_SECONDS*time.Second), NUM_REQUESTS) // 2 requests every 20 seconds
		r.limiters[ip] = limiter
	}

	return limiter
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	// Extract the IP address without the port
	if strings.Contains(ip, ":") {
		ip = ip[:strings.LastIndex(ip, ":")]
		if strings.Count(ip, ":") > 1 { // IPv6 address
			ip = strings.Trim(ip, "[]")
		}
	}

	log.Println("Extracted IP address:", ip)
	return ip
}
