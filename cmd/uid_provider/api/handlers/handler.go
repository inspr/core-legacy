package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"inspr.dev/inspr/cmd/uid_provider/api/models"
	"inspr.dev/inspr/cmd/uid_provider/client"
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
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received create user request", zap.String("host", r.Host))

		data := models.ReceivedDataCreate{}
		logger.Debug("reading request body")
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			logger.Error("error reading body", zap.Error(err))
			rest.ERROR(w, err)
			return
		}
		logger.Debug("creating user in redis", zap.Any("created-user", data.User), zap.String("auth-user-uid", data.UID))
		if err := h.rdb.CreateUser(h.ctx, data.UID, data.Password, data.User); err != nil {
			logger.Error("error creating user in redis", zap.Any("created-user", data.User), zap.String("auth-user", data.UID), zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		logger.Info("user created", zap.String("username", data.UID))
	}).Post().JSON().Recover()
}

// DeleteUserHandler handles user deletion requests
func (h *Handler) DeleteUserHandler() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received delete user request", zap.String("host", r.Host))

		data := models.ReceivedDataDelete{}
		logger.Debug("reading request body")
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			logger.Error("error reading body", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("deleting user in redis", zap.Any("deleted-user", data.UserToBeDeleted), zap.String("auth-user-uid", data.UID))
		if err := h.rdb.DeleteUser(h.ctx, data.UID, data.Password, data.UserToBeDeleted); err != nil {
			logger.Error("error deleting user in redis", zap.Any("created-user", data.UserToBeDeleted), zap.String("auth-user", data.UID), zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		logger.Info("user deleted", zap.String("username", data.UID))
	}).Delete().JSON().Recover()
}

// UpdatePasswordHandler handles requests to update an user password
func (h *Handler) UpdatePasswordHandler() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received update request", zap.String("host", r.Host))

		data := models.ReceivedDataUpdate{}
		logger.Debug("reading request body")
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			logger.Error("error reading body", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("updating user in redis", zap.Any("updated-user", data.UserToBeUpdated), zap.String("auth-user-uid", data.UID))
		if err := h.rdb.UpdatePassword(h.ctx, data.UID, data.Password,
			data.UserToBeUpdated, data.NewPassword); err != nil {
			logger.Error("error updating user in redis", zap.Any("created-user", data.UserToBeUpdated), zap.String("auth-user", data.UID), zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		logger.Info("user updated", zap.String("username", data.UID))
	}).Put().JSON().Recover()
}

// LoginHandler handles login requests
func (h *Handler) LoginHandler() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received login request", zap.String("host", r.Host))

		data := models.ReceivedDataLogin{}
		logger.Debug("reading request body")
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			logger.Error("error reading body", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("logging in user username provided", zap.String("UID", data.UID))
		token, err := h.rdb.Login(h.ctx, data.UID, data.Password)
		if err != nil {
			logger.Debug("error logging in", zap.String("UID", data.UID))
			rest.ERROR(w, err)
			return
		}

		rest.JSON(w, 200, token)
	}).Post().JSON().Recover()
}

// RefreshTokenHandler handles token refresh requests
func (h *Handler) RefreshTokenHandler() rest.Handler {
	return rest.Handler(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received refresh token request", zap.String("host", r.Host))

		data := models.ReceivedDataRefresh{}
		logger.Debug("reading request body")
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			logger.Error("error reading body", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("refreshing token")
		payload, err := h.rdb.RefreshToken(h.ctx, data.RefreshToken)
		if err != nil {
			logger.Error("error refreshing token", zap.Binary("token", data.RefreshToken), zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		rest.JSON(w, 200, payload)
	}).Post().JSON().Recover()
}
