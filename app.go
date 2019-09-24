package main

import(
  ws "github.com/gorilla/websocket"
  "fmt"
  "net/http"
  "encoding/json"

  . "github.com/prestonjcrowe/binance-bot/orderbook"
  . "github.com/prestonjcrowe/binance-bot/bollinger"
)

func GetDepthSnapshot(symbol string, ob *OrderBook) BinanceSnapshot {
    base := "https://www.binance.com/api/v1/depth?symbol=%s&limit=1000"
    url := fmt.Sprintf(base, symbol)
    var target BinanceSnapshot
    resp, err := http.Get(url)
    if (err != nil) {
        panic("error retrieving snapshot");
    }
    defer resp.Body.Close()

    json.NewDecoder(resp.Body).Decode(&target)
    msg := BinanceDepth { Bids: target.Bids, Asks: target.Asks, FinalID: target.FinalID }
    ob.Update(msg)
    return target
}

func listenForOrders(url string, c chan BinanceDepth) {
    var wsDialer ws.Dialer
    wsConn, _, err := wsDialer.Dial(url, nil)
    var lastUpdatedID int
    if err != nil {
        println(err.Error())
    }

    for true {
        var msg BinanceDepth
        if err := wsConn.ReadJSON(&msg); err != nil {
            println(err.Error())
            break
        }
        c <- msg
        if (lastUpdatedID != 0 && lastUpdatedID + 1 != msg.FirstID) {
            panic("MISSED AN UPDATE")
        }
        lastUpdatedID = msg.FinalID
    }
}

func main() {
    // start listening for events in a go routine -> chan
    // after that, get the snapshot, then start consuming from chan
    // snapshot should return lastUpdateId
    var cw CandleWindow
    klineHistory := GetKLineHistory("BTCUSDT", "1m")

    for _, c := range klineHistory {
      cw.Add(c)
      bb := GetBollingerBands(cw, 20)
      bb.Print()
    }
    /*
    var ob OrderBook
    var url string = "wss://stream.binance.com:9443/ws/btcusdt@depth"
    msgChan := make(chan BinanceDepth)

    go listenForOrders(url, msgChan)
    GetDepthSnapshot("BTCUSDT", &ob)
    for true  {
        //ob.Asks.Print()
        msg := <-msgChan
        ob.Update(msg)
        fmt.Printf("Ask price: %s\n", ob.GetAsking().String())
    }*/
}
