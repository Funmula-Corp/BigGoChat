package app

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/amqp"
	"git.biggo.com/Funmula/mattermost-funmula/server/v8/channels/store/storetest/mocks"
	amqp091 "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const testAMQPServer = "amqp://guest:guest@localhost:5672/"

type testAMQPConsumer struct {
	_notifications       []*model.PushNotification
	_notificationsHeader []*map[string]interface{}
	_notificationAcks    []*model.PushNotificationAck
	_numReqs             int
	conn                 *amqp091.Connection
	lock                 *sync.RWMutex
	subscribe            chan struct {
		exited <-chan struct{}
		notify chan<- int
	}
}

func (ac *testAMQPConsumer) Shutdown() {
	ac.conn.Close()
}

func (ac *testAMQPConsumer) notifications() []*model.PushNotification {
	ac.lock.RLock()
	defer ac.lock.RUnlock()
	return ac._notifications
}

func (ac *testAMQPConsumer) notificationsHeader() []*map[string]interface{} {
	ac.lock.RLock()
	defer ac.lock.RUnlock()
	return ac._notificationsHeader
}

func (ac *testAMQPConsumer) notificationAcks() []*model.PushNotificationAck {
	ac.lock.RLock()
	defer ac.lock.RUnlock()
	return ac._notificationAcks
}

func (ac *testAMQPConsumer) numReqs() int {
	ac.lock.RLock()
	defer ac.lock.RUnlock()
	return ac._numReqs
}

// wait numReqs
func (ac *testAMQPConsumer) waitNumReqs(c int) error {
	reqs := ac.numReqs()
	if reqs >= c {
		return nil
	}

	exit := make(chan struct{})
	reqChan := make(chan int)

	sub := struct {
		exited <-chan struct{}
		notify chan<- int
	}{
		exited: exit,
		notify: reqChan,
	}

	ac.subscribe <- sub

	defer close(exit)
	for {
		select {
		case reqs = <-reqChan:
			if reqs >= c {
				return nil
			}
		case <-time.After(2 * time.Second):
			return fmt.Errorf("consumer timeout/less req, expect: %d, get: %d", c, reqs)
		}
	}
}

