package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("/api/heartbeat")
	w.Write([]byte("OK"))
}

func main() {
	router := mux.NewRouter()
	// Routes  consist of a path and a handler function.
	router.HandleFunc("/api/heartbeat", HeartbeatHandler)

	srv := &http.Server{
		Handler: router,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Bind to a port and pass our router in
	log.Print("tinygo live on :8080...")
	log.Fatal(srv.ListenAndServe())
}
