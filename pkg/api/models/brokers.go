package models

import (
	"github.com/inspr/inspr/pkg/utils"
)

// BrokersDI data interface to provide broker information
type BrokersDI struct {
	Installed utils.StringArray `json:"installed"`
	Default   string            `json:"default"`
}
