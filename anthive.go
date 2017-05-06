package main

import (
	"fmt"
	"github.com/alienantfarm/anthive/api/antling"
	"github.com/alienantfarm/anthive/common"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/antling", antling.Post).Methods("POST")

	s := &http.Server{
		Addr:           addr,
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	common.Info.Printf("Running on http://%s/\n", addr)
	common.Error.Fatal(s.ListenAndServe())
}
