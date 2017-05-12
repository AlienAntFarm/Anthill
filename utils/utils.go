package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	MIME_JSON = "application/json"
)

type UnexpectedContentType struct {
	Expected string
	Received string
}

func (uct *UnexpectedContentType) Error() string {
	return fmt.Sprintf(
		"unexpected content type, want: %s, got: %s", uct.Expected, uct.Received,
	)
}

type UnmatchingIds struct {
	Ids [2]int
}

func (ui *UnmatchingIds) Error() string {
	return fmt.Sprintf(
		"ids are not matching got: %d and %d", ui.Ids[0], ui.Ids[1],
	)
}

func Encode(w http.ResponseWriter, i interface{}) error {
	err := json.NewEncoder(w).Encode(i)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Add("Content-Type", MIME_JSON)
	}
	return err
}

func Decode(r *http.Request, i interface{}) error {
	mime := r.Header.Get("Content-Type")
	if mime != MIME_JSON {
		return &UnexpectedContentType{MIME_JSON, mime}
	}
	return json.NewDecoder(r.Body).Decode(i)
}
