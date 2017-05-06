package main

import (
	"fmt"
	"github.com/alienantfarm/anthive/api"
	"github.com/alienantfarm/anthive/common"
	"net/http"
	"time"
)

func main() {
	addr := fmt.Sprintf("%s:%d", common.Config.Host, common.Config.Port)

	s := &http.Server{
		Addr:           addr,
		Handler:        api.Router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	common.Info.Printf("Running on http://%s/\n", addr)
	common.Error.Fatal(s.ListenAndServe())
}
