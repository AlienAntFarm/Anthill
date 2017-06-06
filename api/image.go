package api

import (
	"github.com/alienantfarm/anthive/db"
	"github.com/alienantfarm/anthive/drivers/images"
	"github.com/alienantfarm/anthive/utils"
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"net/http"
)

type ImagePost struct {
	Tag string `json:"tag"`
}

func imagePost(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		ip  *ImagePost = &ImagePost{}
	)
	if err := utils.Decode(r, ip); err != nil {
		glog.Errorf("%s", err)
		err.Dump(w)
		return
	}
	i, err := images.Docker2AIF(ip.Tag)
	if err != nil {
		glog.Errorf("%s", err)
		utils.NewError500(err).Dump(w)
		return
	}
	defer func() { images.RemoveOnFail(i.Archive, err) }()
	query := "INSERT INTO anthive.image (archive, command, environment, cwd, hostname) "
	query += "VALUES ($1, $2, $3, $4, $5) "
	query += "RETURNING anthive.image.id"

	args := []interface{}{i.Archive, pq.Array(i.Cmd), pq.Array(i.Env), i.Cwd, i.Hostname}
	if err = db.Conn().QueryRow(query, args...).Scan(&i.Id); err != nil {
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
	query := "SELECT i.id, i.archive, i.command, i.environment, i.cwd, i.hostname "
	query += "FROM anthive.image as i"

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
		args := []interface{}{
			&image.Id, &image.Archive, pq.Array(&image.Cmd), pq.Array(&image.Env),
			&image.Cwd, &image.Hostname,
		}
		if err := rows.Scan(args...); err != nil {
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
	image := &structs.Image{}

	query := "SELECT i.id, i.archive, i.command, i.environment, i.cwd, i.hostname "
	query += "FROM anthive.image as i "
	query += "WHERE i.id = $1"
	args := []interface{}{
		&image.Id, &image.Archive, pq.Array(&image.Cmd), pq.Array(&image.Env),
		&image.Cwd, &image.Hostname,
	}

	if err := db.Conn().QueryRow(query, id).Scan(args...); err != nil {
		glog.Errorln(err)
		utils.NewError500(err).Dump(w)
		return
	}
	if err := utils.Encode(w, image); err != nil {
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
