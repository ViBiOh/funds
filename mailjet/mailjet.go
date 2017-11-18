package mailjet

import (
	"flag"
	"fmt"

	"github.com/ViBiOh/httputils"
)

const mailjetSendURL = `https://api.mailjet.com/v3/send`

var (
	apiPublicKey  = flag.String(`mailjetPublicKey`, ``, `Mailet Public Key`)
	apiPrivateKey = flag.String(`mailjetPrivateKey`, ``, `Mailet Private Key`)
)

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

// Ping indicate if Mailjet is ready or not
func Ping() bool {
	return *apiPublicKey != ``
}

// SendMail send mailjet mail
func SendMail(fromEmail string, fromName string, subject string, to []string, html string) error {
	recipients := make([]mailjetRecipient, 0, len(to))
	for _, rawTo := range to {
		recipients = append(recipients, mailjetRecipient{Email: rawTo})
	}

	mailjetMail := mailjetMail{FromEmail: fromEmail, FromName: fromName, Subject: subject, Recipients: recipients, HTML: html}
	if _, err := httputils.PostJSONBody(mailjetSendURL, mailjetMail, map[string]string{`Authorization`: httputils.GetBasicAuth(*apiPublicKey, *apiPrivateKey)}); err != nil {
		return fmt.Errorf(`Error while sending data to %s: %v`, mailjetSendURL, err)
	}

	return nil
}
