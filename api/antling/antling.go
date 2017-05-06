package antling

import (
	"github.com/alienantfarm/anthive/common"
	"github.com/alienantfarm/anthive/db"
	"net/http"
)

func Post(w http.ResponseWriter, r *http.Request) {
	a := &db.Antling{}
	a.Save()
	common.Encode(w, a)
}
