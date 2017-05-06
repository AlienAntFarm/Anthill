package api

import (
	"github.com/alienantfarm/anthive/common"
	"github.com/alienantfarm/anthive/db"
	"net/http"
)

type Antling struct {
	Id int `json:"id"`
}

func antlingPost(w http.ResponseWriter, r *http.Request) {
	query := `
    INSERT INTO anthive.antling
    DEFAULT VALUES
    RETURNING anthive.antling.id
  `
	a := &Antling{}
	db.Conn.QueryRow(query).Scan(&a.Id)
	common.Encode(w, a)
}

func init() {
	Router.HandleFunc("/antling", antlingPost).Methods("POST")
}
