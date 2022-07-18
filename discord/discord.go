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
	webhooks := make(map[string]webhook.Client)
	webhooks["Binance"] = webhook.New(snowflake.ID(998671929877741659), "DZf3Md9HkReJmCr8OvuyYj-wo9otJPxCHe4HvY9ikr5qWJQxGJnR9J2KPqRIBD2ULPo4")
	webhooks["Garantex"] = webhook.New(snowflake.ID(998672189274464336), "XOkqMObkLJYLGsE_Iea8ycYRvp-Ib__glMACzr5GM-FCVWCpkTq3Ov3LYGIVgC3XcVUA")
	webhooks["Huobi"] = webhook.New(snowflake.ID(998672704418893834), "6kQTTc-8nveZQzoWxBRdpuQdTUGuEL5xF8HwefHjqBIRYjY1sbl4naplVbgFcX3quAH8")
	webhooks["ByBit"] = webhook.New(snowflake.ID(991651258299588658), "w8-zzyWDeEsdOJDAl0nVHVvPaXa6Q6Z8dlsT0pNyXOtAJCeB1u0_oQqRqUn_2Djm6wTn")

	for {
		pair := <-data

		_, err := webhooks[pair.BuyOrderInfo.Market].CreateEmbeds([]discord.Embed{discord.NewEmbedBuilder().
			SetTitle(pair.BuyOrderInfo.Market).
			SetAuthor(pair.BuyOrderInfo.SellerName, "", "").
			SetURL(pair.BuyOrderInfo.Url).
			SetFooterText(fmt.Sprintf("Spread: %g", pair.SellOrderInfo.SellPrice-pair.BuyOrderInfo.BuyPrice)).
			AddField("Price", fmt.Sprintf("%g", pair.BuyOrderInfo.BuyPrice), false).
			AddField("Quantity", pair.BuyOrderInfo.Quantity, false).
			AddField("MinAmount", pair.BuyOrderInfo.MinAmount, false).
			AddField("MaxAmount", pair.BuyOrderInfo.MaxAmount, false).
			AddField("PaymentMethods", pair.BuyOrderInfo.PaymentMethods, false).
			AddField(pair.SellOrderInfo.Market, fmt.Sprintf("%g", pair.SellOrderInfo.SellPrice), false).
			SetTimestamp(time.Now()).
			Build()})

		if err != nil {
			log.Fatal(err)
		}
	}
}
