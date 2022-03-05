package models

type MarketStat struct {
	NumTrades        int
	TotalPrice       float64
	TotalVolume      float64
	VolumeTimesPrice float64
	NumBuys          int
}
