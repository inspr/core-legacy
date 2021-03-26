package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dappclient "gitlab.inspr.dev/inspr/core/pkg/client"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

type discordMessage struct {
	Content   string `json:"content"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	TTS       bool   `json:"tts"`
	File      []byte `json:"file"`
	Embeded   []byte `json:"embeds"`
}

var webhook = "https://discord.com/api/webhooks/823903452475162666/o9aKLMVOb9-ZhfDD7RJ84OCW6WgU2PQnsEj0CzPgFzKC1icqgqWqF8LZxHSFXEzH1NED"
var channel = "pubsubch"

type expectedDataType struct {
	Message struct {
		Data string `json:"data"`
	} `json:"message"`
	Channel string `json:"channel"`
}

func main() {
	c := &http.Client{}
	client := dappclient.NewAppClient()
	for {
		subMsg := expectedDataType{}
		err := client.ReadMessage(context.Background(), channel, &subMsg)
		if err != nil {
			log.Printf("%#v", err.(*ierrors.InsprError).Err)
			continue
		}

		msg := DiscordMessage{
			Content:   fmt.Sprintf("%v", subMsg.Message.Data),
			Username:  "Notifications",
			AvatarURL: "",
			TTS:       true,
		}

		msgBuff, _ := json.Marshal(msg)

		req, _ := http.NewRequest(http.MethodPost, webhook, bytes.NewBuffer(msgBuff))
		head := http.Header{}
		head.Add("Content-type", "application/json")
		req.Header = head
		_, err = c.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}

		if err := client.CommitMessage(context.Background(), channel); err != nil {
			log.Println(err.Error())
		}
	}
}
