package db

import (
	"database/sql"
	"github.com/alienantfarm/anthive/utils/structs"
)

type Antling structs.Antling

func (a *Antling) Create() error {
	query := "INSERT INTO anthive.antling "
	query += "DEFAULT VALUES "
	query += "RETURNING anthive.antling.id"

	return Conn().QueryRow(query).Scan(&a.Id)
}

func (a *Antling) Get(id string) error {
	query := "SELECT anthive.antling.id "
	query += "FROM anthive.antling "
	query += "WHERE anthive.antling.id = $1"

	return Conn().QueryRow(query, id).Scan(&a.Id)
}

type Antlings structs.Antlings

func (a *Antlings) Get() (err error) {
	var rows *sql.Rows

	query := "SELECT anthive.antling.id "
	query += "FROM anthive.antling"

	if rows, err = Conn().Query(query); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		antling := &structs.Antling{}

		if err = rows.Scan(&antling.Id); err != nil {
			return
		}
		a.Antlings = append(a.Antlings, antling)
	}

	return
}
