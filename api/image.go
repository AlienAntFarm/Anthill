package api

import (
	"github.com/alienantfarm/anthive/drivers/images"
	"github.com/alienantfarm/anthive/ext/db"
	"github.com/alienantfarm/anthive/utils"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
)

func imagePost(w http.ResponseWriter, r *http.Request) {
	ip := struct {
		Tag string `json:"tag"`
	}{}

	if err := utils.Decode(r, &ip); err != nil {
		glog.Errorf("%s", err)
		return
	}

	i, err := images.Docker2AIF(ip.Tag)

	if err != nil {
		glog.Errorf("%s", err)
		return
	}
	defer func() { images.RemoveOnFail(i.Archive, err) }()

	if err := ((*db.Image)(i)).Create(); err != nil {
		glog.Errorf("%s", err)
	} else if err := utils.Encode(w, i); err != nil {
		glog.Errorf("%s", err)
	}
}

func imageGet(w http.ResponseWriter, r *http.Request) {
	images := &db.Images{}
	if err := images.Get(); err != nil {
		glog.Errorf("%s", err)
	} else if err := utils.Encode(w, images); err != nil {
		glog.Errorf("%s", err)
	}
}

func imageGetId(w http.ResponseWriter, r *http.Request) {
	image := &db.Image{}

	if err := image.Get(mux.Vars(r)["id"]); err != nil {
		glog.Errorf("%s", err)
	} else if err := utils.Encode(w, image); err != nil {
		glog.Errorf("%s", err)
	}
}

func init() {
	Router.HandleFunc("/images", imagePost).Methods("POST")
	Router.HandleFunc("/images", imageGet).Methods("GET")
	Router.HandleFunc("/images/{id:[0-9]+}", imageGetId).Methods("GET")
}
