package db

import (
	"database/sql"
	"fmt"
	"github.com/alienantfarm/anthive/common"
	_ "github.com/lib/pq"
	"regexp"
)

var (
	Conn = connect()
)

func connect() *sql.DB {
	dbConfig := common.Config.Database
	re := regexp.MustCompile("password=.* ")

	connString := fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%d sslmode=disable",
		dbConfig.Name, dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port,
	)

	common.Info.Printf(
		"Connecting to database with following options: %s",
		re.ReplaceAllString(connString, "password=******* "),
	)
	db, err := sql.Open("postgres", connString)

	if err != nil {
		common.Error.Fatalf(
			"Something bad happened during database connection: %s", err,
		)
	}
	return db
}
