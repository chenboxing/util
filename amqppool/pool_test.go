package amqppool

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	pool := NewConnPool("guest", "guest", "127.0.0.1", "5672", 3, 5, 3)
	ch, err := pool.Get()
	fmt.Println(ch, err)
	pool.Release(ch)
}
