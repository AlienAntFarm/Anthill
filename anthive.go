package main

import (
	"flag"
	"fmt"
	"github.com/alienantfarm/anthive/api"
	"github.com/alienantfarm/anthive/utils"
	"github.com/golang/glog"
	"net/http"
	"os"
	"time"
)

func main() {
	// reinit args for glog
	os.Args = os.Args[:1]
	conf := utils.Config()

	flag.Set("logtostderr", "true")
	if conf.Debug {
		flag.Set("v", "10") // totally arbitrary but who cares!
	}
	flag.Parse()
	glog.V(1).Infoln("debug mode enabled")
	api.InitScheduler()

	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	s := &http.Server{
		Addr:           addr,
		Handler:        api.Router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}

	glog.Infof("running on http://%s/\n", addr)
	glog.Errorf("%s", s.ListenAndServe())
}
