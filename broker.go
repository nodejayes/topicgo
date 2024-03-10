package topicgo


import (
	"sync"

	"github.com/google/uuid"
	"github.com/nodejayes/anything-parse-json"
)

var m = &sync.Mutex{}
var subscriptions = make(map[string]map[string]func(payload any))

func Publish[T any](topic string, payload T) {
	if len(subscriptions[topic]) < 1 {
		return
	}
	for _, subscription := range subscriptions[topic] {
		subscription(payload)
	}
}

func Subscribe[T any](topic string, handler func(payload T)) func() {
	m.Lock()
	defer m.Unlock()

	id := uuid.NewString()
	if subscriptions[topic] == nil {
		subscriptions[topic] = make(map[string]func(payload any))
	}
	subscriptions[topic][id] = func(data any) {
		convertedData, err := anythingparsejson.Parse[T](data)
		if err != nil {
			return
		}
		handler(convertedData)
	}
	return func() {
		m.Lock()
		defer m.Unlock()

		delete(subscriptions[topic], id)
		if len(subscriptions[topic]) < 1 {
			delete(subscriptions, topic)
		}
	}
}
