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
    Asks OrderList
    Bids OrderList
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

func (ob *OrderBook) GetAsking() decimal.Decimal{
    return ob.Asks.Head.Price
}

func (ob *OrderBook) GetBidding() decimal.Decimal{
    return ob.Bids.Tail.Price
}

func (ob *OrderBook) Update(bd BinanceDepth) {
    for _, bid := range bd.Bids {
        p, err_1 := decimal.NewFromString(bid[0]) // get price
        q, err_2 := decimal.NewFromString(bid[1]) // get quantity

        if (err_1 != nil || err_2 != nil) {
            panic("error parsing bids")
        }

        if (q.IsZero()) {
            ob.Bids.Remove(p)
        } else {
            ob.Bids.Insert(p, q)
        }
    }

    for _, ask := range bd.Asks {
        p, err_1 := decimal.NewFromString(ask[0]) // get price
        q, err_2 := decimal.NewFromString(ask[1]) // get quantity

        if (err_1 != nil || err_2 != nil) {
            panic("error parsing asks")
        }

        if (q.IsZero()) {
            ob.Asks.Remove(p)
        } else {
            ob.Asks.Insert(p, q)
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

type OrderList struct {
    Head *OrderNode
    Tail *OrderNode
    Size int
}

type OrderNode struct {
    Price    decimal.Decimal
    Quantity []decimal.Decimal
    Next     *OrderNode
    Prev     *OrderNode
}

func (ll *OrderList) Insert(price decimal.Decimal, quanity decimal.Decimal) {
    node := ll.FindNode(price)
    if (node != nil) {
        node.Quantity = append(node.Quantity, quanity)
        return
    }

    curr := ll.Head
    newNode := &OrderNode { Price : price }
    newNode.Quantity = append(newNode.Quantity, quanity)
    for curr != nil {
        // insert before curr, sorted ascending
        if (price.Cmp(curr.Price) <= 0) {
            newNode.Prev = curr.Prev
            newNode.Next = curr
            curr.Prev = newNode
            if (newNode.Prev == nil) {
                ll.Head = newNode
            }
            ll.Size++
            return
        }
        curr = curr.Next
    }
    // insert at tail
    if (ll.Tail == nil) {
        ll.Tail = newNode
        ll.Head = newNode
        newNode.Next = nil
        newNode.Prev = nil
    } else {
        newNode.Prev = ll.Tail
        newNode.Next = nil;
        ll.Tail.Next = newNode
        ll.Tail = newNode
        ll.Size++
    }
    return
}

func (ll *OrderList) FindNode(price decimal.Decimal) *OrderNode {
    curr := ll.Head
    for curr != nil {
        if (price.Cmp(curr.Price) == 0) {
            return curr
        }
        curr = curr.Next
    }
    return nil;
}

func (ll *OrderList) Remove(price decimal.Decimal) bool {
    node := ll.FindNode(price)
    if (node == nil) {
        return false
    }
    fmt.Printf("Couldnt find node %s, remove failed\n", price.String())

    if(node.Prev != nil && node.Next != nil) {
        node.Prev.Next = node.Next
        node.Next.Prev = node.Prev
    } else if (node.Next == nil) {
        ll.Tail = node.Prev
        ll.Tail.Next = nil
    } else {
        ll.Head = node.Next
        ll.Head.Prev = nil
    }
    ll.Size--
    return true
}

func (ll *OrderList) Print() {
    if (ll.Head == nil) { return }
    curr := ll.Head
    fmt.Printf("%s(%d)", curr.Price.String(), len(curr.Quantity))
    for curr.Next != nil {
        fmt.Printf(", %s(%d)", curr.Next.Price.String(), len(curr.Quantity))
        curr = curr.Next
    }
    fmt.Printf("\n")
}
