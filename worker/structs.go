package worker

import (
	"encoding/json"
	"time"
)

// AskBid is a struct use to return ask and bid of request pair.
type AskBid struct {
	Success string  `json:"success"`
	Message string  `json:"message"`
	Market  string  `json:"market"`
	Ask     float64 `json:"ask"`
	Bid     float64 `json:"bid"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Volume  float64 `json:"volume"`
}

// MainAskBid1 this is use to get single request
type MainAskBid1 struct {
	Values []AskBid
}

type ExchangeInfo struct {
	ExchangeID   int64  `json:"ExchangeID"`
	ExchangeName string `json:"ExchangeName"`
}
type OrderInfoResponse struct {
	Result       string          `json:"result"`
	Message      string          `json:"message"`
	Details      json.RawMessage `json:"order_details"`
	ExchangeInfo json.RawMessage `json:"exchange_info"`
}

type orderDetail struct {
	Market            string  `json:"market"`
	OrderType         string  `json:"order_type"`
	ActualQuantity    float64 `json:"actual_quantity"`
	QuantityRemaining float64 `json:"QuantityRemaining"`
	ActualRate        float64 `json:"actual_rate"`
	OrderStatus       string  `json:"order_status"`
	Fee               float64 `json:"fee"`
	OrderDate         string  `json:"order_date"`
	Price             float64 `json:"price"`
	PricePerUnit      float64 `json:"pricePerUnit"`
	Reserved          float64 `json:"reserved"`
	Exchange          ExchangeInfo
}

//OrderResult holds result of an order
type OrderResult struct {
	Message     string
	OrderNumber string
}

type orderInfo2 struct {
	Message     string
	Market      string
	OrderType   string
	ActualQty   float64
	ActualRate  float64
	OrderStatus string
	Fee         float64
	OrderDate   string
}

//OrderRowID order row id
type OrderRowID struct {
	Message string
	RowID   int
}

//OrderInfo struct stores order info
type OrderInfo struct {
	Market            string  `json:"market"`
	QuantityRemaining float64 `json:"QuantityRemaining"`
	ActualQuantity    float64 `json:"actual_quantity"`
	ActualRate        float64 `json:"actual_rate"`
	OrderDate         string  `json:"order_date"`
	OrderStatus       string  `json:"order_status"`
	OrderType         string  `json:"order_type"`
	Price             float64 `json:"price"`
	PricePerUnit      float64 `json:"pricePerUnit"`
}

type dbJob struct {
	market                 string
	orderID                string
	accountID              int64
	exchangeID             int64
	apiKey                 string
	startTime              time.Time
	buyOrderTimeout        int64
	partailBuyTimeout      int64
	partailBuyTimeoutPl    float64
	partialBuyDetectedTime time.Time
	askBid                 float64
	jobID                  int64
	updateHash             string
	orderType              string
	actualRate             float64
	cost                   float64
	quantity               float64
	actualQuantity         float64
	buyCapital             float64
	tradeMode              string
	tradeProfile           string
	instanceID             int64
}

//InternalBalanceResponse returns Response of internalBalance
type InternalBalanceResponse struct {
	Result  string                `json:"result"`
	Message string                `json:"message"`
	Details InternalWalletBalance `json:"details"`
}

//InternalBalanceUpdateResponse returns Response of internalBalance Updates
type InternalBalanceUpdateResponse struct {
	Result     string  `json:"result"`
	Message    string  `json:"message"`
	NewBalance float64 `json:"new_balance"`
}

//InternalWalletBalance holds internal wallet balance
type InternalWalletBalance struct {
	AccountID      int64   `json:"account_id"`
	WalletType     string  `json:"wallet_type"`
	ExchangeID     int64   `json:"exchange_id"`
	CoinName       string  `json:"coin_name"`
	Available      float64 `json:"available"`
	Pending        float64 `json:"pending"`
	Reserved       float64 `json:"reserved"`
	Used           float64 `json:"used"`
	Total          float64 `json:"total"`
	MainCapital    float64 `json:"main_capital"`
	CurrentCapital float64 `json:"current_capital"`
}

type balances struct {
	Capital   float64
	Pending   float64
	Available float64
	Reserved  float64
	Used      float64
	Msg       string
}

//OrderResponse struct stores response of Orders
type OrderResponse struct {
	Result      string `json:"result"`
	Message     string `json:"message"`
	OrderNumber string `json:"order_number"`
}
