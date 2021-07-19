package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"

	"inspr.dev/inspr/pkg/ierrors"
)

<<<<<<< HEAD
// AttachProfiler is responsible for adding the pprof routes to the server mux
// passed as a parameter
=======
>>>>>>> 0a33d610 (dev(servers): added the route for pprof in all inspr services, still missing for pods/dapps created)
func AttachProfiler(m *http.ServeMux) {
	m.HandleFunc("/debug/pprof/", pprof.Index)
	m.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	m.HandleFunc("/debug/pprof/profile", pprof.Profile)
	m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	m.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	m.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	m.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	m.Handle("/debug/pprof/block", pprof.Handler("block"))
}

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

// ERROR reports the error back to the user within a JSON format
func ERROR(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *ierrors.InsprError:
		switch e.Code {
		case ierrors.AlreadyExists:
			JSON(w, http.StatusConflict, e)
		case ierrors.NotFound:
			JSON(w, http.StatusNotFound, e)
		case ierrors.InternalServer:
			JSON(w, http.StatusInternalServerError, e)
		case ierrors.InvalidName:
			JSON(w, http.StatusForbidden, e)
		case ierrors.InvalidApp:
			JSON(w, http.StatusForbidden, e)
		case ierrors.InvalidChannel:
			JSON(w, http.StatusForbidden, e)
		case ierrors.InvalidType:
			JSON(w, http.StatusForbidden, e)
		case ierrors.BadRequest:
			JSON(w, http.StatusBadRequest, e)
		case ierrors.Unauthorized:
			JSON(w, http.StatusUnauthorized, e)
		case ierrors.Forbidden:
			JSON(w, http.StatusForbidden, e)
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
func UnmarshalERROR(body io.Reader) error {
	defaultErr := ierrors.
		NewError().
		InternalServer().
		Message("cannot retrieve error from server").
		Build()

	var err *ierrors.InsprError
	decoder := json.NewDecoder(body)
	decoder.Decode(&err)

	if err != nil {
		return err
	}
	return defaultErr
}
