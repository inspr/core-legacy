package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

// RecoverFromPanic will handle panic
func RecoverFromPanic(w http.ResponseWriter) {
	if recoveryMessage := recover(); recoveryMessage != nil {
		ERROR(w, ierrors.NewError().InternalServer().Message("%s", recoveryMessage).Build())
	}
}

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

// UnmarshalERROR generates a golang error with the
// response body created by the ERROR function
func UnmarshalERROR(r io.Reader) error {
	errBody := struct {
		Error string `json:"error"`
	}{}
	decoder := json.NewDecoder(r)
	decoder.Decode(&errBody)
	return errors.New(errBody.Error)
}
