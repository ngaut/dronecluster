package service

import (
	"io"
	"net/http"
	_ "net/http/pprof"

	"github.com/gorilla/mux"
	"github.com/juju/errors"
	"github.com/ngaut/dronecluster/task"
	log "github.com/ngaut/logging"
)

type ApiServer struct {
	Server       string //todo: get from zookeeper
	ExecutorUrls string
}

func (s *ApiServer) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	//todo: post a tyrant job and start task
	err := r.ParseForm()
	if err != nil {
		http.Error(w, errors.ErrorStack(err), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	log.Debug(vars)

	repo := r.Form.Get("repo")
	log.Debug(r.Form, "repo", repo)
	h := &task.JobHelper{Server: s.Server, ExecutorUrls: s.ExecutorUrls}
	job := h.BuildRepoJob(repo)
	if err := h.CreateJob(job); err != nil {
		http.Error(w, errors.ErrorStack(err), http.StatusInternalServerError)
		return
	}

	log.Debugf("%+v", job)

	if err := h.RunJob(job); err != nil {
		http.Error(w, errors.ErrorStack(err), http.StatusInternalServerError)
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "drone cluster")
}

func (s *ApiServer) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/{project_id}", s.WebhookHandler).Methods("POST")
	http.Handle("/", r)
	http.ListenAndServe(":10001", nil)
}
