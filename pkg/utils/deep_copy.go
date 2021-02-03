package utils

import (
	"encoding/json"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

func DCopy(root *meta.App) *meta.App {
	rootObj, _ := json.Marshal(*root)
	temp := meta.App{}

	_ = json.Unmarshal(rootObj, &temp)
	return &temp
}
