package api

import (
	"github.com/alienantfarm/anthive/common"
	"github.com/alienantfarm/anthive/db"
	"net/http"
)

type Antling struct {
	Id int `json:"id"`
}

type Antlings struct {
	Antlings []*Antling `json:"antlings"`
}

func antlingPost(w http.ResponseWriter, r *http.Request) {
	var err error
	query := "INSERT INTO anthive.antling "
	query += "DEFAULT VALUES "
	query += "RETURNING anthive.antling.id"

	a := &Antling{}
	err = db.Conn.QueryRow(query).Scan(&a.Id)
	if err != nil {
		common.Error.Println(err)
		return
	}
	err = common.Encode(w, a)
	if err != nil {
		common.Error.Println(err)
		return
	}
}

func antlingGet(w http.ResponseWriter, r *http.Request) {
	var err error
	query := "SELECT anthive.antling.id "
	query += "FROM anthive.antling"
	rows, err := db.Conn.Query(query)
	if err != nil {
		common.Error.Println(err)
		return
	}
	defer rows.Close()
	antlings := []*Antling{}
	for rows.Next() {
		antling := &Antling{}
		err = rows.Scan(&antling.Id)
		if err != nil {
			common.Error.Println(err)
			return
		}
		antlings = append(antlings, antling)
	}
	err = common.Encode(w, Antlings{antlings})
	if err != nil {
		common.Error.Println(err)
		return
	}
}

func init() {
	Router.HandleFunc("/antlings", antlingPost).Methods("POST")
	Router.HandleFunc("/antlings", antlingGet).Methods("GET")
}
