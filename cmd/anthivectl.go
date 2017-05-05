package main

import (
	"github.com/alienantfarm/anthive/common"
	"github.com/alienantfarm/anthive/db"
)

const table_init = `
CREATE SCHEMA IF NOT EXISTS anthive;

-- Set default search_path to schema
SET search_path TO anthive,public;

-- Creation of tables
CREATE TABLE IF NOT EXISTS antling (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50)
)
`

func main() {
	common.Info.Printf("Init tables")
	_, err := db.Conn.Query(table_init)
	if err != nil {
		common.Error.Fatalf("%s", err)
	}
}
