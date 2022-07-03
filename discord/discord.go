package discord

import (
	"fmt"
	"log"
	"monitor/monitor/markets"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
)

type BestOrderPair struct {
	BuyOrderInfo  markets.Order
	SellOrderInfo markets.Order
}

func DiscordSender(data chan BestOrderPair) {

	client := webhook.New(snowflake.ID(990357629031309392), "UnIMrvtS8_NIdtJpwlNl11Mxn38N1ddlfUKPyQ8iZ4ctDz1NBLtPf74a24TvMbhj61Qv")

	for {
		pair := <-data

		_, err := client.CreateEmbeds([]discord.Embed{discord.NewEmbedBuilder().
			SetTitle(pair.BuyOrderInfo.Market).
			SetAuthor(pair.BuyOrderInfo.SellerName, "", "").
			SetURL(pair.BuyOrderInfo.Url).
			SetFooterText(fmt.Sprintf("Spread: %g", pair.SellOrderInfo.SellPrice-pair.BuyOrderInfo.BuyPrice)).
			AddField("Price", fmt.Sprintf("%g", pair.BuyOrderInfo.BuyPrice), false).
			AddField("Quantity", pair.BuyOrderInfo.Quantity, false).
			AddField("MinAmount", pair.BuyOrderInfo.MinAmount, false).
			AddField("MaxAmount", pair.BuyOrderInfo.MaxAmount, false).
			AddField(pair.SellOrderInfo.Market, fmt.Sprintf("%g", pair.SellOrderInfo.SellPrice), false).
			SetTimestamp(time.Now()).
			Build()})

		if err != nil {
			log.Fatal(err)
		}
	}
}
