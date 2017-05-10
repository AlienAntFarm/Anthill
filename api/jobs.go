package api

import (
	"github.com/alienantfarm/anthive/db"
	"github.com/alienantfarm/anthive/utils"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	JOB_NEW = iota
	JOB_PENDING
	JOB_FINISH
	JOB_ERROR
)

var JOB_STATES = [...]string{
	"NEW",
	"PENDING",
	"FINISH",
	"ERROR",
}

type Job struct {
	Id    int `json:"id"`
	State int `json:"state"`
}

type Jobs struct {
	Jobs []*Job `json:"jobs"`
}

func jobPost(w http.ResponseWriter, r *http.Request) {
	var err error
	query := "INSERT INTO anthive.job "
	query += "DEFAULT VALUES "
	query += "RETURNING anthive.job.id"

	j := &Job{}
	err = db.Conn().QueryRow(query).Scan(&j.Id)
	if err != nil {
		glog.Errorln(err)
		return
	}
	Scheduler.AddJob(j.Id)
	err = utils.Encode(w, j)
	if err != nil {
		glog.Errorln(err)
		return
	}
}

func jobGet(w http.ResponseWriter, r *http.Request) {
	var err error
	query := "SELECT anthive.job.id "
	query += "FROM anthive.job"
	rows, err := db.Conn().Query(query)
	if err != nil {
		glog.Errorln(err)
		return
	}
	defer rows.Close()
	jobs := []*Job{}
	for rows.Next() {
		job := &Job{}
		err = rows.Scan(&job.Id)
		if err != nil {
			glog.Errorln(err)
			return
		}
		jobs = append(jobs, job)
	}
	err = utils.Encode(w, Jobs{jobs})
	if err != nil {
		glog.Errorln(err)
		return
	}
}

func jobGetId(w http.ResponseWriter, r *http.Request) {
	var err error
	id := mux.Vars(r)["id"]
	j := &Job{}

	query := "SELECT anthive.job.id "
	query += "FROM anthive.job "
	query += "WHERE anthive.job.id = $1"

	err = db.Conn().QueryRow(query, id).Scan(&j.Id)
	if err != nil {
		glog.Errorln(err)
		return
	}
	err = utils.Encode(w, j)
	if err != nil {
		glog.Errorln(err)
		return
	}
}

func init() {
	Router.HandleFunc("/jobs", jobPost).Methods("POST")
	Router.HandleFunc("/jobs", jobGet).Methods("GET")
	Router.HandleFunc("/jobs/{id:[0-9]+}", jobGetId).Methods("GET")
}
