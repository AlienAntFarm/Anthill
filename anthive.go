package main

import (
	"flag"
	"fmt"
	"github.com/alienantfarm/anthive/api"
	"github.com/alienantfarm/anthive/utils"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "anthive",
	Short: "Start anthive server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// reinit args for glog
		os.Args = os.Args[:1]
		flag.Set("logtostderr", "true")
		flag.Parse()
	},
	Run: func(cmd *cobra.Command, args []string) {
		conf := utils.Config()
		if conf.Debug {
			flag.Set("v", "10") // totally arbitrary but who cares!
			flag.Parse()
		}
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
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		glog.Fatalf("%s", err)
	}
}
