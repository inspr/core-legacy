package models

import "github.com/inspr/inspr/pkg/meta"

// ChannelDI - Data Input format for requests that pass the channel data
type ChannelDI struct {
	Channel meta.Channel `json:"channel"`
	DryRun  bool         `json:"dry"`
}

// ChannelQueryDI - Data Input format for queries requests
type ChannelQueryDI struct {
	ChName string `json:"chname"`
	DryRun bool   `json:"dry"`
}
