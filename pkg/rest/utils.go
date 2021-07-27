package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"inspr.dev/inspr/pkg/ierrors"
)

// RecoverFromPanic will handle panic
func RecoverFromPanic(w http.ResponseWriter) {
	if recoveryMessage := recover(); recoveryMessage != nil {
		ERROR(w, ierrors.New("%s", recoveryMessage).InternalServer())
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

// TODO REVIEW

// ERROR reports the error back to the user within a JSON format
func ERROR(w http.ResponseWriter, err error) {
	switch ierrors.Code(err) {
	case ierrors.AlreadyExists:
		JSON(w, http.StatusConflict, err)
	case ierrors.NotFound:
		JSON(w, http.StatusNotFound, err)
	case ierrors.InternalServer:
		JSON(w, http.StatusInternalServerError, err)
	case ierrors.InvalidName:
		JSON(w, http.StatusForbidden, err)
	case ierrors.InvalidApp:
		JSON(w, http.StatusForbidden, err)
	case ierrors.InvalidChannel:
		JSON(w, http.StatusForbidden, err)
	case ierrors.InvalidType:
		JSON(w, http.StatusForbidden, err)
	case ierrors.BadRequest:
		JSON(w, http.StatusBadRequest, err)
	case ierrors.Unauthorized:
		JSON(w, http.StatusUnauthorized, err)
	case ierrors.Forbidden:
		JSON(w, http.StatusForbidden, err)
	// default case
	default:
		JSON(w, http.StatusInternalServerError, err)
	}
}

// UnmarshalERROR generates a golang error with the
// response body created by the ERROR function
func UnmarshalERROR(body io.Reader) error {
	defaultErr := ierrors.New("cannot retrieve error from server").
		InternalServer()

	err := ierrors.New("")
	decoder := json.NewDecoder(body)
	decoder.Decode(&err)

	if err != nil {
		return err
	}
	return defaultErr
}
