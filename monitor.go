package main

import (
	"log"
	"monitor/monitor/discord"
	"monitor/monitor/markets"
	"strconv"
)

func FindBestPair(orders map[string]markets.Order) (discord.BestOrderPair, bool) {
	bestPair := discord.BestOrderPair{}
	spread := 0.

	for i, orderFirst := range orders {
		for j, orderSecond := range orders {
			if i == j {
				continue
			}

			curSpread := orderFirst.SellPrice - orderSecond.BuyPrice
			quantityBuy, _ := strconv.ParseFloat(orderSecond.Quantity, 64)

			if curSpread > spread && orderSecond.BuyPrice*quantityBuy > 1000. {
				bestPair.BuyOrderInfo = orderSecond
				bestPair.SellOrderInfo = orderFirst
				spread = curSpread
			}
		}
	}

	return bestPair, spread > 1.0
}

func main() {

	discordOrders := make(chan discord.BestOrderPair)
	go discord.DiscordSender(discordOrders)

	garantexOrders := make(chan markets.Order)
	go markets.MonitorGarantexPrice(garantexOrders)

	binanceOrders := make(chan markets.Order)
	go markets.MonitorBinancePrice(binanceOrders)

	bybitOrders := make(chan markets.Order)
	go markets.MonitorByBitPrice(bybitOrders)

	huobiOrders := make(chan markets.Order)
	go markets.MonitorHuobiPrice(huobiOrders)

	prevBuyOrder := markets.Order{}

	orders := make(map[string]markets.Order)

	for {
		select {
		case garantexOrder := <-garantexOrders:
			orders["garantex"] = garantexOrder
		case binanceOrder := <-binanceOrders:
			orders["binance"] = binanceOrder
		case bybitOrder := <-bybitOrders:
			orders["bybit"] = bybitOrder
		case huobiOrder := <-huobiOrders:
			orders["huobi"] = huobiOrder
		}

		bestPair, report := FindBestPair(orders)

		if report {
			for name, order := range orders {
				log.Printf("%s: %g,", name, order.BuyPrice)
			}
			log.Println()
		}

		if report && prevBuyOrder != bestPair.BuyOrderInfo {
			prevBuyOrder = bestPair.BuyOrderInfo
			discordOrders <- bestPair
		}

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
