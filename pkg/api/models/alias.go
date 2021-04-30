package models

import "github.com/inspr/inspr/pkg/meta"

// AliasDI - Data Input(DI) format for requests that pass the alias data
type AliasDI struct {
	Alias  meta.Alias `json:"alias"`
	Scope  string     `json:"scope"`
	Target string     `json:"target"`
	DryRun bool       `json:"dry"`
}

// AliasQueryDI - Data Input format for queries requests
type AliasQueryDI struct {
	Scope  string `json:"scope"`
	Key    string `json:"key"`
	DryRun bool   `json:"dry"`
}
