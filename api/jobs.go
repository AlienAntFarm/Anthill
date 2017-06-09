package api

import (
	"github.com/alienantfarm/anthive/ext/db"
	"github.com/alienantfarm/anthive/utils"
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
)

func jobPost(w http.ResponseWriter, r *http.Request) {
	jp := struct {
		ImageId int      `json:"image"`
		Cmd     []string `json:"command"`
		Env     []string `json:"env"`
		Cwd     string   `json:"cwd"`
	}{}

	if err := utils.Decode(r, &jp); err != nil {
		glog.Errorf("%s", err)
		err.Dump(w)
		return
	}
	job := &db.Job{Cmd: jp.Cmd, Env: jp.Env, Cwd: jp.Cwd}
	if err := job.Create(jp.ImageId); err != nil {
		glog.Errorln(err)
	}

	if glog.V(2) {
		glog.Infof(utils.MarshalJSON(job))
	}
	Scheduler.ProcessJob((*structs.Job)(job))

	if err := utils.Encode(w, job); err != nil {
		glog.Errorln(err)
	}
}

func jobGet(w http.ResponseWriter, r *http.Request) {
	jobs := &db.Jobs{}
	if err := jobs.Get(structs.JOB_ERROR); err != nil {
		glog.Errorf("%s", err)
	} else if err := utils.Encode(w, jobs); err != nil {
		glog.Errorf("%s", err)
	}
}

func jobGetId(w http.ResponseWriter, r *http.Request) {
	job := &db.Job{}
	if err := job.Get(mux.Vars(r)["id"]); err != nil {
		glog.Errorln(err)
	} else if err := utils.Encode(w, job); err != nil {
		glog.Errorln(err)
	}
}

func init() {
	Router.HandleFunc("/jobs", jobPost).Methods("POST")
	Router.HandleFunc("/jobs", jobGet).Methods("GET")
	Router.HandleFunc("/jobs/{id:[0-9]+}", jobGetId).Methods("GET")
}
