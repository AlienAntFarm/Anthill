package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
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

func (uct *UnexpectedContentType) Dump(w http.ResponseWriter) {
	http.Error(w, uct.Error(), http.StatusBadRequest)
}

type HttpErrorDumper interface {
	error
	Dump(w http.ResponseWriter)
}

type HttpError struct {
	StatusCode int
	error
}

func (he *HttpError) Dump(w http.ResponseWriter) {
	http.Error(w, he.Error(), he.StatusCode)
}

func NewError500(err error) *HttpError {
	return &HttpError{http.StatusInternalServerError, err}
}

type UnmatchingIds struct {
	Ids [2]int
}

func (ui *UnmatchingIds) Error() string {
	return fmt.Sprintf(
		"ids are not matching got: %d and %d", ui.Ids[0], ui.Ids[1],
	)
}

func (ui *UnmatchingIds) Dump(w http.ResponseWriter) {
	http.Error(w, ui.Error(), http.StatusBadRequest)
}
func Encode(w http.ResponseWriter, i interface{}) HttpErrorDumper {
	err := json.NewEncoder(w).Encode(i)
	if err != nil {
		return NewError500(err)
	} else {
		w.Header().Add("Content-Type", MIME_JSON)
		return nil
	}
}

func Decode(r *http.Request, i interface{}) HttpErrorDumper {
	mime := r.Header.Get("Content-Type")
	if mime != MIME_JSON {
		return &UnexpectedContentType{MIME_JSON, mime}
	} else if err := json.NewDecoder(r.Body).Decode(i); err != nil {
		return NewError500(err)
	} else {
		return nil
	}
}

// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func SecureRandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
