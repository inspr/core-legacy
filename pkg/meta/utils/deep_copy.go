package utils

import (
	"encoding/json"
	"reflect"

	"github.com/inspr/inspr/pkg/ierrors"
)

//DeepCopy is a copy method for dapps tree structures
func DeepCopy(orig, dest interface{}) error {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return ierrors.NewError().Message("dest must be a pointer").InvalidName().Build()
	}

	rootObj, _ := json.Marshal(orig)

	err := json.Unmarshal(rootObj, dest)
	return err
}
