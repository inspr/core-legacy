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
			JSON(w, http.StatusConflict, e)
		case ierrors.InternalServer:
			JSON(w, http.StatusInternalServerError, e)
		case ierrors.InvalidName:
			JSON(w, http.StatusForbidden, e)
		case ierrors.InvalidApp:
			JSON(w, http.StatusForbidden, e)
		case ierrors.InvalidChannel:
			JSON(w, http.StatusForbidden, e)
		case ierrors.InvalidChannelType:
			JSON(w, http.StatusForbidden, e)
		case ierrors.BadRequest:
			JSON(w, http.StatusBadRequest, e)
		default:
			JSON(w, http.StatusInternalServerError, e)
		}

	// default case
	case error:
		defaultInsprErr := ierrors.InsprError{
			Message: e.Error(),
			Code:    http.StatusInternalServerError,
		}
		JSON(w, http.StatusInternalServerError, defaultInsprErr)
	}
}
