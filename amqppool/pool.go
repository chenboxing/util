package amqppool

import (
	"fmt"
	"net"
	"sync"
	"time"

	"errors"

	"github.com/streadway/amqp"
)

type Connection struct {
	conn             *amqp.Connection
	numOpenedChannel int
}

func (conn *Connection) openChannel() (*Channel, error) {

	if conn.conn == nil {
		return nil, errors.New("error conn is null")
	}

	amqpCha, err := conn.conn.Channel()
	if err != nil {
		return nil, err
	}
	cha := &Channel{
		conn: conn,
		cha:  amqpCha,
	}

	for _, b := range DefaultBinds {
		if err := amqpCha.ExchangeDeclare(
			b.Exchange.Name,
			b.Exchange.Kind,
			b.Exchange.Durable,
			b.Exchange.AutoDelete,
			b.Exchange.Internal,
			b.Exchange.NoWait,
			b.Exchange.Args,
		); err != nil {
			return nil, err
		}

		for _, q := range b.Queues {
			queue, err := amqpCha.QueueDeclare(
				q.Name,
				q.Durable,
				q.AutoDelete,
				q.Exclusive,
				q.NoWait,
				q.Args,
			)

			if err != nil {
				return nil, err
			}

			if err = amqpCha.QueueBind(queue.Name, q.BindKey, b.Exchange.Name, q.BindNoWait, q.BindArgs); err != nil {
				return nil, err
			}
		}
	}

	//设置channel confirm
	if err := amqpCha.Confirm(false); err != nil {
		return nil, fmt.Errorf("Channel could not be put into confirm mode: %s", err)
	}
	cha.confirms = amqpCha.NotifyPublish(make(chan amqp.Confirmation, 1))

	// err = cha.Bind("test", "topic", "test1", "key")
	// if err != nil {
	// 	return nil, err
	// }

	conn.numOpenedChannel++
	return cha, nil
}

func (conn *Connection) close() error {
	if conn.numOpenedChannel > 0 {
		return errors.New("conn has opened channel")
	}
	if conn.conn != nil {
		conn.conn.Close()
	}
	conn.conn = nil
	return nil
}

type Channel struct {
	conn     *Connection
	cha      *amqp.Channel
	confirms chan amqp.Confirmation //确认 送到了mq broker 才算成功
}

func (cha *Channel) close() {
	cha.cha.Close()
	cha.conn.numOpenedChannel--
	cha.cha = nil
	cha.conn = nil
}

func (cha *Channel) Bind(exchange, kind, queue, key string) error {
	if err := cha.cha.ExchangeDeclare(
		exchange,
		kind,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	q, err := cha.cha.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	if err := cha.cha.QueueBind(q.Name, key, exchange, false, nil); err != nil {
		return err
	}
	return nil
}

type ConnPool struct {
	conns []*Connection
	Dial  func() (*amqp.Connection, error) //建立连接的方法

	l sync.Mutex

	// idle channels
	idleChas chan *Channel

	MaxChannelsPerConnection int
	MaxConnections           int
	MaxIdleChannels          int
	//closed                   bool //暂时不设置关闭pool
}

func NewConnPool(user, password, host, port string, MaxConnections, MaxIdleChannels, MaxChannelsPerConnection int) *ConnPool {
	return &ConnPool{
		Dial: func() (*amqp.Connection, error) {
			return amqp.DialConfig(
				fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port),
				amqp.Config{
					Heartbeat: 10 * time.Second,
					Dial: func(network, addr string) (net.Conn, error) {
						return net.DialTimeout(network, addr, 3*time.Second) //3秒设置超时
					}})
		},
		MaxConnections:           MaxConnections,
		MaxIdleChannels:          MaxIdleChannels,
		MaxChannelsPerConnection: MaxChannelsPerConnection,
		conns:    make([]*Connection, 0, MaxConnections),
		idleChas: make(chan *Channel, MaxIdleChannels),
		//closed:   false,
	}
}

func (pool *ConnPool) Get() (*Channel, error) {
	pool.l.Lock()
	defer pool.l.Unlock()
	select {
	case ch := <-pool.idleChas:
		return ch, nil
	default:
	}
	conn, err := pool.getConn() // pool.idleChas 里面没有就新建
	if err != nil {
		return nil, err
	}
	return conn.openChannel()
}

func (pool *ConnPool) Release(cha *Channel) {
	if cha == nil || cha.cha == nil || cha.conn == nil {
		return
	}
	pool.l.Lock()
	defer pool.l.Unlock()
	select {
	case pool.idleChas <- cha:
	default:
		pool.probeCloseChannel(cha) // pool ch 满了就close
	}
}

func (pool *ConnPool) probeCloseChannel(cha *Channel) {
	conn := cha.conn
	cha.close()
	if conn.numOpenedChannel == 0 && len(pool.conns) > pool.MaxConnections {
		pool.removeConn(conn)
	}
}

//链接错误了，需要移除
func (pool *ConnPool) errCloseChannel(cha *Channel) {
	conn := cha.conn
	cha.close()
	pool.removeConn(conn)
}

func (pool *ConnPool) getConn() (*Connection, error) {
	if len(pool.conns) > 0 {
		for i := 0; i < len(pool.conns); i++ {
			if pool.conns[i].numOpenedChannel < pool.MaxChannelsPerConnection {
				return pool.conns[i], nil
			}
		}
	}

	if pool.MaxConnections > 0 && len(pool.conns) >= pool.MaxConnections {
		return nil, errors.New("error too may conn")
	}

	amqpConn, err := pool.Dial()
	if err != nil {
		return nil, err
	}
	conn := &Connection{conn: amqpConn, numOpenedChannel: 0}
	pool.conns = append(pool.conns, conn)
	return conn, nil
}

func (pool *ConnPool) removeConn(conn *Connection) error {
	foundIdx := -1
	for i := 0; i < len(pool.conns); i++ {
		if conn == pool.conns[i] {
			foundIdx = i
			break
		}
	}
	// delete from pool
	if foundIdx > -1 {
		copy(pool.conns[foundIdx:], pool.conns[foundIdx+1:])
		pool.conns[len(pool.conns)-1] = nil
		pool.conns = pool.conns[:len(pool.conns)-1]
	}

	return conn.close()
}
