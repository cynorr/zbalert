package main

import (
	"net/smtp"
	"fmt"
	"strings"
	"errors"
)


type plainAuth struct {
	identity, username, password string
	host string
}

func UnEncryptedPlainAuth(identity, username, password, host string) *plainAuth {
	return &plainAuth{identity, username, password, host}
}

func (a *plainAuth) Start(server *ServerInfo) (string, []byte, error) {

	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}
	resp := []byte(a.identity + "x00" + a.username + "x00" + a.password)
	return "PLAIN", resp, nil
}

func (a *plainAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	return nil, nil
}

func main() {
	user := "cynorr@tom.com"
	password := "AAAaaa123"
	basehost := "smtp.tom.com:143"
	host := "smtp.tom.com:25"
	to := []string{"cynorr@163.com"}


	auth := UnEncryptedPlainAuth("", user, password, basehost)
	nickname := "Test" //fmt.Sprintf("%s\t%+d%%", message.CoinName, message.Amplitude)
	subject := "hahah" // fmt.Sprintf("%g  ->  %g", message.ReferencePrice, message.TargetPrice)
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := "Alert"
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)

	err := smtp.SendMail(host, auth, user, to, msg)
	if err != nil {
		fmt.Printf("send mail error: %v", err)
	}
}
