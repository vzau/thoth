package discord

import (
	"fmt"

	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
	"github.com/dhawton/log4g"
)

var log = log4g.Category("discord")

func SendToDiscord(webhookToken string, message string) {
	if webhookToken == "" {
		return
	}

	webhook, err := disgohook.NewWebhookClientByToken(nil, nil, webhookToken)
	if err != nil {
		log.Error("Error creating webhook client: %s", err.Error())
		return
	}

	_, err = webhook.SendMessage(api.NewWebhookMessageCreateBuilder().SetContent(message).Build())
	if err != nil {
		log.Error("Error sending message: %s", err.Error())
	}
}

func SendToDiscordf(webhookToken string, message string, values ...interface{}) {
	SendToDiscord(webhookToken, fmt.Sprintf(message, values...))
}
