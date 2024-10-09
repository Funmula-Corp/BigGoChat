package app

import (
	"encoding/json"
	"fmt"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/amqp"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
)

const (
	pushProxyAMQPExchange = "chat.pushNotification"
)

func (a *App) rawSendToPushProxyAMQP(msg *model.PushNotification) (model.PushResponse, error) {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode to JSON: %w", err)
	}
	// TODO: add msg priroity to header

	a.Srv().pushNotificationAMQPClient.Publish(amqp.AMQPMessage{
		Exchange: pushProxyAMQPExchange,
		Key: "send_push",
		Body: msgJSON,
	})
	return model.NewOkPushResponse(), nil
}

func (a *App) SendAckToPushProxyAMQP(ack *model.PushNotificationAck) error {
	if ack == nil {
		return nil
	}

	a.NotificationsLog().Trace("Notification successfully received",
		mlog.String("type", model.NotificationTypePush),
		mlog.String("ack_id", ack.Id),
		mlog.String("push_type", ack.NotificationType),
		mlog.String("post_id", ack.PostId),
		mlog.String("ack_type", ack.NotificationType),
		mlog.String("device_type", ack.ClientPlatform),
		mlog.Int("received_at", ack.ClientReceivedAt),
		mlog.String("status", model.PushReceived),
	)

	ackJSON, err := json.Marshal(ack)
	if err != nil {
		return fmt.Errorf("failed to encode to JSON: %w", err)
	}

	a.Srv().pushNotificationAMQPClient.Publish(amqp.AMQPMessage{
		Exchange: pushProxyAMQPExchange,
		Key: "ack",
		Body: ackJSON,
	})

	return nil
}