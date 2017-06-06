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
	Env     []string `json:"env"`
	Cwd     string   `json:"cwd"`
}

func jobPost(w http.ResponseWriter, r *http.Request) {
	jp := &JobPost{}

	if err := utils.Decode(r, jp); err != nil {
		glog.Errorf("%s", err)
		err.Dump(w)
		return
	}

	query := "WITH j as ("
	query += "  INSERT INTO anthive.job (fk_image, command, environment, cwd) "
	query += "  VALUES ($1, $2, $3, $4) RETURNING anthive.job.id "
	query += ") "
	query += "SELECT j.id, i.id, i.archive, i.command, i.environment, i.cwd, i.hostname "
	query += "FROM j, anthive.image as i "
	query += "WHERE i.id = $1"

	j := &structs.Job{Cmd: jp.Cmd, Env: jp.Env, Cwd: jp.Cwd}
	i := &j.Image

	args := []interface{}{jp.ImageId, pq.Array(j.Cmd), pq.Array(j.Env), j.Cwd}
	argsScan := []interface{}{
		&j.Id, &i.Id, &i.Archive, pq.Array(&i.Cmd), pq.Array(&i.Env), &i.Cwd, &i.Hostname,
	}

	if err := db.Conn().QueryRow(query, args...).Scan(argsScan...); err != nil {
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

	query := "SELECT j.id, j.state, j.cwd, j.command, j.environment, "
	query += "  i.id, i.archive, i.command, i.environment, i.cwd, i.hostname "
	query += "FROM anthive.job AS j, anthive.image AS i "
	query += "WHERE i.id = j.fk_image"

	rows, err := db.Conn().Query(query)
	if err != nil {
		glog.Errorln(err)
		return
	}
	defer rows.Close()

	jobs := []*structs.Job{}
	for rows.Next() {
		j := &structs.Job{}
		i := &j.Image

		args := []interface{}{
			&j.Id, &j.State, &j.Cwd, pq.Array(&j.Cmd), pq.Array(&j.Env),
			&i.Id, &i.Archive, pq.Array(&i.Cmd), pq.Array(&i.Env), &i.Cwd, &i.Hostname,
		}

		if err := rows.Scan(args...); err != nil {
			glog.Errorln(err)
			return
		}
		jobs = append(jobs, j)
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

	query := "SELECT j.id, j.state, j.cwd, j.command, j.environment, "
	query += "  i.id, i.archive, i.command, i.environment, i.cwd, i.hostname "
	query += "FROM anthive.job AS j, anthive.image AS i "
	query += "WHERE i.id = j.fk_image AND j.id = $1"

	args := []interface{}{
		&j.Id, &j.State, &j.Cwd, pq.Array(&j.Cmd), pq.Array(&j.Env),
		&i.Id, &i.Archive, pq.Array(&i.Cmd), pq.Array(&i.Env), &i.Cwd, &i.Hostname,
	}
	if err := db.Conn().QueryRow(query, id).Scan(args...); err != nil {
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
