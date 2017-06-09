package db

import (
	"database/sql"
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/lib/pq"
)

type Job structs.Job

func (j *Job) UpdateAntling() error {
	query := "UPDATE anthive.job "
	query += "SET fk_antling = $1"
	query += "WHERE anthive.job.id = $2"

	update := Conn().QueryRow(query, j.IdAntling, j.Id).Scan // quick alias

	if err := update(); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func (j *Job) UpdateState() error {
	query := "UPDATE anthive.job "
	query += "SET state = $1 "
	query += "WHERE anthive.job.id = $2"

	update := Conn().QueryRow(query, j.State, j.Id).Scan // quick alias
	if err := update(); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func (j *Job) Create(imageId int) error {
	query := "WITH j as ("
	query += "  INSERT INTO anthive.job (fk_image, command, environment, cwd) "
	query += "  VALUES ($1, $2, $3, $4) RETURNING anthive.job.id "
	query += ") "
	query += "SELECT j.id, i.id, i.archive, i.command, i.environment, i.cwd, i.hostname "
	query += "FROM j, anthive.image as i "
	query += "WHERE i.id = $1"

	i := &j.Image // extract image sub struct
	args := []interface{}{imageId, pq.Array(j.Cmd), pq.Array(j.Env), j.Cwd}
	argsScan := []interface{}{
		&j.Id, &i.Id, &i.Archive, pq.Array(&i.Cmd), pq.Array(&i.Env), &i.Cwd, &i.Hostname,
	}

	return Conn().QueryRow(query, args...).Scan(argsScan...)
}

func (j *Job) Get(id string) (err error) {

	query := "SELECT j.id, j.state, j.cwd, j.command, j.environment, "
	query += "  i.id, i.archive, i.command, i.environment, i.cwd, i.hostname "
	query += "FROM anthive.job AS j, anthive.image AS i "
	query += "WHERE i.id = j.fk_image AND j.id = $1"

	i := &j.Image
	args := []interface{}{
		&j.Id, &j.State, &j.Cwd, pq.Array(&j.Cmd), pq.Array(&j.Env),
		&i.Id, &i.Archive, pq.Array(&i.Cmd), pq.Array(&i.Env), &i.Cwd, &i.Hostname,
	}
	return Conn().QueryRow(query, id).Scan(args...)
}

type Jobs structs.Jobs

func (jobs *Jobs) Get(js structs.JobState) (err error) {
	var rows *sql.Rows

	query := "SELECT j.id, j.state, j.cwd, j.command, j.environment, j.fk_antling, "
	query += "  i.id, i.archive, i.command, i.environment, i.cwd, i.hostname "
	query += "FROM anthive.job AS j, anthive.image AS i "
	query += "WHERE j.fk_image = i.id AND j.state <= $1"

	if rows, err = Conn().Query(query, js); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		j := &structs.Job{}
		i := &j.Image

		args := []interface{}{
			&j.Id, &j.State, &j.Cwd, pq.Array(&j.Cmd), pq.Array(&i.Env), &j.IdAntling,
			&i.Id, &i.Archive, pq.Array(&i.Cmd), pq.Array(&i.Env), &i.Cwd, &i.Hostname,
		}
		if err = rows.Scan(args...); err != nil {
			return
		}
		jobs.Jobs = append(jobs.Jobs, j)
	}
	return
}
