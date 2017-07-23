package api

func incubatorGet(w http.ResponseWriter, r *http.Request) {
}

func init() {
	Router.HandleFunc("/incubators", incubatorGet).Methods("GET")
}
