package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"inspr.dev/inspr/cmd/uidp/api/models"
	"inspr.dev/inspr/cmd/uidp/client"
	"inspr.dev/inspr/pkg/logs"
	"inspr.dev/inspr/pkg/rest"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = logs.Logger(zap.Fields(zap.String("section", "uidp-api-handlers")))
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
	l := logger.With(zap.String("subSection", "users"), zap.String("operation", "create"))
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		l = l.With(zap.String("host", r.Host))
		l.Info("received create user request")

		data := models.ReceivedDataCreate{}
		l.Debug("reading request body")
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			l.Error("error reading body", zap.Error(err))
			rest.ERROR(w, err)
			return
		}
		l.With(zap.String("user-creator", data.UID), zap.String("user-created", data.User.UID))
		l.Debug("creating user in redis")
		if err := h.rdb.CreateUser(h.ctx, data.UID, data.Password, data.User); err != nil {
			l.Error("error creating user in redis", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		l.Info("user created")
	}).Post().JSON().Recover()
}

// DeleteUserHandler handles user deletion requests
func (h *Handler) DeleteUserHandler() rest.Handler {
	l := logger.With(zap.String("subSection", "users"), zap.String("operation", "delete"))
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		l = l.With(zap.String("host", r.Host))
		l.Info("received delete user request")

		data := models.ReceivedDataDelete{}
		l.Debug("reading request body")
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			l.Error("error reading body", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		l = l.With(zap.String("user-deleter", data.UID), zap.String("user-deleted", data.UserToBeDeleted))

		l.Debug("deleting user in redis")
		if err := h.rdb.DeleteUser(h.ctx, data.UID, data.Password, data.UserToBeDeleted); err != nil {
			l.Error("error deleting user in redis", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		l.Info("user deleted")
	}).Delete().JSON().Recover()
}

// UpdatePasswordHandler handles requests to update an user password
func (h *Handler) UpdatePasswordHandler() rest.Handler {
	l := logger.With(zap.String("subSection", "users"), zap.String("operation", "update"))
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		l = l.With(zap.String("host", r.Host))
		l.Info("received update request")

		data := models.ReceivedDataUpdate{}
		l.Debug("reading request body")
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			l.Error("error reading body", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		l = l.With(zap.String("user-updater", data.UID), zap.String("user-updated", data.UID))

		l.Debug("updating user in redis")
		if err := h.rdb.UpdatePassword(h.ctx, data.UID, data.Password,
			data.UserToBeUpdated, data.NewPassword); err != nil {
			l.Error("error updating user in redis", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		l.Info("user updated")
	}).Put().JSON().Recover()
}

// LoginHandler handles login requests
func (h *Handler) LoginHandler() rest.Handler {
	l := logger.With(zap.String("subSection", "users"), zap.String("operation", "login"))
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		l = l.With(zap.String("host", r.Host))
		l.Info("received login request")

		data := models.ReceivedDataLogin{}
		l.Debug("reading request body")
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			l.Error("error reading body", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		l = l.With(zap.String("user", data.UID))

		l.Debug("logging in user username provided")
		token, err := h.rdb.Login(h.ctx, data.UID, data.Password)
		if err != nil {
			l.Debug("error logging in")
			rest.ERROR(w, err)
			return
		}

		rest.JSON(w, 200, token)
	}).Post().JSON().Recover()
}

// RefreshTokenHandler handles token refresh requests
func (h *Handler) RefreshTokenHandler() rest.Handler {
	l := logger.With(zap.String("subSection", "token"), zap.String("operation", "refresh"))
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		l = l.With(zap.String("host", r.Host))
		l.Info("received refresh token request")

		data := models.ReceivedDataRefresh{}
		l.Debug("reading request body")
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			l.Error("error reading body", zap.Error(err))
			rest.ERROR(w, err)
			return
		}
		l = l.With(zap.Binary("token", data.RefreshToken))

		l.Debug("refreshing token")
		payload, err := h.rdb.RefreshToken(h.ctx, data.RefreshToken)
		if err != nil {
			l.Error("error refreshing token", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		rest.JSON(w, 200, payload)
	}).Post().JSON().Recover()
}
