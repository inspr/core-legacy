package models

import "github.com/inspr/inspr/pkg/meta/brokers"

// BrokersDi data interface to provide broker information
type BrokersDi struct {
	Installed brokers.BrokerStatusArray `json:"installed"`
	Default   string                    `json:"default"`
}
