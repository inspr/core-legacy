package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	dappclient "inspr.dev/inspr/pkg/client"
	"inspr.dev/inspr/pkg/sidecars/models"
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

func main() {
	c := &http.Client{}
	client := dappclient.NewAppClient()
	client.HandleChannel(
		channel,
		func(ctx context.Context, body io.Reader) error {
			decoder := json.NewDecoder(body)

			subMsg := models.BrokerMessage{}
			err := decoder.Decode(&subMsg)
			if err != nil {
				return err
			}

			msg := discordMessage{
				Content:   fmt.Sprintf("%v", subMsg.Data),
				Username:  "Notifications",
				AvatarURL: "",
				TTS:       true,
			}

			msgBuff, _ := json.Marshal(msg)

			req, _ := http.NewRequest(
				http.MethodPost,
				webhook,
				bytes.NewBuffer(msgBuff),
			)
			head := http.Header{}
			head.Add("Content-type", "application/json")
			req.Header = head
			_, err = c.Do(req)
			if err != nil {
				return err
			}

			return nil

		},
	)
	log.Fatalln(client.Run(context.Background()))
}
