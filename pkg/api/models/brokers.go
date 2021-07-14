package models

import (
	"inspr.dev/inspr/pkg/utils"
)

// BrokersDI data interface to provide broker information
type BrokersDI struct {
	Available utils.StringArray `json:"available"`
	Default   string            `json:"default"`
}

// BrokerConfigDI is the struct that defines the means in which the data used
// in operations related to creating or altering the broker in the insprd/cluster
type BrokerConfigDI struct {
	BrokerName   string `json:"brokername"`
	FileContents []byte `json:"filecontents"`
}
