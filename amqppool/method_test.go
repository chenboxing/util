package amqppool

import (
	"fmt"
	"sync"
	"time"

	"testing"
)

func TestPublish(t *testing.T) {
	Configure("127.0.0.1", "5672", "guest", "guest")
	/*go func() {
		err := Publish(TradeQueue, []byte("123"))
		fmt.Println(err)
	}()

	go func() {
		err := Publish(TradeQueue, []byte("123"))
		fmt.Println(err)
	}()

	go func() {
		err := Publish(TradeQueue, []byte("123"))
		fmt.Println(err)
	}()

	go func() {
		err := Publish(TradeQueue, []byte("123"))
		fmt.Println(err)
	}()

	go func() {
		err := Publish(TradeQueue, []byte("123"))
		fmt.Println(err)
	}()

	go func() {
		err := Publish(TradeQueue, []byte("123"))
		fmt.Println(err)
	}()*/

	// // time.Sleep(1 * time.Second)

	for i := 0; i < 100; i++ {

		go func() {
			err := Publish(TradeQueue, []byte("123"))
			if err != nil {
				//panic(err)
				fmt.Println("==============", err)
			}

			time.Sleep(1 * time.Millisecond)
		}()
	}

	time.Sleep(20 * time.Second)
	// fmt.Println("===========--------===")

	// go func() {
	// 	err := Publish(TradeQueue, []byte("123"))
	// 	fmt.Println(err)
	// }()

	// go func() {
	// 	err := Publish(TradeQueue, []byte("123"))
	// 	fmt.Println(err)
	// }()

	// err := Publish(TradeQueue, []byte("123"))
	// fmt.Println(err)

	time.Sleep(100 * time.Second)

}

func TestConsume(t *testing.T) {
	Configure("127.0.0.1", "5672", "guest", "guest")

	test1, err := Consume(TradeQueue)
	if err != nil {
		fmt.Println(1, err)
		panic(err)
	}

	for msg := range test1 {
		fmt.Println(msg.RoutingKey, string(msg.Body))
		msg.Ack(false)
	}
}

func TestConsume2(t *testing.T) {
	Configure("127.0.0.1", "5672", "guest", "guest")

	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		test1, err := Consume(TradeQueue)
		if err != nil {
			fmt.Println(1, err)
			panic(err)
		}

		for msg := range test1 {
			fmt.Println("1", msg.RoutingKey, string(msg.Body))
			msg.Ack(false)
		}
		wg.Done()
	}()
	time.Sleep(time.Second)

	go func() {
		wg.Add(1)
		test1, err := Consume(RefundQueue)
		if err != nil {
			fmt.Println(1, err)
			panic(err)
		}

		for msg := range test1 {
			fmt.Println("2", msg.RoutingKey, string(msg.Body))
			msg.Ack(false)
		}
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		err := Publish(TradeQueue, []byte("123"))
		fmt.Println(err)
		err = Publish(RefundQueue, []byte("123"))
		fmt.Println(err)
		wg.Done()
	}()
	fmt.Println("=======")
	wg.Wait()
}
