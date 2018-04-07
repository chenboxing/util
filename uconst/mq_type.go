package uconst

import "time"

type TradeMsg struct {
	UlMchID       string    `json:"ul_mch_id"`
	OutTradeNo    string    `json:"out_trade_no"`
	TransactionID string    `json:"transaction_id"` //添加微信交易单号
	AppID         string    `json:"appid"`          //TODO 去掉
	Channel       string    `json:"channel"`
	TradeType     string    `json:"trade_type"`   //TODO 去掉
	PaymentCode   string    `json:"payment_code"` //定位费率字段改为payment_code
	TotalFee      int64     `json:"total_fee"`
	CompleteAt    time.Time `json:"complete_at"` //交易完成时间时间戳
}

type RefundMsg struct {
	UlMchID       string    `json:"ul_mch_id"`
	OutRefundNo   string    `json:"out_refund_no"`
	OutTradeNo    string    `json:"out_trade_no"`
	TransactionID string    `json:"transaction_id"` //添加微信交易单号
	Channel       string    `json:"channel"`
	RefundFee     int64     `json:"refund_fee"`
	CompleteAt    time.Time `json:"complete_at"` //交易完成时间时间戳
}

type WithdrawMsg struct {
	OutWithdrawNo string           `json:"out_withdraw_no"` //ul对外提现单号
	MchWithdrawNo string           `json:"mch_withdraw_no"` //商户提现单号
	UlMchID       string           `json:"ul_mch_id"`       //商户ID
	TotalFee      int64            `json:"total_fee"`       //总扣款额
	HandlingFee   int64            `json:"handling_fee"`    //手续费
	Channel       string           `json:"channel"`         //通道
	Fees          map[string]int64 `json:"fees"`            //按日期的冻结金额
}

type BookMsg struct {
	JobID            string    `json:"job_id"`
	Channel          string    `json:"channel"`
	MchID            int64     `json:"mch_id"`
	DtID             int64     `json:"dt_id"`
	TradeFee         int64     `json:"trade_fee"`          //+
	RefundFee        int64     `json:"refund_fee"`         //+
	WithdrawFee      int64     `json:"withdraw_fee"`       //+
	MchProfit        int64     `json:"mch_profit"`         //交易＋ 退款－ 提现－
	MchTradeFee      int64     `json:"mch_trade_profit"`   //交易－ 退款＋
	MchD0Fee         int64     `json:"mch_d0_profit"`      //提现－
	DtTradeProfit    int64     `json:"dt_trade_profit"`    //交易＋ 退款－
	DtD0Profit       int64     `json:"dt_d0_profit"`       //提现＋
	BpTradeProfit    int64     `json:"bp_trade_profit"`    //交易＋ 退款－
	BpD0Profit       int64     `json:"bp_d0_profit"`       //提现＋
	PayChannelProfit int64     `json:"pay_channel_profit"` //交易＋ 退款－
	BankD0Cost       int64     `json:"bank_d0_const"`      //银行提现成本 提现费率部分＋单笔提现手续费 +
	TradeCount       int       `json:"trade_count"`
	RefundCount      int       `json:"refund_count"`
	WithdrawCount    int       `json:"withdraw_count"`
	CompleteAt       time.Time `json:"complete_at"`
}
