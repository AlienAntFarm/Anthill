package db

import (
	"database/sql"
	"fmt"
	"github.com/alienantfarm/anthive/utils"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"regexp"
)

var conn *sql.DB

func connect() *sql.DB {
	dbConfig := utils.Config.Database
	re := regexp.MustCompile("password=.* ")

	connString := fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%d sslmode=disable",
		dbConfig.Name, dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port,
	)

	glog.Infof("db connection: %s", re.ReplaceAllString(connString, "password=******* "))
	db, err := sql.Open("postgres", connString)

	if err != nil {
		glog.Fatalf(
			"Something bad happened during database connection: %s", err,
		)
	}
	return db
}

func Conn() *sql.DB {
	if conn == nil {
		conn = connect()
	}
	return conn
}
