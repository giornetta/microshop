package respond

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// JSON serializes the given data as JSON and sends it as a http Response
func JSON(w http.ResponseWriter, status int, v interface{}) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		Err(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func Err(w http.ResponseWriter, status int, err error) {
	JSON(w, status, map[string]string{
		"error": err.Error(),
	})
}
