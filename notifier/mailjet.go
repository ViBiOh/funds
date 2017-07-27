package notifier

import (
	"fmt"
	"log"
	"os"

	"github.com/ViBiOh/funds/fetch"
)

const mailjetSendURL = `https://api.mailjet.com/v3/send`

var apiPublicKey string
var apiPrivateKey string

type mailjetRecipient struct {
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
func MailjetSend(fromEmail string, fromName string, subject string, to []string, html string) error {
	recipients := make([]mailjetRecipient, 0, len(to))
	for _, rawTo := range to {
		recipients = append(recipients, mailjetRecipient{Email: rawTo})
	}

	mailjetMail := mailjetMail{FromEmail: fromEmail, FromName: fromName, Subject: subject, Recipients: recipients, HTML: html}
	if _, err := fetch.PostJSONBody(mailjetSendURL, mailjetMail, apiPublicKey, apiPrivateKey); err != nil {
		return fmt.Errorf(`Error while sending Mailjet mail: %v`, err)
	}

	return nil
}
