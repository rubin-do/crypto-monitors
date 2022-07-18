package discord

import (
	"monitor/monitor/markets"
	"testing"
	"time"
)

func TestDiscord(t *testing.T) {
	discordOrders := make(chan BestOrderPair)
	go DiscordSender(discordOrders)

	discordOrders <- BestOrderPair{
		BuyOrderInfo: markets.Order{
			Market:     "Binance",
			SellerName: "TestSeller",
			BuyPrice:   50.0,
			SellPrice:  51.0,
			Quantity:   "123.4",
			MinAmount:  "228.0",
			MaxAmount:  "1337.0",
			Url:        "https://www.google.com",
		},
		SellOrderInfo: markets.Order{
			Market:    "TestOutput",
			BuyPrice:  120.0,
			SellPrice: 121.0,
		},
	}
	time.Sleep(5)
}
