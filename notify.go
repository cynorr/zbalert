package main

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

func main() {
	user := "369262524@qq.com"
	password := "aa"
	basehost := "smtp.qq.com"
	host := "smtp.qq.com:25"
	to := []string{"cynorr@163.com"}
	auth := smtp.PlainAuth("", user, password, basehost)
	nickname := "test"
	subject := "test mail"
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := "This is the email body."
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)

	for i := 0; i < 100; i++ {
		err := smtp.SendMail(host, auth, user, to, msg)
		if err != nil {
			fmt.Printf("send mail error: %v", err)
		}
	}

	fmt.Println("End ... ")
}
