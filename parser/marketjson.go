package parser

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"marketracker/models"
	"sync"
)

func DecodeTradeJSON(jsonChan <-chan string, tradesChn chan<- models.Trade, group *sync.WaitGroup) {
	var trade models.Trade
	for {
		jsonString, ok := <-jsonChan
		if ok == false {
			group.Done()
			break
		} else {
			err := jsoniter.Unmarshal([]byte(jsonString), &trade)
			if err != nil {
				fmt.Println(err)
				continue
			}
			tradesChn <- trade
		}
	}
}
