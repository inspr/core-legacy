package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"inspr.dev/inspr/cmd/uid_provider/api/models"
	"inspr.dev/inspr/cmd/uid_provider/client"
	"inspr.dev/inspr/pkg/rest"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = zap.NewProduction(zap.Fields(zap.String("section", "uidp-api-handlers")))
}

// Handler is a structure which cointains methods to handle
// requests received by the UID Provider API
type Handler struct {
	ctx context.Context
	rdb client.RedisManager
}

// NewHandler instantiates a new Handler structure
func NewHandler(ctx context.Context, rdb client.RedisManager) *Handler {
	logger.Info("creating a new UIDP API handler")

	return &Handler{
		rdb: rdb,
		ctx: ctx,
	}
}

// CreateUserHandler handles user creation requests
func (h *Handler) CreateUserHandler() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("creating a new user")

		data := models.ReceivedDataCreate{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			rest.ERROR(w, err)
			return
		}

		if err := h.rdb.CreateUser(h.ctx, data.UID, data.Password, data.User); err != nil {
			rest.ERROR(w, err)
			return
		}

		logger.Info("user created", zap.String("username", data.UID))
	}).Post().JSON().Recover()
}

// DeleteUserHandler handles user deletion requests
func (h *Handler) DeleteUserHandler() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("deleting a user")

		data := models.ReceivedDataDelete{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			rest.ERROR(w, err)
			return
		}

		if err := h.rdb.DeleteUser(h.ctx, data.UID, data.Password, data.UserToBeDeleted); err != nil {
			rest.ERROR(w, err)
			return
		}

		logger.Info("user deleted", zap.String("username", data.UID))
	}).Delete().JSON().Recover()
}

// UpdatePasswordHandler handles requests to update an user password
func (h *Handler) UpdatePasswordHandler() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("updating a user")

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

		logger.Info("user updated", zap.String("username", data.UID))
	}).Put().JSON().Recover()
}

// LoginHandler handles login requests
func (h *Handler) LoginHandler() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("login of an user")

		data := models.ReceivedDataLogin{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			rest.ERROR(w, err)
			return
		}

		logger.Info("username provided", zap.String("UID", data.UID))

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
		logger.Info("refreshing the token")

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
