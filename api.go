package dronecluster

import (
	"io"
	"net/http"
	_ "net/http/pprof"

	"github.com/gorilla/mux"
)

type ApiServer struct {
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	//todo: post a tyrant job and start task
}

func newapiserver() *ApiServer {
	return &ApiServer{}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "drone cluster")
}

func (s *ApiServer) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/{project_id}", WebhookHandler).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(":10000", nil)
}
