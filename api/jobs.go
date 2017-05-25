package api

import (
	"github.com/alienantfarm/anthive/db"
	"github.com/alienantfarm/anthive/utils"
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"net/http"
)

type JobPost struct {
	ImageId int      `json:"image"`
	Cmd     []string `json:"command"`
}

func jobPost(w http.ResponseWriter, r *http.Request) {
	jp := &JobPost{}

	if err := utils.Decode(r, jp); err != nil {
		glog.Errorf("%s", err)
		err.Dump(w)
		return
	}

	query := "WITH job as ("
	query += "  INSERT INTO anthive.job (fk_image, command) "
	query += "  VALUES ($1, $2) RETURNING anthive.job.id "
	query += ") "
	query += "SELECT job.id, anthive.image.archive "
	query += "FROM job, anthive.image "
	query += "WHERE anthive.image.id = $1"

	j := &structs.Job{}
	j.Image.Id = jp.ImageId
	j.Image.Cmd = jp.Cmd

	qr := db.Conn().QueryRow(query, jp.ImageId, pq.Array(jp.Cmd))
	if err := qr.Scan(&j.Id, &j.Image.Archive); err != nil {
		glog.Errorln(err)
		return
	}
	if glog.V(2) {
		glog.Infof(utils.MarshalJSON(j))
	}
	Scheduler.ProcessJob(j)

	if err := utils.Encode(w, j); err != nil {
		glog.Errorln(err)
		return
	}
}

func jobGet(w http.ResponseWriter, r *http.Request) {

	query := "SELECT job.id, job.state, job.command, image.id, image.archive "
	query += "FROM anthive.job AS job, anthive.image AS image "
	query += "WHERE image.id = job.fk_image"

	rows, err := db.Conn().Query(query)
	if err != nil {
		glog.Errorln(err)
		return
	}
	defer rows.Close()

	jobs := []*structs.Job{}
	for rows.Next() {
		job := &structs.Job{}
		image := &job.Image

		err = rows.Scan(
			&job.Id, &job.State, pq.Array(&image.Cmd), &image.Id, &image.Archive,
		)
		if err != nil {
			glog.Errorln(err)
			return
		}
		jobs = append(jobs, job)
	}

	if err := utils.Encode(w, structs.Jobs{jobs}); err != nil {
		glog.Errorln(err)
		return
	}
}

func jobGetId(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	j := &structs.Job{}
	i := &j.Image

	query := "SELECT job.id, job.state, job.command, image.id, image.archive "
	query += "FROM anthive.job AS job, anthive.image AS image "
	query += "WHERE image.id = job.fk_image AND job.id = $1"

	err := db.Conn().QueryRow(query, id).Scan(
		&j.Id, &j.State, pq.Array(&i.Cmd), &i.Id, &i.Archive,
	)
	if err != nil {
		glog.Errorln(err)
		return
	}

	if err := utils.Encode(w, j); err != nil {
		glog.Errorln(err)
		return
	}
}

func init() {
	Router.HandleFunc("/jobs", jobPost).Methods("POST")
	Router.HandleFunc("/jobs", jobGet).Methods("GET")
	Router.HandleFunc("/jobs/{id:[0-9]+}", jobGetId).Methods("GET")
}
