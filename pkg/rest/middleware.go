package rest

import "net/http"

// JSON specifies in the header that the response content is a json
func (h Handler) JSON() Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h(w, r)
	}
}

type CRUDHandler interface {
	HandleCreate() Handler
	HandleDelete() Handler
	HandleUpdate() Handler
	HandleGet() Handler
}

func HandleCRUD(handler CRUDHandler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			handler.HandleGet().JSON().Recover()(w, r)

		case http.MethodPost:
			handler.HandleCreate().JSON().Recover()(w, r)

		case http.MethodPut:
			handler.HandleUpdate().JSON().Recover()(w, r)

		case http.MethodDelete:
			handler.HandleDelete().JSON().Recover()(w, r)

		default:
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}
