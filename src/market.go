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
)

var MarketUrl string = "http://api.zb.com/data/v1/markets"
var KLineUrl string = "http://api.zb.com/data/v1/kline"
const size int = 60

type Coin struct {
	name string
	timestamp int
	line [size]int
}

func (coin *Coin) Init(name string) {
	coin.name = name
}

func (coin *Coin) Pull() {
	now := time.Now().Unix()
	since := now*1000 - 60000*int64(size)
	var buffer bytes.Buffer
	buffer.WriteString(KLineUrl)
	buffer.WriteString("?market=")
	buffer.WriteString(coin.name)
	buffer.WriteString("&type=1min&since=")
	buffer.WriteString(strconv.FormatInt(since, 10))
	Url := buffer.String()
	fmt.Println(Url)
	resp, err := http.Get(Url)
	if err != nil {
		fmt.Printf("Coin Err")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//fmt.Println(string(body))
	var m map[string] interface{}
	json.Unmarshal(body, &m)
	klist := m["data"]
	fmt.Println(klist)
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
	coin := Coin{}
	coin.Init("btc_usdt")
	fmt.Println(time.Now().Unix())
	coin.Pull()
}