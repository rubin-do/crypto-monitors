package discord

import (
	"monitor/monitor/markets"
	"testing"
)

func TestDiscord(t *testing.T) {
	discordOrders := make(chan BestOrderPair)
	go DiscordSender(discordOrders)

	discordOrders <- BestOrderPair{
		BuyOrderInfo: markets.Order{
			Market:     "TestSource",
			SellerName: "TestSeller",
			Price:      50.0,
			Quantity:   "123.4",
			MinAmount:  "228.0",
			MaxAmount:  "1337.0",
			Url:        "https://www.google.com",
		},
		SellOrderInfo: markets.Order{
			Market: "TestOutput",
			Price:  121.0,
		},
	}
}
