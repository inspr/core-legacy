package models

import "github.com/inspr/inspr/pkg/meta/brokers"

// BrokersDI data interface to provide broker information
type BrokersDI struct {
	Installed brokers.BrokerStatusArray `json:"installed"`
	Default   string                    `json:"default"`
}
