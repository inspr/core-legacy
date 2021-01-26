package models

import (
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// AppDI - Data Input(DI) format for requests that pass the app data
type AppDI struct {
	App   meta.App `json:"app"`
	Ctx   string   `json:"ctx"`
	Setup bool     `json:"setup"`
}

// AppQueryDI - Data Input format for queries requests
type AppQueryDI struct {
	Query string `json:"query"`
	Setup bool   `json:"setup"`
}
