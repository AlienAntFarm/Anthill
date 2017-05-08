package utils

import (
	"encoding/json"
	"net/http"
)

func Encode(w http.ResponseWriter, i interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(w)
	return e.Encode(i)
}