func newTestAMQPConsumer() *testAMQPConsumer {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, _, err := amqp.Connect(ctx, testAMQPServer)
	if err != nil {
		panic(err)
	}

	channel, err := client.Channel()
	if err != nil {
		panic(err)
	}
	err = channel.ExchangeDeclare(pushProxyAMQPExchange, amqp091.ExchangeTopic, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	queue, err := channel.QueueDeclare(model.NewId(), false, true, false, false, nil)
	if err != nil {
		panic(err)
	}
	err = channel.QueueBind(queue.Name, sendToPushProxyAMQPKey, pushProxyAMQPExchange, false, nil)
	if err != nil {
		panic(err)
	}
	err = channel.QueueBind(queue.Name, ackToPushProxyAMQPKey, pushProxyAMQPExchange, false, nil)
	if err != nil {
		panic(err)
	}
	consume, err := channel.Consume(queue.Name, "test consumer", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	ac := &testAMQPConsumer{
		conn: client,
		lock: new(sync.RWMutex),
		subscribe: make(chan struct {
			exited <-chan struct{}
			notify chan<- int
		}),
	}

	reqChan := make(chan int)

	go func() {
		subscribers := []struct {
			exited <-chan struct{}
			notify chan<- int
		}{}

		broadcase := func(num int) {
			for _, subscriber := range subscribers {
				select {
				case <-subscriber.exited:
					// closed, do nothing
				default:
					subscriber.notify <- num
				}
			}
		}

		for {
			select {
			case subscriber := <-ac.subscribe:
				subscribers = append(subscribers, subscriber)
			case num, ok := <-reqChan:
				if !ok {
					return
				}
				go broadcase(num)
			}
		}
	}()

	go func() {
		defer close(reqChan)
		for message := range consume {
			ac.lock.Lock()
			ac._numReqs += 1
			if message.RoutingKey == sendToPushProxyAMQPKey {
				var notification model.PushNotification
				if json.Unmarshal(message.Body, &notification) == nil {
					ac._notifications = append(ac._notifications, &notification)
					header := map[string]interface{}(message.Headers)
					ac._notificationsHeader = append(ac._notificationsHeader, &header)
				}
			} else if message.RoutingKey == ackToPushProxyAMQPKey {
				var notificationAck model.PushNotificationAck
				if json.Unmarshal(message.Body, &notificationAck) == nil {
					ac._notificationAcks = append(ac._notificationAcks, &notificationAck)
				}
			}
			message.Ack(false)
			reqChan <- ac._numReqs
			ac.lock.Unlock()
		}
	}()

	return ac
}

func TestClearPushNotificationSyncAMQP(t *testing.T) {
	th := SetupWithStoreMock(t)
	defer th.TearDown()

	consumer := newTestAMQPConsumer()
	defer consumer.Shutdown()

	sess1 := &model.Session{
		Id:        "id1",
		UserId:    "user1",
		DeviceId:  "test1",
		ExpiresAt: model.GetMillis() + 100000,
	}
	sess2 := &model.Session{
		Id:        "id2",
		UserId:    "user1",
		DeviceId:  "test2",
		ExpiresAt: model.GetMillis() + 100000,
	}

	mockStore := th.App.Srv().Store().(*mocks.Store)
	mockUserStore := mocks.UserStore{}
	mockUserStore.On("Count", mock.Anything).Return(int64(10), nil)
	mockUserStore.On("GetUnreadCount", mock.AnythingOfType("string"), mock.AnythingOfType("bool")).Return(int64(1), nil)
	mockPostStore := mocks.PostStore{}
	mockPostStore.On("GetMaxPostSize").Return(65535, nil)
	mockSystemStore := mocks.SystemStore{}
	mockSystemStore.On("GetByName", "UpgradedFromTE").Return(&model.System{Name: "UpgradedFromTE", Value: "false"}, nil)
	mockSystemStore.On("GetByName", "InstallationDate").Return(&model.System{Name: "InstallationDate", Value: "10"}, nil)
	mockSystemStore.On("GetByName", "FirstServerRunTimestamp").Return(&model.System{Name: "FirstServerRunTimestamp", Value: "10"}, nil)

	mockSessionStore := mocks.SessionStore{}
	mockSessionStore.On("GetSessionsWithActiveDeviceIds", mock.AnythingOfType("string")).Return([]*model.Session{sess1, sess2}, nil)
	mockSessionStore.On("UpdateDeviceId", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return("testdeviceID", nil)
	mockStore.On("User").Return(&mockUserStore)
	mockStore.On("Post").Return(&mockPostStore)
	mockStore.On("System").Return(&mockSystemStore)
	mockStore.On("Session").Return(&mockSessionStore)
	mockStore.On("GetDBSchemaVersion").Return(1, nil)

	// When CRT is disabled
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.EmailSettings.PushNotificationAMQPServer = testAMQPServer
		*cfg.ServiceSettings.CollapsedThreads = model.CollapsedThreadsDisabled
	})

	err := th.App.clearPushNotificationSync(th.Context, sess1.Id, "user1", "channel1", "")
	require.Nil(t, err)
	assert.Nil(t, consumer.waitNumReqs(1))
	// Server side verification.
	// We verify that 1 request has been sent, and also check the message contents.
	require.Equal(t, 1, consumer.numReqs())
	assert.Equal(t, "channel1", consumer.notifications()[0].ChannelId)
	assert.Equal(t, model.PushTypeClear, consumer.notifications()[0].Type)

	// When CRT is enabled, Send badge count adding both "User unreads" + "User thread mentions"
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.ServiceSettings.ThreadAutoFollow = true
		*cfg.ServiceSettings.CollapsedThreads = model.CollapsedThreadsDefaultOn
	})

	mockPreferenceStore := mocks.PreferenceStore{}
	mockPreferenceStore.On("Get", mock.AnythingOfType("string"), model.PreferenceCategoryDisplaySettings, model.PreferenceNameCollapsedThreadsEnabled).Return(&model.Preference{Value: "on"}, nil)
	mockStore.On("Preference").Return(&mockPreferenceStore)

	mockThreadStore := mocks.ThreadStore{}
	mockThreadStore.On("GetTotalUnreadMentions", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything).Return(int64(3), nil)
	mockStore.On("Thread").Return(&mockThreadStore)

	err = th.App.clearPushNotificationSync(th.Context, sess1.Id, "user1", "channel1", "")
	require.Nil(t, err)
	assert.Nil(t, consumer.waitNumReqs(2))
	assert.Equal(t, consumer.notifications()[1].Badge, 4)
}

