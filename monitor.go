package main

import (
	"log"
	"monitor/monitor/discord"
	"monitor/monitor/markets"
	"time"
)

func FindBestPair(orders []markets.Order) (discord.BestOrderPair, bool) {
	best_pair := discord.BestOrderPair{}
	spread := 0.

	for i, order_first := range orders {
		for j, order_second := range orders {
			if i == j {
				continue
			}

			cur_spread := order_first.Price - order_second.Price

			if cur_spread > spread {
				best_pair.BuyOrderInfo = order_second
				best_pair.SellOrderInfo = order_first
				spread = cur_spread
			}
		}
	}

	return best_pair, (spread > 1.0)
}

func main() {

	discordOrders := make(chan discord.BestOrderPair)
	go discord.DiscordSender(discordOrders)

	garantex_orders := make(chan markets.Order)
	go markets.MonitorGarantexPrice(garantex_orders)

	binance_orders := make(chan markets.Order)
	go markets.MonitorBinancePrice(binance_orders)

	bybit_orders := make(chan markets.Order)
	go markets.MonitorByBitPrice(bybit_orders)

	prevBuyOrder := markets.Order{}

	for {
		garantex_order := <-garantex_orders
		binance_order := <-binance_orders
		bybit_order := <-bybit_orders

		log.Printf("Binance price: %g, garantex price: %g, bybit price: %g", binance_order.Price, garantex_order.Price, bybit_order.Price)

		best_pair, report := FindBestPair([]markets.Order{garantex_order, binance_order, bybit_order})

		log.Println(best_pair)

		if report && prevBuyOrder != best_pair.BuyOrderInfo {
			prevBuyOrder = best_pair.BuyOrderInfo
			discordOrders <- best_pair
		}

		time.Sleep(10 * time.Second)
	}
}

// discordOrders <- discord.BestOrderPair{
// 	BuyOrderInfo: markets.Order{
// 		Market:     "TestSource",
// 		SellerName: "TestSeller",
// 		Price:      50.0,
// 		Quantity:   "123.4",
// 		MinAmount:  "228.0",
// 		MaxAmount:  "1337.0",
// 		Url:        "https://www.google.com",
// 	},
// 	SellOrderInfo: markets.Order{
// 		Market: "TestOutput",
// 		Price:  121.0,
// 	},
// }
