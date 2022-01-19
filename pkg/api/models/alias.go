package models

import "inspr.dev/inspr/pkg/meta"

// AliasDI - Data Input(DI) format for requests that pass the alias data
type AliasDI struct {
	Alias  meta.Alias `json:"alias"`
	DryRun bool       `json:"dry"`
}

// AliasQueryDI - Data Input format for queries requests
type AliasQueryDI struct {
	Name   string `json:"name"`
	DryRun bool   `json:"dry"`
}
