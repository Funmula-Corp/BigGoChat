package amqp

import (
	"testing"
)

var service *AMQPClient
var messages []*[]byte = []*[]byte{}

func TestMain(m *testing.M) {
	service = MakeAMQPClient("amqp://guest:guest@localhost:5672")
	m.Run()
}

func getMessages(size int) []*[]byte {
	for i := len(messages); i < size; i++ {
		elem := []byte("The quick brown fox jumps over the lazy dog")
		messages = append(messages, &elem)
	}

	return messages[:size]
}

// go test -v -bench . -benchtime=1000000x
func BenchmarkXxx(b *testing.B) {
	msgs := getMessages(b.N)
	b.ResetTimer()
	for _, msg := range msgs {
		service.Publish(AMQPMessage{
			Exchange: "test",
			Key:      "test",
			Body:     *msg,
		})
	}
}
