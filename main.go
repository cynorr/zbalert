package main

import (
	"zbalert/market"
	"time"
	"zbalert/notify"
)

func main() {

	var Coins = market.InitCoins()

	qmail := notify.QQMailSMTP{"369262524@qq.com","fsdlvjiwnsocbifb"}
	to := []string{"cynorr@163.com"}

	for {
		for _, coin := range Coins {
			time.Sleep(1 * time.Second)
			if coin.Name == "" {
				break
			}
			alert := coin.Pull()
			if alert != nil {
				qmail.PushAlert(alert, to)
			}
		}
	}



}
