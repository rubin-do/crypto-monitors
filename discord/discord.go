package discord

import (
//	"fmt"
	"github.com/DisgoOrg/disgohook"
//	"strconv"
)

func DiscordSender(data chan string) {
	webhook, err := disgohook.NewWebhookClientByToken(nil, nil, "990357629031309392/UnIMrvtS8_NIdtJpwlNl11Mxn38N1ddlfUKPyQ8iZ4ctDz1NBLtPf74a24TvMbhj61Qv")

	if err != nil {
		panic(err)
	}

	for {
		webhook.SendContent(<-data)
	}
}
