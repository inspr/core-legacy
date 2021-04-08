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
		UID string
		Usr client.User
	}

	data := ReceivedDataCreate{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rest.ERROR(rw, err)
		return
	}

	if err := h.rdb.CreateUser(h.ctx, data.UID, data.Usr); err != nil {
		rest.ERROR(rw, err)
		return
	}
}

func (h *Handler) DeleteUserHandler(rw http.ResponseWriter, r *http.Request) {
	type ReceivedDataDelete struct {
		UID            string
		UsrToBeDeleted string
	}

	data := ReceivedDataDelete{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rest.ERROR(rw, err)
		return
	}

	if err := h.rdb.DeleteUser(h.ctx, data.UID, data.UsrToBeDeleted); err != nil {
		rest.ERROR(rw, err)
		return
	}
}

func (h *Handler) UpdatePasswordHandler(rw http.ResponseWriter, r *http.Request) {
	type ReceivedDataUpdate struct {
		UID            string
		UsrToBeUpdated string
		NewPwd         string
	}

	data := ReceivedDataUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rest.ERROR(rw, err)
		return
	}

	if err := h.rdb.UpdatePassword(h.ctx, data.UID, data.UsrToBeUpdated, data.NewPwd); err != nil {
		rest.ERROR(rw, err)
		return
	}
}

func (h *Handler) LoginHandler(rw http.ResponseWriter, r *http.Request) {
	type ReceivedDataLogin struct {
		UID string
		Pwd string
	}

	data := ReceivedDataLogin{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rest.ERROR(rw, err)
		return
	}

	token, err := h.rdb.Login(h.ctx, data.UID, data.Pwd)
	if err != nil {
		rest.ERROR(rw, err)
		return
	}

	rest.JSON(rw, 200, token)
}

func (h *Handler) RefreshTokenHandler(rw http.ResponseWriter, r *http.Request) {
	type ReceivedDataRefresh struct {
		RefreshToken string
	}

	data := ReceivedDataRefresh{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rest.ERROR(rw, err)
		return
	}

	payload, err := h.rdb.RefreshToken(h.ctx, data.RefreshToken)
	if err != nil {
		rest.ERROR(rw, err)
		return
	}

	rest.JSON(rw, 200, payload)
}
