package rest

import (
	"encoding/json"
	"fmt"
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

// ERROR reports the error back to the user withing a JSON format
func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
	} else {
		JSON(w, statusCode, err)
	}
}
