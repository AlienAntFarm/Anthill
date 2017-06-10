package main

import (
	"fmt"
	"github.com/alienantfarm/anthive/api"
	"github.com/alienantfarm/anthive/drivers/images"
	"github.com/alienantfarm/anthive/ext/minio"
	"github.com/alienantfarm/anthive/utils"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"net/http"
	"time"
)

const IMAGES_PREFIX = "/static/images"

func run(cmd *cobra.Command, args []string) {
	api.InitScheduler()
	addr := fmt.Sprintf("%s:%d", utils.Config.Host, utils.Config.Port)
	router := api.Router
	client := minio.Client()
	glog.Infof("%q", client)
	// serve images
	imagesHandler := http.FileServer(http.Dir(utils.Config.Assets.Images))
	router.PathPrefix(IMAGES_PREFIX).Handler(
		http.StripPrefix(IMAGES_PREFIX, imagesHandler),
	)

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

func runDocker2AIF(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		glog.Fatalf("you must specify a docker tag to build the image upon")
	}
	_, err := images.Docker2AIF(args[0])
	if err != nil {
		glog.Fatalf("%s", err)
	}
}

func main() {
	utils.Command.Run = run
	utils.OCICommand.Run = runDocker2AIF

	if err := utils.Command.Execute(); err != nil {
		glog.Fatalf("%s", err)
	}
}
