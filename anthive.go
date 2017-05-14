package main

import (
	"fmt"
	"github.com/alienantfarm/anthive/api"
	"github.com/alienantfarm/anthive/utils"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"net/http"
	"time"
)

func run(cmd *cobra.Command, args []string) {
	api.InitScheduler()
	addr := fmt.Sprintf("%s:%d", utils.Config.Host, utils.Config.Port)
	router := api.Router
	fileHandler := http.FileServer(http.Dir(utils.Config.Assets.Static))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileHandler))
	s := &http.Server{
		Addr:           addr,
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	glog.Infof("running on http://%s/\n", addr)
	glog.Errorf("%s", s.ListenAndServe())
}

func main() {
	utils.Command.Run = run
	if err := utils.Command.Execute(); err != nil {
		glog.Fatalf("%s", err)
	}
}