func TestUpdateMobileAppBadgeSyncAMQP(t *testing.T) {
	th := SetupWithStoreMock(t)
	defer th.TearDown()

	consumer := newTestAMQPConsumer()
	defer consumer.Shutdown()

	sess1 := &model.Session{
		Id:        "id1",
		UserId:    "user1",
		DeviceId:  "test1",
		ExpiresAt: model.GetMillis() + 100000,
	}
	sess2 := &model.Session{
		Id:        "id2",
		UserId:    "user1",
		DeviceId:  "test2",
		ExpiresAt: model.GetMillis() + 100000,
	}

	mockStore := th.App.Srv().Store().(*mocks.Store)
	mockUserStore := mocks.UserStore{}
	mockUserStore.On("Count", mock.Anything).Return(int64(10), nil)
	mockUserStore.On("GetUnreadCount", mock.AnythingOfType("string"), mock.AnythingOfType("bool")).Return(int64(1), nil)
	mockPostStore := mocks.PostStore{}
	mockPostStore.On("GetMaxPostSize").Return(65535, nil)
	mockSystemStore := mocks.SystemStore{}
	mockSystemStore.On("GetByName", "UpgradedFromTE").Return(&model.System{Name: "UpgradedFromTE", Value: "false"}, nil)
	mockSystemStore.On("GetByName", "InstallationDate").Return(&model.System{Name: "InstallationDate", Value: "10"}, nil)
	mockSystemStore.On("GetByName", "FirstServerRunTimestamp").Return(&model.System{Name: "FirstServerRunTimestamp", Value: "10"}, nil)

	mockSessionStore := mocks.SessionStore{}
	mockSessionStore.On("GetSessionsWithActiveDeviceIds", mock.AnythingOfType("string")).Return([]*model.Session{sess1, sess2}, nil)
	mockSessionStore.On("UpdateDeviceId", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return("testdeviceID", nil)
	mockStore.On("User").Return(&mockUserStore)
	mockStore.On("Post").Return(&mockPostStore)
	mockStore.On("System").Return(&mockSystemStore)
	mockStore.On("Session").Return(&mockSessionStore)
	mockStore.On("GetDBSchemaVersion").Return(1, nil)

	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.EmailSettings.PushNotificationAMQPServer = testAMQPServer
		*cfg.ServiceSettings.CollapsedThreads = model.CollapsedThreadsDisabled
	})

	err := th.App.updateMobileAppBadgeSync(th.Context, "user1")
	require.Nil(t, err)
	assert.Nil(t, consumer.waitNumReqs(2))
	// Server side verification.
	// We verify that 2 requests have been sent, and also check the message contents.
	require.Equal(t, 2, consumer.numReqs())
	assert.Equal(t, 1, consumer.notifications()[0].ContentAvailable)
	assert.Equal(t, model.PushTypeUpdateBadge, consumer.notifications()[0].Type)
	assert.Equal(t, 1, consumer.notifications()[1].ContentAvailable)
	assert.Equal(t, model.PushTypeUpdateBadge, consumer.notifications()[1].Type)
}

func TestSendAckToPushProxyAMQP(t *testing.T) {
	th := SetupWithStoreMock(t)
	defer th.TearDown()

	consumer := newTestAMQPConsumer()
	defer consumer.Shutdown()

	mockStore := th.App.Srv().Store().(*mocks.Store)
	mockUserStore := mocks.UserStore{}
	mockUserStore.On("Count", mock.Anything).Return(int64(10), nil)
	mockPostStore := mocks.PostStore{}
	mockPostStore.On("GetMaxPostSize").Return(65535, nil)
	mockSystemStore := mocks.SystemStore{}
	mockSystemStore.On("GetByName", "UpgradedFromTE").Return(&model.System{Name: "UpgradedFromTE", Value: "false"}, nil)
	mockSystemStore.On("GetByName", "InstallationDate").Return(&model.System{Name: "InstallationDate", Value: "10"}, nil)
	mockSystemStore.On("GetByName", "FirstServerRunTimestamp").Return(&model.System{Name: "FirstServerRunTimestamp", Value: "10"}, nil)

	mockStore.On("User").Return(&mockUserStore)
	mockStore.On("Post").Return(&mockPostStore)
	mockStore.On("System").Return(&mockSystemStore)
	mockStore.On("GetDBSchemaVersion").Return(1, nil)

	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.EmailSettings.PushNotificationAMQPServer = testAMQPServer
	})

	ack := &model.PushNotificationAck{
		Id:               "testid",
		NotificationType: model.PushTypeMessage,
	}
	err := th.App.SendAckToPushProxy(ack)
	require.NoError(t, err)
	assert.Nil(t, consumer.waitNumReqs(1))
	// Server side verification.
	// We verify that 1 request has been sent, and also check the message contents.
	require.Equal(t, 1, consumer.numReqs())
	assert.Equal(t, ack.Id, consumer.notificationAcks()[0].Id)
	assert.Equal(t, ack.NotificationType, consumer.notificationAcks()[0].NotificationType)
}

