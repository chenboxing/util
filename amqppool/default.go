package amqppool

import "github.com/streadway/amqp"

const (
	ExchangeName          = "clear_exchange"
	ExchangeKind          = "direct" //每个队列一个key使用完全匹配即可 direct
	TradeQueue            = "trade_manage"
	TradeQueueKey         = "trade"
	RefundQueue           = "refund_manage"
	RefundQueueKey        = "refund"
	WithdrawQueue         = "d0withdraw_manage"
	WithdrawQueueKey      = "d0withdraw"
	TradeClearQueue       = "trade_clear_manage"
	TradeClearQueueKey    = "trade_clear"
	RefundClearQueue      = "refund_clear_manage"
	RefundClearQueueKey   = "refund_clear"
	WithdrawClearQueue    = "withdraw_clear_manage"
	WithdrawClearQueueKey = "withdraw_clear"
	BookQueue             = "book_manage"
	BookQueueKey          = "book"
)

type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type Queue struct {
	BindKey    string
	BindNoWait bool
	BindArgs   amqp.Table
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type Bind struct {
	Exchange *Exchange
	Queues   []*Queue
}

var (
	DefaultBinds = []*Bind{
		&Bind{
			Exchange: &Exchange{
				Name:       ExchangeName,
				Kind:       ExchangeKind,
				Durable:    true,
				AutoDelete: false,
				Internal:   false,
				NoWait:     false,
				Args:       nil,
			},
			Queues: []*Queue{
				&Queue{
					BindKey:    TradeQueueKey,
					BindNoWait: false,
					BindArgs:   nil,
					Name:       TradeQueue,
					Durable:    true,
					AutoDelete: false,
					Exclusive:  false,
					NoWait:     false,
					Args:       nil,
				},
				&Queue{
					BindKey:    RefundQueueKey,
					BindNoWait: false,
					BindArgs:   nil,
					Name:       RefundQueue,
					Durable:    true,
					AutoDelete: false,
					Exclusive:  false,
					NoWait:     false,
					Args:       nil,
				},
				&Queue{
					BindKey:    WithdrawQueueKey,
					BindNoWait: false,
					BindArgs:   nil,
					Name:       WithdrawQueue,
					Durable:    true,
					AutoDelete: false,
					Exclusive:  false,
					NoWait:     false,
					Args:       nil,
				},
				// &Queue{
				// 	BindKey:    TradeClearQueueKey,
				// 	BindNoWait: false,
				// 	BindArgs:   nil,
				// 	Name:       TradeClearQueue,
				// 	Durable:    true,
				// 	AutoDelete: false,
				// 	Exclusive:  false,
				// 	NoWait:     false,
				// 	Args:       nil,
				// },
				// &Queue{
				// 	BindKey:    RefundClearQueueKey,
				// 	BindNoWait: false,
				// 	BindArgs:   nil,
				// 	Name:       RefundClearQueue,
				// 	Durable:    true,
				// 	AutoDelete: false,
				// 	Exclusive:  false,
				// 	NoWait:     false,
				// 	Args:       nil,
				// },
				// &Queue{
				// 	BindKey:    WithdrawClearQueueKey,
				// 	BindNoWait: false,
				// 	BindArgs:   nil,
				// 	Name:       WithdrawClearQueue,
				// 	Durable:    true,
				// 	AutoDelete: false,
				// 	Exclusive:  false,
				// 	NoWait:     false,
				// 	Args:       nil,
				// },
				&Queue{
					BindKey:    BookQueueKey,
					BindNoWait: false,
					BindArgs:   nil,
					Name:       BookQueue,
					Durable:    true,
					AutoDelete: false,
					Exclusive:  false,
					NoWait:     false,
					Args:       nil,
				},
			},
		},
	}
)
