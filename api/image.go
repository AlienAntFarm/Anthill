package api

import (
	"github.com/alienantfarm/anthive/db"
	"github.com/alienantfarm/anthive/drivers/images"
	"github.com/alienantfarm/anthive/utils"
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
)

type ImagePost struct {
	Tag string `json:"tag"`
}

func imagePost(w http.ResponseWriter, r *http.Request) {
	ip := &ImagePost{}
	if err := utils.Decode(r, ip); err != nil {
		glog.Errorf("%s", err)
		err.Dump(w)
		return
	}
	archive, err := images.Docker2AIF(ip.Tag)
	if err != nil {
		glog.Errorf("%s", err)
		utils.NewError500(err).Dump(w)
		return
	}
	query := "INSERT INTO anthive.image (file)"
	query += "VALUES ($1) "
	query += "RETURNING anthive.image.id"

	i := &structs.Image{Archive: archive}
	if err := db.Conn().QueryRow(query).Scan(&i.Id); err != nil {
		glog.Errorf("%s", err)
		utils.NewError500(err).Dump(w)
		return
	}

	if err := utils.Encode(w, i); err != nil {
		glog.Errorln(err)
		err.Dump(w)
		return
	}
}

func imageGet(w http.ResponseWriter, r *http.Request) {
	query := "SELECT anthive.image.id, anthive.image.archive "
	query += "FROM anthive.image"

	rows, err := db.Conn().Query(query)
	if err != nil {
		glog.Errorln(err)
		utils.NewError500(err).Dump(w)
		return
	}
	defer rows.Close()
	images := []*structs.Image{}
	for rows.Next() {
		image := &structs.Image{}

		if err := rows.Scan(&image.Id, &image.Archive); err != nil {
			glog.Errorln(err)
			return
		}
		images = append(images, image)
	}
	if err := utils.Encode(w, structs.Images{images}); err != nil {
		glog.Errorln(err)
		err.Dump(w)
		return
	}
}

func imageGetId(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	i := &structs.Image{}

	query := "SELECT anthive.image.id, anthive.image.archive "
	query += "FROM anthive.image "
	query += "WHERE anthive.image.id = $1"

	if err := db.Conn().QueryRow(query, id).Scan(&i.Id, &i.Archive); err != nil {
		glog.Errorln(err)
		utils.NewError500(err).Dump(w)
		return
	}
	if err := utils.Encode(w, i); err != nil {
		glog.Errorln(err)
		err.Dump(w)
		return
	}
}

func init() {
	Router.HandleFunc("/images", imagePost).Methods("POST")
	Router.HandleFunc("/images", imageGet).Methods("GET")
	Router.HandleFunc("/images/{id:[0-9]+}", imageGetId).Methods("GET")
}
