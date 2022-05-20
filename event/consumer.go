package events

import (
	"fmt"
)

type consumer struct {
	_hash uint32
	_id   EventID
	_chan chan *Event
}

func newConsumer(hash uint32, id EventID, echan chan *Event) *consumer {
	return &consumer{
		_hash: hash,
		_id:   id,
		_chan: echan,
	}
}

func (c *consumer) String() string {
	return fmt.Sprintf("hash: %v, event: %v (channel: %v)", c.hash(), c.id(), c.stream())
}

func (c *consumer) hash() uint32 {
	return c._hash
}

func (c *consumer) id() EventID {
	return c._id
}

func (c *consumer) stream() chan *Event {
	return c._chan
}
