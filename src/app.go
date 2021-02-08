package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type LongUrl struct {
	Url string `json:"url"`
}

type TinyUrl struct {
	Key string `json:"key"`
	Url string `json:"url"`
}

var urlCache = make(map[string]TinyUrl)

func main() {
	router := mux.NewRouter()
	// Routes  consist of a path and a handler function.
	router.Path("/api/heartbeat").Methods(http.MethodGet).HandlerFunc(HeartbeatHandler)
	router.Path("/{key}").Methods(http.MethodGet).HandlerFunc(UrlHandler)

	router.Path("/new").Methods(http.MethodPost).HandlerFunc(CreateTinyUrl)

	srv := &http.Server{
		Handler: router,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Bind to a port and pass our router in
	log.Fatal(srv.ListenAndServe())
}

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("/api/heartbeat")
	w.WriteHeader(http.StatusOK)
}

func GetUuid() string {
	for {
		// Converts key to string and grab the first 8 characters
		// TODO find a better way!
		key := uuid.New().String()[:8]
		_, prs := urlCache[key]
		// prs checks for collision if already in cache
		if !prs {
			return key
		}
	}
}

func CreateTinyUrl(w http.ResponseWriter, r *http.Request) {
	u := LongUrl{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := GetUuid()
	tiny := TinyUrl{Key: key, Url: u.Url}
	resp, err := json.Marshal(&tiny)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	urlCache[key] = tiny

	log.Print("/new, ", u.Url, ", ", key)
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func UrlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tiny, prs := urlCache[vars["key"]]
	if prs == true {
		log.Print("/", vars["key"], ", ", tiny.Url)
		http.Redirect(w, r, tiny.Url, http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}
