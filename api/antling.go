package api

import (
	"github.com/alienantfarm/anthive/db"
	"github.com/alienantfarm/anthive/utils"
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func antlingPost(w http.ResponseWriter, r *http.Request) {
	var err error
	query := "INSERT INTO anthive.antling "
	query += "DEFAULT VALUES "
	query += "RETURNING anthive.antling.id"

	a := &structs.Antling{0, []*structs.Job{}}
	err = db.Conn().QueryRow(query).Scan(&a.Id)
	if err != nil {
		glog.Errorln(err)
		return
	}
	Scheduler.AddAntling(a.Id)
	if glog.V(2) {
		glog.Infof("created antling with id %d", a.Id)
	}
	w.WriteHeader(http.StatusCreated)
	err = utils.Encode(w, a)
	if err != nil {
		glog.Errorf("%s", err)
		return
	}
}

func antlingGet(w http.ResponseWriter, r *http.Request) {
	var err error
	query := "SELECT anthive.antling.id "
	query += "FROM anthive.antling"
	rows, err := db.Conn().Query(query)
	if err != nil {
		glog.Errorln(err)
		return
	}
	defer rows.Close()
	antlings := []*structs.Antling{}
	for rows.Next() {
		antling := &structs.Antling{}
		err = rows.Scan(&antling.Id)
		if err != nil {
			glog.Errorln(err)
			return
		}
		antlings = append(antlings, antling)
	}
	err = utils.Encode(w, structs.Antlings{antlings})
	if err != nil {
		glog.Errorln(err)
		return
	}
}

func antlingGetId(w http.ResponseWriter, r *http.Request) {
	var err error
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	a := &structs.Antling{0, structs.ListJobs(Scheduler.queue[id])}

	query := "SELECT anthive.antling.id "
	query += "FROM anthive.antling "
	query += "WHERE anthive.antling.id = $1"

	err = db.Conn().QueryRow(query, id).Scan(&a.Id)
	if err != nil {
		glog.Errorln(err)
		return
	}
	if glog.V(2) {
		glog.Infof(utils.MarshalJSON(a))
	}
	err = utils.Encode(w, a)
	if err != nil {
		glog.Errorln(err)
		return
	}
}

func antlingPatchId(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	antling := &structs.Antling{}

	// retrieve objects from request (an antling with a list of jobs)
	if err := utils.Decode(r, antling); err != nil {
		glog.Errorf("%s", err)
		err.Dump(w)
		return
	}
	// check if ids are matching, dump an error otherwise
	if antling.Id != id {
		errUI := &utils.UnmatchingIds{[2]int{antling.Id, id}}
		glog.Errorf("%s", errUI)
		errUI.Dump(w)
		return
	}

	// compare jobs statuses to not overwhelm the scheduler channel
	helper := []*structs.Job{}
	jobs := Scheduler.queue[id]

	for _, job := range antling.Jobs { // BEWARE: Do not modify jobs from scheduler
		job.IdAntling = id
		if j, ok := jobs[job.Id]; !ok {
			// job not here anymore just pass
		} else if j.State != job.State {
			helper = append(helper, job)
		}
	}

	// send jobs to update to the scheduler
	for _, job := range helper {
		if glog.V(2) {
			glog.Infof("Updating state of job %d", job.Id)
			glog.Infof(utils.MarshalJSON(job))
		}
		Scheduler.ProcessJob(job)
	}
}

func init() {
	Router.HandleFunc("/antlings", antlingPost).Methods("POST")
	Router.HandleFunc("/antlings", antlingGet).Methods("GET")
	Router.HandleFunc("/antlings/{id:[0-9]+}", antlingGetId).Methods("GET")
	Router.HandleFunc("/antlings/{id:[0-9]+}", antlingPatchId).Methods("PATCH")
}
