package main

import (
	"fmt"
	"log"
	"math"
	"monitor/monitor/discord"
	"monitor/monitor/markets"
	"time"
)

func main() {

	messages := make(chan string)
	go discord.DiscordSender(messages)

	garantex_prices := make(chan float64)
	go markets.MonitorGarantexPrice(garantex_prices)

	binance_prices := make(chan float64)
	go markets.MonitorBinancePrice(binance_prices)

	//bybit_prices := make(chan float64)
	//go markets.MonitorByBitPrice(bybit_prices)

	for {
		garantex_price := <-garantex_prices
		binance_price := <-binance_prices
		//<-bybit_prices

		log.Printf("Binance price: %g, garantex price: %g", binance_price, garantex_price)

		if math.Abs(garantex_price-binance_price) > 1.0 {

			message := fmt.Sprintf("binance price: %g\ngarantex price: %g\nspread: %g",
				binance_price,
				garantex_price,
				math.Abs(garantex_price-binance_price),
			)

			fmt.Println(message)

			messages <- message
		}

		time.Sleep(10 * time.Second)
	}
}
