package amqppool

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

type Delivery amqp.Delivery

var chPool *ConnPool

func Configure(host string, port string, user string, password string) (err error) {

	chPool = NewConnPool(user, password, host, port, 20, 150, 20)
	return nil
}

//使用 consume 的正确姿势是 for range  ，重连后需要 重新获取 chan Delivery
func Consume(queue string) (<-chan amqp.Delivery, error) {
	ch, err := chPool.Get()
	if err != nil {
		return nil, err
	}

	//defer chPool.Release(ch) //消费者 不释放

	return ch.cha.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func Publish(queue string, msg []byte) error {
	var exchange, key string
	ch, err := chPool.Get()
	if err != nil {
		return err
	}
	defer chPool.Release(ch)

	exchange = ExchangeName
	for _, q := range DefaultBinds[0].Queues {
		if q.Name == queue {
			key = q.BindKey
			break
		}
	}

	if key == "" {
		return errors.New("queue " + queue + " is not exist")
	}

	// err = ch.cha.Tx()
	// if err != nil {
	// 	return err
	// }
	if ch == nil || ch.cha == nil {
		panic(fmt.Errorf("%v ,%v", ch, ch.cha))
	}
	err = ch.cha.Publish(
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: 2,
			ContentType:  "text/plain",
			Body:         msg,
		},
	)
	if err != nil {
		chPool.errCloseChannel(ch)
		return err
	}
	// if err != nil {
	// 	ch.cha.TxRollback()
	// 	return err
	// }
	// err = ch.cha.TxCommit()
	// if err != nil {
	// 	ch.cha.TxRollback()
	// 	return err
	// }
	if confirmed := <-ch.confirms; confirmed.Ack {
		//log.Infof("confirmed delivery with delivery tag: %d\n", confirmed.DeliveryTag)
	} else {
		return fmt.Errorf("failed delivery of delivery queue:%s msg:%s", queue, string(msg))
	}
	return err
}
