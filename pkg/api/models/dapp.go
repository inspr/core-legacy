package models

import (
	"github.com/inspr/inspr/pkg/meta"
)

// AppDI - Data Input(DI) format for requests that pass the app data
type AppDI struct {
	App    meta.App `json:"app"`
	Scope  string   `json:"scope"`
	DryRun bool     `json:"dry"`
}

// AppQueryDI - Data Input format for queries requests
type AppQueryDI struct {
	Scope  string `json:"scope"`
	DryRun bool   `json:"dry"`
}
