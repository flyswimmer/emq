package channel

import (
	"errors"
	"sync"
)

type Broker struct {
	mutex sync.RWMutex
	chans map[string][]chan Msg
}

func NewBroker() *Broker {
	return &Broker{
		chans: make(map[string][]chan Msg),
	}
}

type Msg struct {
	Topic   string
	Content string
}

func (b *Broker) Send(m Msg) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for _, ch := range b.chans[m.Topic] {
		select {
		case ch <- m:
		default:
			return errors.New("消息队列已满")
		}
	}
	return nil
}

func (b *Broker) Subscribe(topic string, capability int) (<-chan Msg, error) {
	res := make(chan Msg, capability)
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.chans[topic] = append(b.chans[topic], res)
	return res, nil
}
