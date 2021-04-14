package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/uid_provider/api/models"
	"gitlab.inspr.dev/inspr/core/cmd/uid_provider/client"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// Handler is a structure which cointains methods to handle
// requests received by the UID Provider API
type Handler struct {
	rdb client.RedisManager
	ctx context.Context
}

// NewHandler instantiates a new Handler structure
func NewHandler(rdb client.RedisManager, ctx context.Context) *Handler {
	return &Handler{
		rdb: rdb,
		ctx: ctx,
	}
}

// CreateUserHandler handles user creation requests
func (h *Handler) CreateUserHandler() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		data := models.ReceivedDataCreate{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			rest.ERROR(w, err)
			return
		}

		if err := h.rdb.CreateUser(h.ctx, data.UID, data.Password, data.User); err != nil {
			rest.ERROR(w, err)
			return
		}
	}).Post().JSON().Recover()
}

// DeleteUserHandler handles user deletion requests
func (h *Handler) DeleteUserHandler() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		data := models.ReceivedDataDelete{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			rest.ERROR(w, err)
			return
		}

		if err := h.rdb.DeleteUser(h.ctx, data.UID, data.Password, data.UserToBeDeleted); err != nil {
			rest.ERROR(w, err)
			return
		}
	}).Post().JSON().Recover()
}

// UpdatePasswordHandler handles requests to update an user password
func (h *Handler) UpdatePasswordHandler() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		data := models.ReceivedDataUpdate{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			rest.ERROR(w, err)
			return
		}

		if err := h.rdb.UpdatePassword(h.ctx, data.UID, data.Password,
			data.UserToBeUpdated, data.NewPassword); err != nil {

			rest.ERROR(w, err)
			return
		}
	}).Put().JSON().Recover()
}

// LoginHandler handles login requests
func (h *Handler) LoginHandler() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		data := models.ReceivedDataLogin{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			rest.ERROR(w, err)
			return
		}

		token, err := h.rdb.Login(h.ctx, data.UID, data.Password)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		rest.JSON(w, 200, token)
	}).Post().JSON().Recover()
}

// RefreshTokenHandler handles token refresh requests
func (h *Handler) RefreshTokenHandler() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		data := models.ReceivedDataRefresh{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			rest.ERROR(w, err)
			return
		}

		payload, err := h.rdb.RefreshToken(h.ctx, data.RefreshToken)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		rest.JSON(w, 200, payload)
	}).Post().JSON().Recover()
}
