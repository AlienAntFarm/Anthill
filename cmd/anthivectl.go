package main

//go:generate go run include.go sql/*

import (
	"github.com/alienantfarm/anthive/assets"
	"github.com/alienantfarm/anthive/utils"
	"github.com/alienantfarm/anthive/db"
	"github.com/spf13/cobra"
)

func runAsset(assetName string) {
	asset := assets.Get(assetName)
	utils.Info.Printf("\n%s", asset)
	_, err := db.Conn.Query(asset)
	if err != nil {
		utils.Error.Fatalf("%s", err)
	}
}

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
		runAsset("sql/init.sql")
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean every tables from anthive database",
	Run: func(cmd *cobra.Command, args []string) {
		runAsset("sql/clean.sql")
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Recreate every tables from anthive database",
	Run: func(cmd *cobra.Command, args []string) {
		runAsset("sql/remove.sql")
		runAsset("sql/init.sql")
	},
}

func main() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(resetCmd)
	if err := rootCmd.Execute(); err != nil {
		utils.Error.Fatalf("%s", err)
	}
}
