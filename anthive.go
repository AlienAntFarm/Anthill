package main

import (
	"fmt"
	"github.com/alienantfarm/anthive/common"
	"net/http"
	"time"
)

func index(w http.ResponseWriter, req *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	addr := fmt.Sprintf("%s:%d", common.Config.Host, common.Config.Port)
	s := &http.Server{
		Addr:           addr,
		Handler:        http.HandlerFunc(index),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	common.Info.Printf("Running on http://%s/\n", addr)
	common.Error.Fatal(s.ListenAndServe())
}
