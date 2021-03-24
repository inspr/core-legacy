package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	dappclient "gitlab.inspr.dev/inspr/core/pkg/client"
)

type DiscordMessage struct {
	Content   string `json:"content"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	TTS       bool   `json:"tts"`
	File      []byte `json:"file"`
	Embeded   []byte `json:"embeds"`
}

type DiscordWebhook struct {
	Type     int    `json:"type"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Avatar   string `json:"type"`
	ChanneId string `json:"channel_id"`
	GuildId  string `json:"guild_id"`
	AppId    string `json:"application_id"`
	Token    string `json:"token"`
}

var webhook = "https://discord.com/api/webhooks/823903452475162666/o9aKLMVOb9-ZhfDD7RJ84OCW6WgU2PQnsEj0CzPgFzKC1icqgqWqF8LZxHSFXEzH1NED"
var channel = "discodMessages"

func main() {
	c := &http.Client{}
	client := dappclient.NewAppClient()
	for {
		// req, _ := http.NewRequest(http.MethodGet, webhook, nil)
		// resp, err := client.Do(req)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// data := &DiscordWebhook{}
		// decoder := json.NewDecoder(resp.Body)
		// err = decoder.Decode(&data)
		// fmt.Println(data.Token)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		subMsg, err := client.ReadMessage(ctx, channel)
		if err != nil {
			continue
		}
		msg := DiscordMessage{
			Content:   fmt.Sprintf("%v", subMsg.Data),
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
			fmt.Println(err)
		}
	}

}
