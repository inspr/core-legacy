package auth

import (
	"encoding/json"

	"gitlab.inspr.dev/inspr/core/pkg/auth/models"
)

// Desserialize converts a interface to a Payload model
func Desserialize(jwtLoad interface{}) models.Payload {

	jwtJSON, err := json.Marshal(jwtLoad)
	if err != nil {
		return models.Payload{}
	}

	var payload models.Payload
	err = json.Unmarshal(jwtJSON, &payload)
	if err != nil {
		return models.Payload{}
	}
	return payload
}
