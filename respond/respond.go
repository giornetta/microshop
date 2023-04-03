package respond

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/giornetta/microshop/errors"
)

// JSON serializes the given data as JSON and sends it as a HTTP Response
func JSON(w http.ResponseWriter, status int, v interface{}) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		Err(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

// Err serializes the given error as JSON and sends it as a HTTP Response,
// choosing the appropriate status code.
func Err(w http.ResponseWriter, err error) {
	var status int = http.StatusInternalServerError
	var msg string

	if e, ok := err.(errors.WithStatusCode); ok {
		status = e.StatusCode()
		msg = e.Error()
	} else {
		msg = new(errors.ErrInternal).Error()
	}

	JSON(w, status, map[string]string{
		"error": msg,
	})
}
