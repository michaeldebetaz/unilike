package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Logging() Middleware {
	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {
		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			// Do middleware things
			start := time.Now()
			defer logRequest(r, start)

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

func logRequest(r *http.Request, start time.Time) {
	msg := fmt.Sprintf("%s %s - %s", r.Method, r.URL.Path, time.Since(start))
	slog.Info(msg)
}

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
