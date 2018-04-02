package mailjet

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/ViBiOh/httputils/pkg/request"
	"github.com/ViBiOh/httputils/pkg/tools"
)

const mailjetSendURL = `https://api.mailjet.com/v3/send`

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

// App stores informations
type App struct {
	apiPublicKey  string
	apiPrivateKey string
}

// NewApp creates new App from Flags' config
func NewApp(config map[string]*string) *App {
	return &App{
		apiPublicKey:  *config[`apiPublicKey`],
		apiPrivateKey: *config[`apiPrivateKey`],
	}
}

// Flags adds flags for given prefix
func Flags(prefix string) map[string]*string {
	return map[string]*string{
		`apiPublicKey`:  flag.String(tools.ToCamel(fmt.Sprintf(`%sMailjetPublicKey`, prefix)), ``, `Mailjet Public Key`),
		`apiPrivateKey`: flag.String(tools.ToCamel(fmt.Sprintf(`%sMailjetPrivateKey`, prefix)), ``, `Mailjet Private Key`),
	}
}

// Ping indicate if Mailjet is ready or not
func (a *App) Ping() bool {
	return a.apiPublicKey != ``
}

// SendMail send mailjet mail
func (a *App) SendMail(fromEmail string, fromName string, subject string, to []string, html string) error {
	recipients := make([]mailjetRecipient, 0, len(to))
	for _, rawTo := range to {
		recipients = append(recipients, mailjetRecipient{Email: rawTo})
	}

	mailjetMail := mailjetMail{FromEmail: fromEmail, FromName: fromName, Subject: subject, Recipients: recipients, HTML: html}
	if _, err := request.DoJSON(mailjetSendURL, mailjetMail, map[string]string{`Authorization`: request.GetBasicAuth(a.apiPublicKey, a.apiPrivateKey)}, http.MethodPost); err != nil {
		return fmt.Errorf(`Error while sending data to %s: %v`, mailjetSendURL, err)
	}

	return nil
}
