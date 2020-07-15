package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bitfinexcom/bfx-utility-scripts/pkg/csvgenerator"
	"github.com/bitfinexcom/bfx-utility-scripts/pkg/request"
	"github.com/shopspring/decimal"
)

var (
	tPairsURL         = "https://api-pub.bitfinex.com/v2/tickers?symbols=ALL"
	tradingVolumeURL  = "https://api-pub.bitfinex.com/v2/candles/trade:%s:%s/hist?limit=%s"
	candleVolumeIndex = 5
)

func main() {
	pairs, err := getTradingPairs()
	if err != nil {
		log.Fatalf("getTradingPairs: %s", err)
	}

	csvData := [][]string{{"Pair", "1 Day", "1 Week", "1 Month", "90 Days"}}
	processedPairs := 0
	allPairsLen := len(pairs)
	log.Printf("Collected %d trading pairs.\n", allPairsLen)

	for _, pair := range pairs {
		log.Printf("Processing:%s | Done:%d | Total:%d\n", pair, processedPairs, allPairsLen)

		csvRow, err := getCSVRow(pair)
		if err != nil {
			log.Fatalf("failed at %s. Err:%s", pair, err)
		}

		csvData = append(csvData, csvRow)
		processedPairs++
	}

	fileName := fmt.Sprintf("tp_volume_report_%d.csv", time.Now().UTC().UnixNano())
	if err := csvgenerator.Generate(fileName, csvData); err != nil {
		log.Fatal(err)
	}
}

func getCSVRow(pair string) ([]string, error) {
	oneDayVol, err := getTradingVolume("1D", pair, "1")
	if err != nil {
		return nil, err
	}

	oneWeekVol, err := getTradingVolume("1D", pair, "7")
	if err != nil {
		return nil, err
	}

	oneMonthVol, err := getTradingVolume("1D", pair, "31")
	if err != nil {
		return nil, err
	}

	ninetyDaysVol, err := getTradingVolume("1D", pair, "90")
	if err != nil {
		return nil, err
	}

	csvRow := []string{pair}
	csvRow = append(csvRow, oneDayVol.String())
	csvRow = append(csvRow, oneWeekVol.String())
	csvRow = append(csvRow, oneMonthVol.String())
	csvRow = append(csvRow, ninetyDaysVol.String())
	return csvRow, nil
}

func getTradingPairs() ([]string, error) {
	pairs := [][]interface{}{}

	req := request.New("GET", tPairsURL, nil).Do().Decode(&pairs)
	if req.Err != nil {
		return nil, req.Err
	}

	tPairs := []string{}

	for _, v := range pairs {
		pair := v[0]
		pairStr, ok := pair.(string)
		if !ok {
			log.Fatalf("Could not convert pair:%v to string", pair)
		}

		if strings.HasPrefix(pairStr, "t") {
			tPairs = append(tPairs, pairStr)
		}
	}

	return tPairs, nil
}

func getTradingVolume(timeFrame, pair, limit string) (decimal.Decimal, error) {
	sumVolume := decimal.NewFromFloat(0)
	candles := [][]interface{}{}

	req := request.
		New("GET", fmt.Sprintf(tradingVolumeURL, timeFrame, pair, limit), nil).
		Do().
		Decode(&candles)

	if req.Err != nil {
		return sumVolume, req.Err
	}

	for _, candle := range candles {
		vol := candle[candleVolumeIndex]
		volFlt64, ok := vol.(float64)
		if !ok {
			return sumVolume, fmt.Errorf("Could not convert volume:%v to float64", vol)
		}

		sumVolume = sumVolume.Add(decimal.NewFromFloat(volFlt64))
	}

	time.Sleep(1 * time.Second)
	return sumVolume, nil
}
