package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"

	"inspr.dev/inspr/pkg/ierrors"
)

// AttachProfiler is responsible for adding the pprof routes to the server mux
// passed as a parameter
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
	default: // default case
		JSON(w, http.StatusInternalServerError, err)
	}
}

// UnmarshalERROR generates an ierror error with the response body created by
// the ERROR function.
//
// Note that this function will always return an error no
// matter if it found the error on the body or not, that means that if the
// error details cannot be found on the body of the response body it will
// generate a unkown error
func UnmarshalERROR(body io.Reader) error {
	err := defaultErr
	decoder := json.NewDecoder(body)
	decoder.Decode(&err)
	return err
}
