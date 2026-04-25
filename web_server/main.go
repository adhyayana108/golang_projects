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

// health check handler

var serverStartTime = time.Now()

func HealthHandler(w http.ResponseWriter , r *http.Request){
	if r.Method != http.MethodGet {
		http.Error(w , `{"error":"Method not allowed"}` , http.StatusMethodNotAllowed )
		return
	}

	w.Header().Set("Content-Type" . "application/json")
	uptime := time,since(serverStartTime).Round(time.second)

	payload := map[string]string {
		"status":    "OK",
        "uptime":    uptime.String(),
        "timestamp": time.Now().UTC().Format(time.RFC3339),
        "version":   "1.0.0",
	}

	json.NewEncoder(w).Enocode(payload)
}

// hello handler 

func helloHandler(w http.ResponseWriter , r *http.Request)
{
	if r.URL.Path != "/hello" {
		http.Error(w, "404-Page Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w , "Method not supported" , http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Hello! The Go Web Server is running successfully.")
}

