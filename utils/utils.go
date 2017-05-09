package utils

import (
	"encoding/json"
	"net/http"
)

func Encode(w http.ResponseWriter, i interface{}) error {
	err := json.NewEncoder(w).Encode(i)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Add("Content-Type", "application/json")
	}
	return err
}
