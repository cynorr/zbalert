package notify

import (
	"fmt"
	"net/smtp"
	"strings"
	"zbalert/market"
)

type QQMailSMTP struct {
	User string
	Password string
}

func (mail *QQMailSMTP) PushAlert(message *market.Alert, to []string) {
	basehost := "smtp.qq.com"
	auth := smtp.PlainAuth("", mail.User, mail.Password, basehost)

	host := "smtp.qq.com:25"
	nickname := fmt.Sprintf("%s\t%+d%%", message.CoinName, message.Amplitude)
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
