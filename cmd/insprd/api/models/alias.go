package models

import "inspr.dev/inspr/pkg/meta"

// AliasDI - Data Input(DI) format for requests that pass the alias data
type AliasDI struct {
	Alias  meta.Alias `json:"alias"`
	Ctx    string     `json:"ctx"`
	Target string     `json:"target"`
	DryRun bool       `json:"dry"`
}

// AliasQueryDI - Data Input format for queries requests
type AliasQueryDI struct {
	Ctx    string `json:"ctx"`
	Key    string `json:"key"`
	DryRun bool   `json:"dry"`
}
