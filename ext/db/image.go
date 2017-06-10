package db

import (
	"database/sql"
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/lib/pq"
)

type Image structs.Image

func (i *Image) Create() error {
	query := "INSERT INTO anthive.image (archive, command, environment, cwd, hostname) "
	query += "VALUES ($1, $2, $3, $4, $5) "
	query += "RETURNING anthive.image.id"

	args := []interface{}{
		i.Archive, pq.Array(i.Cmd), pq.Array(i.Env), i.Cwd, i.Hostname,
	}

	return Client().QueryRow(query, args...).Scan(&i.Id)
}

func (i *Image) Get(id string) error {
	query := "SELECT i.id, i.archive, i.command, i.environment, i.cwd, i.hostname "
	query += "FROM anthive.image as i "
	query += "WHERE i.id = $1"

	args := []interface{}{
		&i.Id, &i.Archive, pq.Array(&i.Cmd), pq.Array(&i.Env), &i.Cwd, &i.Hostname,
	}

	return Client().QueryRow(query, id).Scan(args...)
}

type Images structs.Images

func (i *Images) Get() (err error) {
	var rows *sql.Rows

	query := "SELECT i.id, i.archive, i.command, i.environment, i.cwd, i.hostname "
	query += "FROM anthive.image as i"

	if rows, err = Client().Query(query); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		image := &structs.Image{}
		args := []interface{}{
			&image.Id, &image.Archive, pq.Array(&image.Cmd), pq.Array(&image.Env),
			&image.Cwd, &image.Hostname,
		}
		if err = rows.Scan(args...); err != nil {
			return
		}
		i.Images = append(i.Images, image)
	}
	return
}
