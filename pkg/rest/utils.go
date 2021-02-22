package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// JSON writes the data into the response writer with a JSON format
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

// ERROR reports the error back to the user within a JSON format
func ERROR(w http.ResponseWriter, statusCode int, err error) {
	JSON(w, statusCode, err)
}

// UnmarshalERROR generates a golang erro with the
// response body created by the ERROR function
func UnmarshalERROR(r io.Reader) error {
	errBody := struct {
		Error string `json:"error"`
	}{}
	decoder := json.NewDecoder(r)
	decoder.Decode(&errBody)
	return errors.New(errBody.Error)

}
