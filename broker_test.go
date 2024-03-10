package topicgo_test

import (
	"sync"
	"testing"
	"time"

	"github.com/nodejayes/topicgo"
)

func TestCanSubscribeAndBecomeData(t *testing.T) {
	runningTest := false
	testMessage := ""
	topic := "basic"
	var wg sync.WaitGroup

	wg.Add(1)
	topicgo.Subscribe(topic, func(data string) {
		runningTest = true
		testMessage = data
		wg.Done()
	})

	time.Sleep(1 * time.Second)
	topicgo.Publish(topic, "hello")
	wg.Wait()
	if !runningTest {
		t.Errorf("Subscription not triggerd")
	}
	if testMessage != "hello" {
		t.Errorf("data not transfered into subscription")
	}
}

func TestHandleMultiplePublishCalls(t *testing.T) {
	counter := 0
	topic := "multiple"
	var wg sync.WaitGroup

	wg.Add(3)
	topicgo.Subscribe(topic, func(data string) {
		counter++
		wg.Done()
	})

	publishes := 0
	for {
		if publishes > 2 {
			break
		}
		time.Sleep(1 * time.Second)
		topicgo.Publish(topic, "hello")
		publishes++
	}
	wg.Wait()
	if counter != publishes {
		t.Errorf("some Subscriptions not handled counter was %v publishes was %v", counter, publishes)
	}
}

func TestUnsubscribe(t *testing.T) {
	counter := 0
	topic := "unsubscribe"
	var wg sync.WaitGroup

	wg.Add(1)
	unsubscribe := topicgo.Subscribe(topic, func(data string) {
		counter++
		wg.Done()
	})

	publishes := 0
	for {
		if publishes > 0 {
			unsubscribe()
		}
		if publishes > 2 {
			break
		}
		time.Sleep(1 * time.Second)
		topicgo.Publish(topic, "hello")
		publishes++
	}
	wg.Wait()
	if counter == publishes {
		t.Errorf("to much Subscriptions handled counter was %v publishes was %v", counter, publishes)
	}
}
