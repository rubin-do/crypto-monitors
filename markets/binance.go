package markets

import (
	"strconv"
	"context"
	"github.com/adshao/go-binance/v2"
	"fmt"
)

func MonitorBinancePrice(prices chan float64) {

loop:
	client := binance.NewClient("", "")
	
	res, err := client.NewDepthService().Symbol("USDTRUB").
		Do(context.Background())
	
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)

	price, err := strconv.ParseFloat(res.Asks[0].Price, 64)

	if err != nil {
		fmt.Println(err)
		return
	}

	prices <- price

	goto loop
}