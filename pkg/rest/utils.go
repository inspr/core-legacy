package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
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
func ERROR(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *ierrors.InsprError:
		switch e.Code {
		case ierrors.NotFound:
			JSON(w, http.StatusNotFound, e)
		case ierrors.AlreadyExists:

		case ierrors.InternalServer:

		case ierrors.InvalidName:

		case ierrors.InvalidChannel:

		case ierrors.InvalidChannelType:

		case ierrors.BadRequest:

		}
	default:
		JSON(w, http.StatusInternalServerError, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
	}
}
