package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBroker_Send(t *testing.T) {
	b := NewBroker()
	go func() {
		for {
			err := b.Send(Msg{
				Topic:   "hello",
				Content: time.Now().String(),
			})
			if err != nil {
				t.Log(err)
				return
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("消费者 %d", i)
		go func(name string) {
			defer wg.Done()
			msgs, err := b.Subscribe("world", 100)
			if err != nil {
				t.Log(err)
				return
			}
			for msg := range msgs {
				fmt.Println(name, msg.Content)
			}
		}(name)
	}
	wg.Wait()
}
