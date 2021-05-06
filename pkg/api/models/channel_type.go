package models

import "github.com/inspr/inspr/pkg/meta"

// TypeDI - Data Input format for requests that pass the Type data
type TypeDI struct {
	Type   meta.Type `json:"Type"`
	Scope  string    `json:"scope"`
	DryRun bool      `json:"dry"`
}

// TypeQueryDI - Data Input format for queries requests
type TypeQueryDI struct {
	Scope  string `json:"scope"`
	CtName string `json:"ctname"`
	DryRun bool   `json:"dry"`
}