// TestAllPushNotifications is a master test which sends all various types
// of notifications and verifies they have been properly sent.
func TestAllPushNotificationsAMQP(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping all push notifications test in short mode")
	}

	th := Setup(t).InitBasic()
	defer th.TearDown()

	// Create 10 users, each having 2 sessions.
	type userSession struct {
		user    *model.User
		session *model.Session
	}
	var testData []userSession
	for i := 0; i < 10; i++ {
		u := th.CreateUser()
		sess, err := th.App.CreateSession(th.Context, &model.Session{
			UserId:    u.Id,
			DeviceId:  "deviceID" + u.Id,
			ExpiresAt: model.GetMillis() + 100000,
		})
		require.Nil(t, err)
		// We don't need to track the 2nd session.
		_, err = th.App.CreateSession(th.Context, &model.Session{
			UserId:    u.Id,
			DeviceId:  "deviceID" + u.Id,
			ExpiresAt: model.GetMillis() + 100000,
		})
		require.Nil(t, err)
		_, err = th.App.AddTeamMember(th.Context, th.BasicTeam.Id, u.Id)
		require.Nil(t, err)
		th.AddUserToChannel(u, th.BasicChannel)
		testData = append(testData, userSession{
			user:    u,
			session: sess,
		})
	}

	consumer := newTestAMQPConsumer()
	defer consumer.Shutdown()

	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.EmailSettings.PushNotificationContents = model.GenericNotification
		*cfg.EmailSettings.PushNotificationAMQPServer = testAMQPServer
	})

	var wg sync.WaitGroup
	for i, data := range testData {
		wg.Add(1)
		// Ranging between 3 types of notifications.
		switch i % 3 {
		case 0:
			go func(user model.User) {
				defer wg.Done()
				notification := &PostNotification{
					Post:    th.CreatePost(th.BasicChannel),
					Channel: th.BasicChannel,
					ProfileMap: map[string]*model.User{
						user.Id: &user,
					},
					Sender: &user,
				}
				// testing all 3 notification types.
				th.App.sendPushNotification(notification, &user, true, false, model.CommentsNotifyAny)
			}(*data.user)
		case 1:
			go func(id string) {
				defer wg.Done()
				th.App.UpdateMobileAppBadge(id)
			}(data.user.Id)
		case 2:
			go func(sessID, userID string) {
				defer wg.Done()
				th.App.clearPushNotification(sessID, userID, th.BasicChannel.Id, "")
			}(data.session.Id, data.user.Id)
		}
	}
	wg.Wait()

	// Hack to let the worker goroutines complete.
	assert.Nil(t, consumer.waitNumReqs(17))
	// Server side verification.
	assert.Equal(t, 17, consumer.numReqs())
	var numClears, numMessages, numUpdateBadges int
	for _, n := range consumer.notifications() {
		switch n.Type {
		case model.PushTypeClear:
			numClears++
			assert.Equal(t, th.BasicChannel.Id, n.ChannelId)
		case model.PushTypeMessage:
			numMessages++
			assert.Equal(t, th.BasicChannel.Id, n.ChannelId)
			assert.Contains(t, n.Message, "mentioned you")
		case model.PushTypeUpdateBadge:
			numUpdateBadges++
			assert.Equal(t, "none", n.Sound)
			assert.Equal(t, 1, n.ContentAvailable)
		}
	}
	assert.Equal(t, 8, numMessages)
	assert.Equal(t, 3, numClears)
	assert.Equal(t, 6, numUpdateBadges)
}

func TestSendPushToPushProxyPriorityAMQP(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	consumer := newTestAMQPConsumer()
	defer consumer.Shutdown()

	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.EmailSettings.PushNotificationAMQPServer = testAMQPServer
	})

	msg := &model.PushNotification{
		Version:  "2",
		Type:     model.PushTypeTest,
		Priority: model.PostPriorityUrgent,
	}
	msg.SetDeviceIdAndPlatform("platform:id")

	_, err := th.App.rawSendToPushProxy(msg)
	require.Nil(t, err)
	assert.Nil(t, consumer.waitNumReqs(1))
	headers := consumer.notificationsHeader()
	require.Len(t, headers, 1)
	require.Equal(t, model.PostPriorityUrgent, (*headers[0])["priority"])
}
