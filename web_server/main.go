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

//form handler 

func formHandler(w http.ResponseWriter , r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w ,`{"error":"Only POST is accepted"}` , http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"ParseForm failed: %v"}`, err), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	address := r.FormValue("address")
	email := r.FormValue("email")

	if name == ""|| email== "" {
		http.Error(w , `{"error": "Name and email are required fields."}`, http.StatusBadRequest)
		return
	
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.Status.OK)

	payload := map[string]string{
		"status":  "success",
		"name":    name,
		"address": address,
		"email":   email,
	}
 
	json.NewEncoder(w).Encode(payload)
}

func main() {
	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
    mux.Handle("/", fileServer)

	mux.HandleFunc("/hello", helloHandler)
    mux.HandleFunc("/form", formHandler)
    mux.HandleFunc("/api/health", healthHandler)

	server := &http.Server{
        Addr:         ":" + port,
        Handler:      loggingMiddleware(mux),
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

	fmt.Printf("Go Web Server\n")
	fmt.Printf("Listening on http://localhost:%s \n", port)
	fmt.Printf("Routes: \n")
	fmt.Printf("GET  /              → index.html\n")
	fmt.Printf("GET  /form.html     → contact form\n")
	fmt.Printf("GET  /hello         → hello endpoint\n")
	fmt.Printf("POST /form          → form handler\n")
	fmt.Printf("GET  /api/health    → health check\n")
	
 
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

