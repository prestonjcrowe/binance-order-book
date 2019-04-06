package main

import(
  ws "github.com/gorilla/websocket"
  coinbasepro "github.com/preichenberger/go-coinbasepro"
  //"github.com/shopspring/decimal"

  "log"
)

func main() {
  var wsDialer ws.Dialer
  wsConn, _, err := wsDialer.Dial("wss://ws-feed.pro.coinbase.com", nil)
  if err != nil {
    println(err.Error())
  }

  subscribe := coinbasepro.Message{
    Type:      "subscribe",
    Channels: []coinbasepro.MessageChannel{
      coinbasepro.MessageChannel{
        Name: "heartbeat",
        ProductIds: []string{
          "BTC-USD",
        },
      },
      coinbasepro.MessageChannel{
        Name: "level2",
        ProductIds: []string{
          "BTC-USD",
        },
      },
    },
  }
  if err := wsConn.WriteJSON(subscribe); err != nil {
    println(err.Error())
  }

  for true {
    /*_, message, err := wsConn.ReadMessage()
    if err != nil {
        println(err.Error())
    }
    log.Printf(": %s\n", message)*/

    message := coinbasepro.Message{}
    if err := wsConn.ReadJSON(&message); err != nil {
      println(err.Error())
      break
    }

    changes := message.Changes
    for _, change := range changes {
        log.Printf(": %s | %s | %s\n", change.Price, change.Size, change.Side)
    }
  }
}
