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

type slackMessage struct {
	Text string `json:"text"`
}

type expectedDataType struct {
	Message string `json:"message"`
	Channel string `json:"channel"`
}

var webhook = "https://hooks.slack.com/services/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"
var channel = "pubsubch"

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

		msg := slackMessage{
			Text: fmt.Sprintf("%v", subMsg.Message),
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
