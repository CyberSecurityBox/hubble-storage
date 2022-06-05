package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hubble-storage/src/env"
	"hubble-storage/src/postgresql"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func SendWithVerdict(flow *postgresql.Flow, verdicts []string, log *logrus.Logger) error {
	var flag = false
	var slackPayload BlockPayload
	var err error

	for _, v := range verdicts {
		if flow.Verdict == v {
			flag = true
			break
		}
	}

	if flag {
		slackPayload.Blocks  = make([]blocks, 2)
		slackPayload.Blocks[0].Type = "section"
		slackPayload.Blocks[0].Text = &text{
			Type: "mrkdwn",
			Text: "You have a new alert :loaloa:",
		}
		slackPayload.Blocks[1].Type = "section"
		slackPayload.Blocks[1].Fields = make([]fields, 6)
		slackPayload.Blocks[1].Fields[0].Type = "mrkdwn"
		slackPayload.Blocks[1].Fields[0].Text = fmt.Sprintf("*Verdict:*\n%s", flow.Verdict)
		slackPayload.Blocks[1].Fields[1].Type = "mrkdwn"
		slackPayload.Blocks[1].Fields[1].Text = fmt.Sprintf("*Layer:*\n%s", flow.Layer)
		slackPayload.Blocks[1].Fields[2].Type = "mrkdwn"
		slackPayload.Blocks[1].Fields[2].Text = fmt.Sprintf("*Source Name:*\n%s/%s", flow.SourceNamespace, flow.SourcePodName)
		slackPayload.Blocks[1].Fields[3].Type = "mrkdwn"
		slackPayload.Blocks[1].Fields[3].Text = fmt.Sprintf("*Source Port:*\n%d", flow.L4SourcePort)
		slackPayload.Blocks[1].Fields[4].Type = "mrkdwn"
		slackPayload.Blocks[1].Fields[4].Text = fmt.Sprintf("*Destination Name:*\n%s/%s", flow.DestinationNamespace, flow.DestinationPodName)
		slackPayload.Blocks[1].Fields[5].Type = "mrkdwn"
		slackPayload.Blocks[1].Fields[5].Text = fmt.Sprintf("*Destination Port:*\n%d", flow.L4DestinationPort)

		// Send to slack
		if err = sendSlackNotification(env.WebhookUrl, slackPayload, log); err != nil {
			return err
		}
	}

	return nil
}

// SendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func sendSlackNotification(webhookUrl string, payload BlockPayload, log *logrus.Logger) error {
	var (
		err error
		req *http.Request
		resp *http.Response
		slackBody []byte
		client *http.Client
		buf *bytes.Buffer
	)

	if slackBody, err = json.Marshal(&payload); err != nil {
		return err
	}

	fmt.Println(string(slackBody))

	if req, err = http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody)); err != nil {
		log.WithError(err).Error("Can't create slack new request")
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	client = &http.Client{
		Timeout: 10 * time.Second,
	}

	if resp, err = client.Do(req); err != nil {
		return err
	}

	buf = new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		log.WithError(err).Error("Can't read buffer")
		return err
	}
	if buf.String() != "ok" {
		err = errors.New("non-ok response returned from slack")
		log.WithError(err).Error(err.Error())
		return err
	}

	return nil
}
