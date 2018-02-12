package market

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"errors"
	"math"
	"strings"
)

const size int = 15
const step float64 = 0.01
const threshold float64 = 0.000001
const klineurl string = "http://api.zb.com/monitor.kline.Nodes/v1/kline"

type KLine struct {
	MoneyType string `json: "moneyType"`
	Symbol string	`json: "symbol"`
	Nodes [][]float64 `json: "monitor.kline.Nodes"`
}

type Monitor struct {
	timeStamp int64
	TradeType string
	amplitude float64
	kline *KLine
	Alert *Alert
}

type Alert struct {
	CoinName string
	Amplitude int
	Duration int
	ReferencePrice float64
	TargetPrice float64
}

func (monitor *Monitor) Pull(duration, size int) error {
	url := fmt.Sprintf("%s?market=%s&type=%dmin&size=%d", klineurl, monitor.TradeType, duration, size)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var kline KLine
	json.Unmarshal(body, &kline)
	if len(market.Data) == size && len(market.Data[0]) == 6 {
		monitor.kline = &kline
	}
	return errors.New("Invalid K line")
}

func (monitor *Monitor) Trigger() (bool, error) {
	err := monitor.Pull(1, size)
	if err != nil {
		return false, err
	}
	
	LatestTimestamp := int64(monitor.kline.Nodes[size-1][0])
	LatestPrice := monitor.kline.Nodes[size-1][4]

	if monitor.timeStamp != LatestTimestamp {
		monitor.timeStamp = LatestTimestamp
		var IndexOfMaxPrice, IndexOfMinPrice, IndexOfPeakPrice int

		for index := size - 2; index >= 0; index-- {
			TimeStamp := int64(monitor.kline.Nodes[index][0])
			TmpPeriod := (LatestTimestamp - TimeStamp) / 60000

			if TmpPeriod > 15 { break }
			if math.Abs( LatestPrice / monitor.kline.Nodes[index][4] - 1.0 ) > threshold {
				if IndexOfMaxPrice == 0 ||  monitor.kline.Nodes[index][4] > monitor.kline.Nodes[IndexOfMaxPrice][4] {
					IndexOfMaxPrice = index
				} else if IndexOfMinPrice ==0 ||  monitor.kline.Nodes[index][4] < monitor.kline.Nodes[IndexOfMinPrice][4] {
					IndexOfMinPrice = index
				}
			}
		}

		if IndexOfMaxPrice > IndexOfMinPrice {
			IndexOfPeakPrice = IndexOfMaxPrice
		} else {
			IndexOfPeakPrice = IndexOfMinPrice
		}

		if IndexOfPeakPrice != 0 {
			var ReferencePrice = monitor.kline.Nodes[IndexOfPeakPrice][4]
			var Amplitude = LatestPrice / ReferencePrice
			var Period = (LatestTimestamp - int64(monitor.kline.Nodes[IndexOfPeakPrice][0])) / 60000

			if math.Abs(Amplitude - monitor.amplitude)> step {
				monitor.amplitude = Amplitude
				HumanAmplitude := int((Amplitude * 100) - 100)
				HumanCoinName := strings.ToUpper(monitor.TradeType[:len(monitor.TradeType)-5])

				monitor.Alert = &Alert{ HumanCoinName,
					HumanAmplitude,
					int(Period),
					ReferencePrice,
					LatestPrice,
				}
				fmt.Println(monitor.Alert)
				return true, nil
			}
		}

	}
	return false, nil
}