package bitrue

import (
	"encoding/json"
	"github.com/ericlagergren/decimal"
	"github.com/spf13/cast"
	"log"
)

const https = "https://www.bitrue.com"

var minAmount = make(map[string]float64)

type SymbolReturn struct {
	Symbols *[]SymbolData `json:"symbols"`
}

type SymbolData struct {
	Symbol         string
	Status         string
	BasePrecision  int `json:"baseAssetPrecision"`
	QuotePrecision int
}

type KlineData struct {
	Id     int64
	Amount float64
	Vol    float64
	High   float64
	Low    float64
	Close  float64
	Open   float64
}

type ReqKline struct {
	EventRep string `json:"event_rep"`
	Symbol   string `json:"cb_id"`
	Channel  string
	Ts       int64
	Data     []KlineData
}

type Ticker struct {
	Amount float64
}

type Kline struct {
	Channel string
	Ts      int64
	Data    KlineData `json:"tick"`
}

type DepthData struct {
	Bids [][2]float64
	Asks [][2]float64 `json:"asks"`
}

func (depthData *DepthData) UnmarshalJSON(data []byte) error {
	dep := struct {
		Bids [][2]*decimal.Big `json:"buys"`
		Asks [][2]*decimal.Big `json:"asks"`
	}{}
	//depth.Bids = make([][]float64, 0)
	//depth.Asks= make([][]float64, 0)

	err := json.Unmarshal(data, &dep)

	for i := 0; i < len(dep.Bids); i++ {
		a, _ := dep.Bids[i][0].Float64()
		b, _ := dep.Bids[i][1].Float64()
		depthData.Bids = append(depthData.Bids, [2]float64{a, b})
	}

	for i := 0; i < len(dep.Asks); i++ {
		a, _ := dep.Asks[i][0].Float64()
		b, _ := dep.Asks[i][1].Float64()
		depthData.Asks = append(depthData.Asks, [2]float64{a, b})
	}

	return err
}

type DepthWs struct {
	Channel string
	Ts      int64
	Data    *DepthData `json:"tick"`
}

// [price ,amount]
type Depth struct {
	LastUpdateId int64
	Bids         [][2]*decimal.Big
	Asks         [][2]*decimal.Big
}

func (depth *Depth) Section() float64 {
	askPrice1, _ := depth.Asks[0][0].Float64()
	bidPrice1, _ := depth.Bids[0][0].Float64()
	return askPrice1 - bidPrice1
}

func (depth *Depth) BidsPrice(row int) float64 {
	price, _ := depth.Bids[row][0].Float64()
	return price
}

func (depth *Depth) AsksPrice(row int) float64 {
	price, _ := depth.Asks[row][0].Float64()
	return price
}

func (depth *Depth) DepthBidsAmountAll(row int) float64 {
	amount := 0.0
	for i := 0; i <= row; i++ {
		a, _ := depth.Bids[i][1].Float64()
		amount += a
	}
	return amount
}

func (depth *Depth) DepthAsksAmountAll(row int) float64 {
	amount := 0.0
	for i := 0; i <= row; i++ {
		a, _ := depth.Asks[i][1].Float64()
		amount += a
	}
	return amount
}

type PriceTicker struct {
	Symbol string
	Price  *decimal.Big
}

type BookTicker struct {
	Symbol   string
	BidPrice *decimal.Big
	BidQty   *decimal.Big
	AskPrice *decimal.Big
	AskQty   *decimal.Big
}

type OrderData struct {
	Symbol      string
	OrderId     int64 `json:",string"`
	Price       decimal.Big
	OrigQty     decimal.Big
	ExecutedQty string
	Sid         string
	Type        string
	Status      string
}

func (order *OrderData) Filled() float64 {
	return cast.ToFloat64(order.ExecutedQty)
}

func (order *OrderData) String() string {
	data, _ := json.Marshal(order)
	return string(data)
}

type BalanceData struct {
	Currency string `json:"asset"`
	Free     decimal.Big
	Locked   decimal.Big
}

func (bd *BalanceData) GetFree() float64 {
	f, ok := bd.Free.Float64()
	if !ok {
		log.Println("cast to float err:", bd.Free.String())
	}
	return f
}

func (bd *BalanceData) GetLock() float64 {
	f, ok := bd.Locked.Float64()
	if !ok {
		log.Println("cast to float err:", bd.Free.String())
	}
	return f
}

type Balance struct {
	UpdateTime int64
	Balances   []*BalanceData
}

type DeleteReturn struct {
	Symbol  string
	OrderId int64
}
