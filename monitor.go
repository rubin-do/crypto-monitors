package main

import (
	"fmt"
	"monitor/monitor/markets"
	"monitor/monitor/discord"
	"time"
	"math"
)


func main() {
	ten, err := time.ParseDuration("10m")
	
	if err != nil {
		fmt.Println(err)
		return
	}

	lastReport := time.Now().Add(-1 * 600 * time.Second)

	messages := make(chan string)
	go discord.DiscordSender(messages)

	garantex_prices := make(chan float64)
	go markets.MonitorGarantexPrice(garantex_prices)

	binance_prices := make(chan float64)
	go markets.MonitorBinancePrice(binance_prices)

	for {
		garantex_price := <-garantex_prices
		binance_price := <- binance_prices

		if math.Abs(garantex_price - binance_price) > 1.0 && time.Since(lastReport) >= ten {
			lastReport = time.Now()

			message := fmt.Sprintf("binance price: %g\ngarantex price: %g\nspread: %g",
				binance_price,
				garantex_price,
				math.Abs(garantex_price - binance_price),
			)

			fmt.Println(message)

		}
		
		time.Sleep(10 * time.Second)
	}
}
