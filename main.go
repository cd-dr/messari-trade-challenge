package main

import (
	"fmt"
	"marketracker/ingester"
	"marketracker/models"
	"marketracker/parser"
	"sync"
)

func main() {
	jsonChn := make(chan string)
	tradesChn := make(chan models.Trade)
	var parserGroup sync.WaitGroup
	var processThread sync.WaitGroup
	numParserThread := 8

	marketDict := map[int]*models.MarketStat{}

	//process the trades
	processThread.Add(1)
	go func() {
		for {
			trade, ok := <-tradesChn
			if ok == false {
				processThread.Done()
				break
			} else {
				buy := 0
				if trade.IsBuy {
					buy = 1
				}
				if _, ok := marketDict[trade.Market]; ok {
					marketDict[trade.Market].NumTrades++
					marketDict[trade.Market].TotalPrice += trade.Price
					marketDict[trade.Market].TotalVolume += trade.Volume
					marketDict[trade.Market].VolumeTimesPrice += trade.Volume * trade.Price
					marketDict[trade.Market].NumBuys += buy
				} else {
					marketDict[trade.Market] = &models.MarketStat{
						NumTrades:        1,
						NumBuys:          buy,
						TotalVolume:      trade.Volume,
						TotalPrice:       trade.Price,
						VolumeTimesPrice: trade.Volume * trade.Price,
					}
				}
			}
		}
	}()
	go func() {
		for i := 0; i < numParserThread; i++ {
			parserGroup.Add(1)
			go parser.DecodeTradeJSON(jsonChn, tradesChn, &parserGroup)
		}
		parserGroup.Wait()
		close(tradesChn)
	}()
	go ingester.IngestFromStdin(jsonChn)

	processThread.Wait()
	for key, val := range marketDict {
		fmt.Printf("{\"market\":%d,\"total_volume\":%f,\"mean_price\":%f,\"mean_volume\":%f,\"volume_weighted_average_price\":%f,\"percentage_buy\":%f}\n",
			key, val.TotalVolume, val.TotalPrice/float64(val.NumTrades),
			val.TotalVolume/float64(val.NumTrades), val.VolumeTimesPrice/val.TotalVolume, float64(val.NumBuys)/float64(val.NumTrades))
	}
}
