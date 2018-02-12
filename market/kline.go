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
const threshold float64 = 0.05
const klineurl string = "http://api.zb.com/monitor.kline.Nodes/v1/kline"

type KLine struct {
	MoneyType string `json: "moneyType"`
	Symbol string	`json: "symbol"`
	Nodes [][]float64 `json: "monitor.kline.Nodes"`
}

type Monitor struct {
	timeStamp int64
	tradeType string
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
	url := fmt.Sprintf("%s?market=%s&type=%dmin&size=%d", klineurl, monitor.tradeType, duration, size)
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
	return errors.New("Invalid K line monitor.kline.Nodes")
}

func (monitor *Monitor) Trigger() (bool, error) {
	err := monitor.Pull(1, size)
	if err != nil {
		return false, err
	}
	
	LatestTimestamp := int64(monitor.kline.Nodes[size-1][0])
	LatestPrice := monitor.kline.Nodes[size-1][4]
	var Amplitude float64
	var Period int64
	var ReferencePrice float64
	var alert *Alert2 = nil

	if monitor.timeStamp != LatestTimestamp {
		monitor.timeStamp = LatestTimestamp
		IndexOfMaxPrice := size - 2
		IndexOfMinPrice := size - 2
		for index := size - 3; index >= 0; index-- {
			TimeStamp := int64(monitor.kline.Nodes[index][0])
			TmpPeriod := (LatestTimestamp - TimeStamp) / 60000
			if TmpPeriod > 15 {
				break
			}
			if monitor.kline.Nodes[index][4] > monitor.kline.Nodes[IndexOfMaxPrice][4] {
				IndexOfMaxPrice = index
			} else if monitor.kline.Nodes[index][4] < monitor.kline.Nodes[IndexOfMinPrice][4] {
				IndexOfMinPrice = index
			}
		}

		UpAmplitude := LatestPrice / monitor.kline.Nodes[IndexOfMinPrice][4]
		DownAmplitude := LatestPrice / monitor.kline.Nodes[IndexOfMaxPrice][4]
		if UpAmplitude > 1.00 + threshold && DownAmplitude < 1.00 - threshold {
			if IndexOfMinPrice > IndexOfMaxPrice {
				Amplitude = UpAmplitude
				ReferencePrice = monitor.kline.Nodes[IndexOfMinPrice][4]
				Period = (LatestTimestamp - int64(monitor.kline.Nodes[IndexOfMinPrice][0])) / 60000
			} else {
				Amplitude = DownAmplitude
				ReferencePrice = monitor.kline.Nodes[IndexOfMaxPrice][4]
				Period = (LatestTimestamp - int64(monitor.kline.Nodes[IndexOfMaxPrice][0])) / 60000
			}
		} else if UpAmplitude > 1.00 + threshold {
			Amplitude = UpAmplitude
			ReferencePrice = monitor.kline.Nodes[IndexOfMinPrice][4]
			Period = (LatestTimestamp - int64(monitor.kline.Nodes[IndexOfMinPrice][0])) / 60000
		} else if DownAmplitude < 1.00 - threshold {
			Amplitude = DownAmplitude
			ReferencePrice = monitor.kline.Nodes[IndexOfMaxPrice][4]
			Period = (LatestTimestamp - int64(monitor.kline.Nodes[IndexOfMaxPrice][0])) / 60000
		}

		if Amplitude != 0 && ( Amplitude * monitor.amplitude <= 0 || math.Abs(Amplitude) - math.Abs(monitor.amplitude)> step ){
			monitor.amplitude = Amplitude
			HumanAmplitude := int((Amplitude * 100) - 100)
			HumanCoinName := strings.ToUpper(monitor.tradeType[0:len(monitor.tradeType)-5])
			monitor.Alert = &Alert{ HumanCoinName,
				HumanAmplitude,
				int(Period),
				ReferencePrice,
				LatestPrice,
			}
			fmt.Println(alert)
		}
	}
	return false, nil
}