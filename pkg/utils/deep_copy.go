package utils

import (
	"encoding/json"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

func DCopy(root *meta.App) (*meta.App, error) {
	rootObj, err := json.Marshal(*root)
	temp := meta.App{}

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rootObj, &temp)
	return &temp, err
}
