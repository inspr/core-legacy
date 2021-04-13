package models

import "github.com/inspr/inspr/pkg/meta"

// ChannelTypeDI - Data Input format for requests that pass the channelType data
type ChannelTypeDI struct {
	ChannelType meta.ChannelType `json:"channeltype"`
	Ctx         string           `json:"ctx"`
	DryRun      bool             `json:"dry"`
}

// ChannelTypeQueryDI - Data Input format for queries requests
type ChannelTypeQueryDI struct {
	Ctx    string `json:"ctx"`
	CtName string `json:"ctname"`
	DryRun bool   `json:"dry"`
}
