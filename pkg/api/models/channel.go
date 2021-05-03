package models

import "github.com/inspr/inspr/pkg/meta"

// ChannelDI - Data Input format for requests that pass the channel data
type ChannelDI struct {
	Channel meta.Channel `json:"channel"`
	Scope   string       `json:"scope"`
	DryRun  bool         `json:"dry"`
}

// ChannelQueryDI - Data Input format for queries requests
type ChannelQueryDI struct {
	Scope  string `json:"scope"`
	ChName string `json:"chname"`
	DryRun bool   `json:"dry"`
}
