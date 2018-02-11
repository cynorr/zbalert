package notify

import (
	"fmt"
	"net/smtp"
	"strings"
	"zbalert/market"
)

type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

type SMTP struct {
	User string
	Password string
	Host string
}

func (mail *SMTP) PushAlert(message *market.Alert, to []string) {
	basehost := mail.Host
	auth := unencryptedAuth{smtp.PlainAuth("", mail.User, mail.Password, basehost)}

	host := fmt.Sprintf("%s:25", mail.Host)
	nickname := fmt.Sprintf("%s  %+d%%", message.CoinName, message.Amplitude)
	subject := fmt.Sprintf("%g  ->  %g", message.ReferencePrice, message.TargetPrice)
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := "Alert"
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + mail.User + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail(host, auth, mail.User, to, msg)
	if err != nil {
		fmt.Printf("send mail error: %v", err)
	}
}
