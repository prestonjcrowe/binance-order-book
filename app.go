package main

import(
  ws "github.com/gorilla/websocket"
  "fmt"
  "net/http"
  "encoding/json"

  . "github.com/prestonjcrowe/coinbase-bot/orderbook"
)

func GetDepthSnapshot(symbol string) BinanceSnapshot {
    base := "https://www.binance.com/api/v1/depth?symbol=%s&limit=1000"
    url := fmt.Sprintf(base, symbol)
    var target BinanceSnapshot
    resp, err := http.Get(url)
    if (err != nil) {
        panic("error retrieving snapshot");
    }
    defer resp.Body.Close()

    json.NewDecoder(resp.Body).Decode(&target)
    return target
}

func main() {
    var wsDialer ws.Dialer
    var ob OrderBook
    var url string = "wss://stream.binance.com:9443/ws/btcusdt@depth"

    ob.Init()
    snapshot := GetDepthSnapshot("BTCUSDT")
    fmt.Printf("Snapshot size: %d\n", len(snapshot.Bids))

    wsConn, _, err := wsDialer.Dial(url, nil)
    if err != nil {
        println(err.Error())
    }

    for true {
        var msg BinanceDepth
        if err := wsConn.ReadJSON(&msg); err != nil {
            println(err.Error())
            break
        }
        msg.Print()
        ob.Update(msg)
    }
}
