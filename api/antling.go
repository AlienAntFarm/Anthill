package api

import (
	"github.com/alienantfarm/anthive/db"
	"github.com/alienantfarm/anthive/utils"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Antling struct {
	Id   int   `json:"id"`
	Jobs []int `json:"jobs"`
}

type Antlings struct {
	Antlings []*Antling `json:"antlings"`
}

func antlingPost(w http.ResponseWriter, r *http.Request) {
	var err error
	query := "INSERT INTO anthive.antling "
	query += "DEFAULT VALUES "
	query += "RETURNING anthive.antling.id"

	a := &Antling{0, []int{}}
	err = db.Conn().QueryRow(query).Scan(&a.Id)
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	Scheduler.AddAntling(a.Id)
	if glog.V(2) {
		glog.Infof("created antling with id %d", a.Id)
	}
	err = utils.Encode(w, a)
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
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
	antlings := []*Antling{}
	for rows.Next() {
		antling := &Antling{}
		err = rows.Scan(&antling.Id)
		if err != nil {
			glog.Errorln(err)
			return
		}
		antlings = append(antlings, antling)
	}
	err = utils.Encode(w, Antlings{antlings})
	if err != nil {
		glog.Errorln(err)
		return
	}
}

func antlingGetId(w http.ResponseWriter, r *http.Request) {
	var err error
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	a := &Antling{0, Scheduler.GetJobs(id)}

	query := "SELECT anthive.antling.id "
	query += "FROM anthive.antling "
	query += "WHERE anthive.antling.id = $1"

	err = db.Conn().QueryRow(query, id).Scan(&a.Id)
	if err != nil {
		glog.Errorln(err)
		return
	}
	err = utils.Encode(w, a)
	if err != nil {
		glog.Errorln(err)
		return
	}
}

func init() {
	Router.HandleFunc("/antlings", antlingPost).Methods("POST")
	Router.HandleFunc("/antlings", antlingGet).Methods("GET")
	Router.HandleFunc("/antlings/{id:[0-9]+}", antlingGetId).Methods("GET")
}
