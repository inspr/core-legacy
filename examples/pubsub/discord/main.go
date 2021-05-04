package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	dappclient "github.com/inspr/inspr/pkg/client"
)

type discordMessage struct {
	Content   string `json:"content"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	TTS       bool   `json:"tts"`
	File      []byte `json:"file"`
	Embedded  []byte `json:"embeds"`
}

var webhook = "https://discord.com/api/webhooks/823903452475162666/o9aKLMVOb9-ZhfDD7RJ84OCW6WgU2PQnsEj0CzPgFzKC1icqgqWqF8LZxHSFXEzH1NED"
var channel = "pubsubch"

type expectedDataType struct {
	Message string `json:"message"`
	Channel string `json:"channel"`
}

func main() {
	c := &http.Client{}
	client := dappclient.NewAppClient()
	client.HandleChannel(channel, func(ctx context.Context, body io.Reader) error {
		decoder := json.NewDecoder(body)

		subMsg := expectedDataType{}
		err := decoder.Decode(&subMsg)
		if err != nil {
			return err
		}

		msg := discordMessage{
			Content:   fmt.Sprintf("%v", subMsg.Message),
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
			return err
		}
		return nil

	})
	log.Fatalln(client.Run(context.Background()))
}
