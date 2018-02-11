package market

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"bytes"
	"strconv"
	"math"
)

var MarketUrl string = "http://api.zb.com/data/v1/markets"
var KLineUrl string = "http://api.zb.com/data/v1/kline"
var market Market

const size int = 15
const step float64 = 0.01
const threshold float64 = 0.05

type Alert struct {
	CoinName string
	Amplitude int
	Duration int
	ReferencePrice float64
	TargetPrice float64
}

type Coin struct {
	Name string
	Timestamp int64
	Amplitude float64
}

type Market struct {
	MoneyType string `json: "moneyType"`
	Symbol string	`json: "symbol"`
	Data [][]float64 `json: "data"`
}

func (coin *Coin) Pull() *Alert {
	var buffer bytes.Buffer
	var alert *Alert = nil
	buffer.WriteString(KLineUrl)
	buffer.WriteString("?market=")
	buffer.WriteString(coin.Name)
	buffer.WriteString("&type=1min&size=")
	buffer.WriteString(strconv.Itoa(size))
	Url := buffer.String()
	resp, err := http.Get(Url)
	if err != nil {
		fmt.Printf("Coin Err")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Parser error")
	}
	defer resp.Body.Close()

	json.Unmarshal(body, &market)
	if len(market.Data) == size && len(market.Data[0]) == 6 {
		alert = coin.trigger(market.Data)
	}
	return alert
}

func (coin *Coin) trigger(data [][]float64) *Alert{
	LatestTimestamp := int64(data[size-1][0])
	LatestPrice := data[size-1][4]
	var Amplitude float64
	var Period int64
	var ReferencePrice float64
	var alert *Alert = nil

	if coin.Timestamp != LatestTimestamp {
		coin.Timestamp = LatestTimestamp
		IndexOfMaxPrice := size - 2
		IndexOfMinPrice := size - 2
		for index := size - 3; index >= 0; index -- {
			TimeStamp := int64(data[index][0])
			TmpPeriod := (LatestTimestamp - TimeStamp) / 60000
			if TmpPeriod > 15 {
				break
			}
			if data[index][4] > data[IndexOfMaxPrice][4] {
				IndexOfMaxPrice = index
			} else if data[index][4] < data[IndexOfMinPrice][4] {
				IndexOfMinPrice = index
			}
		}

		UpAmplitude := LatestPrice / data[IndexOfMinPrice][4]
		DownAmplitude := LatestPrice / data[IndexOfMaxPrice][4]
		if UpAmplitude > 1.00 + threshold && DownAmplitude < 1.00 - threshold {
			if IndexOfMinPrice > IndexOfMaxPrice {
				Amplitude = UpAmplitude
				ReferencePrice = data[IndexOfMinPrice][4]
				Period = (LatestTimestamp - int64(data[IndexOfMinPrice][0])) / 60000
			} else {
				Amplitude = DownAmplitude
				ReferencePrice = data[IndexOfMaxPrice][4]
				Period = (LatestTimestamp - int64(data[IndexOfMaxPrice][0])) / 60000
			}
		} else if UpAmplitude > 1.00 + threshold {
			Amplitude = UpAmplitude
			ReferencePrice = data[IndexOfMinPrice][4]
			Period = (LatestTimestamp - int64(data[IndexOfMinPrice][0])) / 60000
		} else if DownAmplitude < 1.00 - threshold {
			Amplitude = DownAmplitude
			ReferencePrice = data[IndexOfMaxPrice][4]
			Period = (LatestTimestamp - int64(data[IndexOfMaxPrice][0])) / 60000
		}

		if Amplitude != 0 && ( Amplitude * coin.Amplitude <= 0 || math.Abs(Amplitude) - math.Abs(coin.Amplitude)> step ){
			coin.Amplitude = Amplitude
			HumanAmplitude := int((Amplitude * 100) - 100)
			alert = &Alert{strings.ToUpper(coin.Name[0:len(coin.Name)-5]),
				HumanAmplitude,
				int(Period),
				ReferencePrice,
				LatestPrice,
				}
			fmt.Println(strings.ToUpper(coin.Name[0:len(coin.Name)-5]), Period, HumanAmplitude, ReferencePrice, LatestPrice)
		}
	}
	return alert
}

func InitCoins() [60]Coin{
	var Coins [60]Coin
	resp, err := http.Get(MarketUrl)
	if err != nil {
		fmt.Printf("Error")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Parser Error")
	}
	defer resp.Body.Close()
	var m map[string] interface{}
	json.Unmarshal(body, &m)

	index := 0
	for k, _ := range m {
		if strings.HasSuffix(k, "_usdt") {
			Coins[index] = Coin{Name: k}
			index ++
		}
	}
	return Coins
}


//func main() {
//
//
//	InitCoins()
//
//	for _, coin := range Coins {
//		time.Sleep(1 * time.Second)
//		if coin.Name == "" {
//			break
//		}
//		coin.Pull()
//	}
//
//	fmt.Println("Starting ... ")
//	time.Sleep(1)
//	coin := Coin{name: "lbtc_usdt"}
//	coin.Pull()
//	fmt.Println("End ... ")
//}
