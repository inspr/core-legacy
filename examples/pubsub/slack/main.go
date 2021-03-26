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

type SlackMessage struct {
	Text string `json:"text"`
}

var webhook = "https://hooks.slack.com/services/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"
var channel = "pubsubch"

func main() {
	c := &http.Client{}
	client := dappclient.NewAppClient()
	for {

		subMsg, err := client.ReadMessage(context.Background(), channel)
		if err != nil {
			log.Printf("%#v", err.(*ierrors.InsprError).Err)
			continue
		}
		fmt.Println(subMsg.Data)

		msg := SlackMessage{
			Text: fmt.Sprintf("%v", subMsg.Data),
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
