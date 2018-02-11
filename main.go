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
		for index, _ := range Coins {
			time.Sleep(1 * time.Second)
			if Coins[index].Name == "" {
				break
			}
			alert := Coins[index].Pull()
			if alert != nil {
				go qmail.PushAlert(alert, to)
			}
		}
	}
}
