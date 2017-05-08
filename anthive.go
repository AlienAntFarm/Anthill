package main

import (
	"fmt"
	"github.com/alienantfarm/anthive/api"
	"github.com/alienantfarm/anthive/utils"
	"net/http"
	"time"
)

func main() {
	addr := fmt.Sprintf("%s:%d", utils.Config.Host, utils.Config.Port)

	s := &http.Server{
		Addr:           addr,
		Handler:        api.Router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	utils.Info.Printf("running on http://%s/\n", addr)
	utils.Error.Fatal(s.ListenAndServe())
}
