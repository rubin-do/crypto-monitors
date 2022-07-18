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
	
	needToReport := spread > 1.0
	if bestPair.SellOrderInfo.Market == "Huobi" && spread > 7.0 {
		needToReport = false
	}

	return bestPair, needToReport
}

func main() {

	discordOrders := make(chan discord.BestOrderPair)
	go discord.DiscordSender(discordOrders)

	allOrders := make(chan markets.Order)

	go markets.MonitorGarantexPrice(allOrders)
	go markets.MonitorBinancePrice(allOrders)
	go markets.MonitorByBitPrice(allOrders)
	go markets.MonitorHuobiPrice(allOrders)

	prevBuyOrder := markets.Order{}

	orders := make(map[string]markets.Order)

	for {
		currentOrder := <-allOrders

		switch currentOrder.Market {
		case "Garantex":
			orders["garantex"] = currentOrder
		case "Binance":
			orders["binance"] = currentOrder
		case "ByBit":
			orders["bybit"] = currentOrder
		case "Huobi":
			orders["huobi"] = currentOrder
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
