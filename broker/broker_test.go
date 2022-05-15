package broker

import (
	"testing"
	"time"
)

func TestBroker(t *testing.T) {
	Subscribe("test.topic", func(topic string, message interface{}) {
		t.Logf("*A* => %v: [%s] %v\n", time.Now(), topic, message)
	})

	Publish("test.topic", "1")
	Publish("test.topic", "2")
	Publish("test.topic", 3)
	Publish("test.topic", struct {
		Val int
	}{
		Val: 4,
	})

	Subscribe("test.topic", func(topic string, message interface{}) {
		t.Logf("*B* => %v: [%s] %v\n", time.Now(), topic, message)
	})

	time.Sleep(1 * time.Second)
}
