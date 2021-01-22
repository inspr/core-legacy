package models

import "gitlab.inspr.dev/inspr/core/pkg/meta"

// ChannelDI - Data Input format for requests that pass the channel data
type ChannelDI struct {
	Channel meta.Channel `json:"channel"`
	Ctx     string       `json:"ctx"`
}

// ChannelQueryDI - Data Input format for queries requests
type ChannelQueryDI struct {
	Query  string `json:"query"`
	ChName string `json:"chname"`
}
