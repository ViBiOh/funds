package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const mailjetSendURL = `https://api.mailjet.com/v3/send`
const jsonContentType = `application/json`

var httpClient = http.Client{Timeout: 30 * time.Second}
var apiPublicKey string
var apiPrivateKey string

type mailjetRecipient struct {
	ID    int    `json:"MessageID"`
	Email string `json:"Email"`
}

type mailjetMail struct {
	FromEmail  string             `json:"FromEmail"`
	FromName   string             `json:"FromName"`
	Subject    string             `json:"Subject"`
	Recipients []mailjetRecipient `json:"Recipients"`
	HTML       string             `json:"Html-part"`
}

type mailjetResponse struct {
	Sent []mailjetRecipient `json:"Sent"`
}

// InitMailjet inits API auth
func InitMailjet() {
	apiPublicKey = os.Getenv("MAILJET_APIKEY_PUBLIC")
	apiPrivateKey = os.Getenv("MAILJET_APIKEY_PRIVATE")

	if apiPublicKey != `` {
		log.Print(`Mailjet configured`)
	}
}

// MailjetSend send mailjet mail
func MailjetSend(fromEmail string, fromName string, subject string, to []string, html string) (int, error) {
	recipients := make([]mailjetRecipient, 0, len(to))
	for _, rawTo := range to {
		recipients = append(recipients, mailjetRecipient{Email: rawTo})
	}

	mailRequest := mailjetMail{FromEmail: fromEmail, FromName: fromName, Subject: subject, Recipients: recipients, HTML: html}
	mailRequestJSON, err := json.Marshal(mailRequest)
	if err != nil {
		return 0, fmt.Errorf(`Marshall: %v`, err)
	}

	request, err := http.NewRequest(`POST`, mailjetSendURL, bytes.NewBuffer(mailRequestJSON))
	if err != nil {
		return 0, fmt.Errorf(`Request: %v`, err)
	}

	request.SetBasicAuth(apiPublicKey, apiPrivateKey)
	request.Header.Add(`Content-Type`, jsonContentType)

	resp, err := httpClient.Do(request)
	if err != nil {
		return 0, fmt.Errorf(`Send: %v`, err)
	}

	defer resp.Body.Close()

	responseContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf(`Read: %v`, err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return 0, fmt.Errorf(`Got status %d while sending mail %s`, resp.StatusCode, string(responseContent))
	}

	mailResponse := mailjetResponse{}
	if err := json.Unmarshal(responseContent, &mailResponse); err != nil {
		return 0, fmt.Errorf(`Unmarshal of %s: %v`, string(responseContent), err)
	}

	return mailResponse.Sent[0].ID, nil
}
