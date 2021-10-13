package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

//DefaultIncomingWebHook for posting to slack
var DefaultIncomingWebHook = os.Getenv("SLACK_INCOMING_WEBHOOK")

//Attachments are slack attachments
type Attachments struct {
	Text string `json:"text,omitempty"`
}

//SlackMessage encapsulatest the message to send to slack
type SlackMessage struct {
	Text        string        `json:"text,omitempty"`
	Attachments []Attachments `json:"attachments,omitempty"`
}

//ResultToSlack will send result to slack
func ResultToSlack(outURL, errURL, action, randomID, status, webhook string) {

	m := ComposeSlackMessage(outURL, errURL, action, randomID, status)
	m.PostToSlack(webhook)

}

//ComposeSlackMessage  composes the mesage to slack
func ComposeSlackMessage(outputURL, errorURL, action, id, actionStatus string) SlackMessage {
	//SlackMessage encapsulatest the message to send to slack
	//id = fmt.Sprintf("%s", id)
	//action = fmt.Sprintf("%s", action)

	topLevelMessage := fmt.Sprintf(`Status for %s %s : %s`, action, id, actionStatus)

	doc := Attachments{Text: fmt.Sprintf(`<%s|See Output Logs for %s id %s>`, outputURL, action, id)}
	doc1 := Attachments{Text: fmt.Sprintf(`<%s|See Error Logs for %s id %s>`, errorURL, action, id)}

	attachements := []Attachments{doc, doc1}

	message := SlackMessage{Text: topLevelMessage, Attachments: attachements}
	return message
}

//PostToSlack post the message to slack
func (m SlackMessage) PostToSlack(webhook string) {
	slackIt, err := json.Marshal(m)

	log.Printf("message %s", slackIt)
	if err != nil {
		log.Printf("Error occured %v", err)
		return
	}
	if webhook == "" {
		webhook = DefaultIncomingWebHook
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(slackIt))
	if err != nil {
		log.Printf("Error occured while invoking callback URL %s, Error is  %v", webhook, err)
		return
	}
	log.Printf("Posted successfully to %s\n", webhook)
	log.Printf("Status code is %d\n", resp.StatusCode)
	defer resp.Body.Close()

}
