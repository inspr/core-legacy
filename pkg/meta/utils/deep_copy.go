package utils

import (
	"encoding/json"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

//DeepCopy is a copy method for dapps tree structures
func DeepCopy(root *meta.App) *meta.App {
	rootObj, _ := json.Marshal(*root)
	temp := meta.App{}

	_ = json.Unmarshal(rootObj, &temp)
	return &temp
}
