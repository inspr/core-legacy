package utils

import (
	"encoding/json"
)

//DeepCopy is a copy method for dapps tree structures
func DeepCopy(orig, copy interface{}) error {
	rootObj, _ := json.Marshal(orig)

	err := json.Unmarshal(rootObj, copy)
	return err
}
