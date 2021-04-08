package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/uid_provider/client"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

type Handler struct {
	rdb client.RedisManager
	ctx context.Context
}

func NewHandler(rdb client.RedisManager) *Handler {
	return &Handler{
		rdb: rdb,
	}
}

func (h *Handler) CreateUserHandler(rw http.ResponseWriter, r *http.Request) {
	type ReceivedDataCreate struct {
		uid string
		usr client.User
	}

	data := ReceivedDataCreate{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rest.ERROR(rw, err)
		return
	}

	if err := h.rdb.CreateUser(h.ctx, data.uid, data.usr); err != nil {
		rest.ERROR(rw, err)
		return
	}
}

func (h *Handler) DeleteUserHandler(rw http.ResponseWriter, r *http.Request) {
	type ReceivedDataDelete struct {
		uid            string
		usrToBeDeleted string
	}

	data := ReceivedDataDelete{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rest.ERROR(rw, err)
		return
	}

	if err := h.rdb.DeleteUser(h.ctx, data.uid, data.usrToBeDeleted); err != nil {
		rest.ERROR(rw, err)
		return
	}
}

func (h *Handler) UpdatePasswordHandler(rw http.ResponseWriter, r *http.Request) {
	type ReceivedDataUpdate struct {
		uid            string
		usrToBeUpdated string
		newPwd         string
	}

	data := ReceivedDataUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rest.ERROR(rw, err)
		return
	}

	if err := h.rdb.UpdatePassword(h.ctx, data.uid, data.usrToBeUpdated, data.newPwd); err != nil {
		rest.ERROR(rw, err)
		return
	}
}

func (h *Handler) LoginHandler(rw http.ResponseWriter, r *http.Request) {
	type ReceivedDataLogin struct {
		uid string
		pwd string
	}

	data := ReceivedDataLogin{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rest.ERROR(rw, err)
		return
	}

	if err := h.rdb.Login(h.ctx, data.uid, data.pwd); err != nil {
		rest.ERROR(rw, err)
		return
	}
}

func (h *Handler) RefreshTokenHandler(rw http.ResponseWriter, r *http.Request) {
	type ReceivedDataRefresh struct {
		refreshToken string
	}

	data := ReceivedDataRefresh{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rest.ERROR(rw, err)
		return
	}

	payload, err := h.rdb.RefreshToken(h.ctx, data.refreshToken)
	if err != nil {
		rest.ERROR(rw, err)
		return
	}

	rest.JSON(rw, 200, payload)
}
