package app

import (
	"encoding/json"
	"fmt"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/amqp"
)

func (a *App) rawSendToPushProxyAMQP(msg *model.PushNotification) (model.PushResponse, error) {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode to JSON: %w", err)
	}

	a.Srv().pushNotificationAMQPClient.Publish(amqp.AMQPMessage{
		Exchange: "chat.pushNotification",
		Key: "push",
		Body: msgJSON,
	})
	return model.NewOkPushResponse(), nil
}