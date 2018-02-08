package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"time"
	"bytes"
	"strconv"
	"math"
)

var MarketUrl string = "http://api.zb.com/data/v1/markets"
var KLineUrl string = "http://api.zb.com/data/v1/kline"
const size int = 15
const step float64 = 0.01
const threshold float64 = 0.00005

type Coin struct {
	name string
	timestamp int64
	amplitude float64
}

type Market struct {
	MoneyType string `json: "moneyType"`
	Symbol string	`json: "symbol"`
	Data [][]float64 `json: "data"`
}

func (coin *Coin) Init(name string) {
	coin.name = name
}

func (coin *Coin) Pull() {
	//now := time.Now().Unix()
	//since := now*1000 - 60000*int64(size)
	var buffer bytes.Buffer
	buffer.WriteString(KLineUrl)
	buffer.WriteString("?market=")
	buffer.WriteString(coin.name)
	buffer.WriteString("&type=1min&size=")
	buffer.WriteString(strconv.Itoa(size))
	Url := buffer.String()
	resp, err := http.Get(Url)
	if err != nil {
		fmt.Printf("Coin Err")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var market Market
	json.Unmarshal(body, &market)
	coin.trigger(market.Data)
}

func (coin *Coin) trigger(data [][]float64) {
	LatestTimestamp := int64(data[size-1][0])
	LatestPrice := data[size-1][4]
	var Amplitude float64
	var Period int64
	var ReferencePrice float64

	if coin.timestamp != LatestTimestamp {
		coin.timestamp = LatestTimestamp
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

		fmt.Println(UpAmplitude, DownAmplitude, Amplitude)

		if Amplitude != 0 && ( Amplitude * coin.amplitude <= 0 || math.Abs(Amplitude) - math.Abs(coin.amplitude)> step ){
			coin.amplitude = Amplitude
			fmt.Println(strings.ToUpper(coin.name[0:len(coin.name)-5]), Period, Amplitude, ReferencePrice, LatestPrice)
		}
	}
}

func InitCoins() {
	resp, err := http.Get(MarketUrl)
	if err != nil {
		fmt.Printf("Error")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var m map[string] interface{}
	json.Unmarshal(body, &m)
	for k, _ := range m {
		if strings.HasSuffix(k, "_usdt") {
			fmt.Println(k)
		}
	}
}


func main() {
	//InitCoins()
	var coins_name []string
	InitCoins(coins_name)

	coin := Coin{}
	coin.Init("lbtc_usdt")
	fmt.Println(time.Now().Unix())
	coin.Pull()
}
