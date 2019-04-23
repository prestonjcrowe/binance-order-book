package orderbook
import (
    "github.com/shopspring/decimal"
    "fmt"
)

type OrderType int

const (
    BID OrderType = 0
    ASK OrderType = 1
)

type OrderBook struct {
    Bids map[string][]decimal.Decimal
    Asks map[string][]decimal.Decimal
}

type BinanceDepth struct {
    Type    string     `json:"e"`
    Time    int        `json:"E"`
    Symbol  string     `json:"s"`
    FirstID int        `json:"U"`
    FinalID int        `json:"u"`
    Bids    [][]string `json:"b"`
    Asks    [][]string `json:"a"`
}

type BinanceSnapshot struct {
    FinalID int            `json:"lastUpdateId"`
    Bids    [][]string     `json:"bids"`
    Asks    [][]string     `json:"asks"`
}

func (ob *OrderBook) Init() {
    ob.Bids = make(map[string][]decimal.Decimal)
    ob.Asks = make(map[string][]decimal.Decimal)
}

func (ob *OrderBook) Update(bd BinanceDepth) {
    for _, bid := range bd.Bids {
        price := bid[0]
        q, err := decimal.NewFromString(bid[1]) // get quantity
        if (err != nil) {
            panic("error parsing quantity");
        }

        if (q == decimal.Zero) {
            delete(ob.Bids, price)
        } else {
            ob.Bids[price]  = append(ob.Bids[price], q)
        }
    }

    for _, ask := range bd.Asks {
        price := ask[0]
        q, err := decimal.NewFromString(ask[1]) // get quantity
        if (err != nil) {
            panic("error parsing quantity");
        }

        if (q == decimal.Zero) {
            delete(ob.Asks, price)
        } else {
            ob.Asks[price]  = append(ob.Asks[price], q)
        }
    }
}

func (bd *BinanceDepth) Print() {
    fmt.Printf("first: %d | last: %d\n", bd.FirstID, bd.FinalID)
    for _, bid := range bd.Bids {
        fmt.Printf("\tprice: %s | quantity: %s\n", bid[0], bid[1]);
    }

    for _, ask := range bd.Asks {
        fmt.Printf("\tprice: %s | quantity: %s\n", ask[0], ask[1]);
    }
}
