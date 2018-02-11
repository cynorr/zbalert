package main

import (
	"zbalert/market"
	"time"
	"zbalert/notify"
)

func main() {


	var Coins = market.InitCoins()

	const mailSize = 4
	var mails [mailSize]notify.SMTP
	mails[0] = notify.SMTP{"cynorr@2980.com", "AAAaaa123", "smtp.2980.com"}
	mails[1] = notify.SMTP{"cynorr1@tom.com", "AAAaaa123", "smtp.tom.com"}
	mails[2] = notify.SMTP{"cynorr2@tom.com", "AAAaaa123", "smtp.tom.com"}
	mails[3] = notify.SMTP{"cynorr@2980.com", "AAAaaa1234", "smtp.2980.com"}

	to := []string{"cynorr@163.com"}

	i := 0

	for {
		for index, _ := range Coins {
			time.Sleep(1 * time.Second)
			if Coins[index].Name == "" {
				break
			}
			alert := Coins[index].Pull()
			if alert != nil {
				mails[i%mailSize].PushAlert(alert, to)
				i++
			}
		}
	}
}
