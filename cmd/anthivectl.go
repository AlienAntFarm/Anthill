package main

//go:generate go run include.go sql/*

import (
	"github.com/alienantfarm/anthive/assets"
	"github.com/alienantfarm/anthive/common"
	"github.com/alienantfarm/anthive/db"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "anthivectl",
	Short: "Simple cli to deal with various part of anthive",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init tables and types for anthive database",
	Run: func(cmd *cobra.Command, args []string) {
		common.Info.Printf("Init tables")
		_, err := db.Conn.Query(assets.Get("sql/init.sql"))
		if err != nil {
			common.Error.Fatalf("%s", err)
		}
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean every tables from anthive database",
	Run: func(cmd *cobra.Command, args []string) {
		common.Info.Printf("Clean tables")
		_, err := db.Conn.Query(assets.Get("sql/clean.sql"))
		if err != nil {
			common.Error.Fatalf("%s", err)
		}
	},
}

func main() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(cleanCmd)
	if err := rootCmd.Execute(); err != nil {
		common.Error.Fatalf("%s", err)
	}
}
