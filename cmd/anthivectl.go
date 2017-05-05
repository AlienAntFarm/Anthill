package main

//go:generate go run include.go sql/*

import (
	"github.com/alienantfarm/anthive/assets"
	"github.com/alienantfarm/anthive/common"
	"github.com/alienantfarm/anthive/db"
)

func main() {
	common.Info.Printf("Init tables")
	_, err := db.Conn.Query(assets.Get("sql/init.sql"))
	if err != nil {
		common.Error.Fatalf("%s", err)
	}
}
