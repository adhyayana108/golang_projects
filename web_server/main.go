import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// middleware

type responseWriter struct{
	http.ResponseWriter
	statusCode int 
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w , statusOK}
}

func (rw *responseWriter) WriteHeader(code int){
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter , r *http.Request){
		start := time.Now()
		rw := newResponseWriter(w)
		next.ServerHTTP(rw,r)
		duration := time.Since(start)
		log.Printf("[%s] %s %s | Status %d | Duration: %v",
		r.Method , r.RemoteAddr , r.URL.Path, rw.statusCode, duration)
	})
}


