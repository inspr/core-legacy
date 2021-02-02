package models

import "gitlab.inspr.dev/inspr/core/pkg/meta"

// ChannelTypeDI - Data Input format for requests that pass the channelType data
type ChannelTypeDI struct {
	ChannelType meta.ChannelType `json:"channeltype"`
	Ctx         string           `json:"ctx"`
	Valid       bool             `json:"valid"`
	Dry         bool             `json:"dry"`
}

// ChannelTypeQueryDI - Data Input format for queries requests
type ChannelTypeQueryDI struct {
	Ctx    string `json:"ctx"`
	CtName string `json:"ctname"`
	Valid  bool   `json:"valid"`
	Dry    bool   `json:"dry"`
}
