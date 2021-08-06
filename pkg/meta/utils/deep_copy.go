package utils

import (
	"encoding/json"
	"reflect"

	"inspr.dev/inspr/pkg/ierrors"
)

//DeepCopy is a copy method for dapps tree structures
func DeepCopy(orig, dest interface{}) error {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return ierrors.New("dest must be a pointer").InvalidName()
	}

	rootObj, _ := json.Marshal(orig)

	err := json.Unmarshal(rootObj, dest)
	return err
}
