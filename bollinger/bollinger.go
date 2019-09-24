package orderbook
import (
    "github.com/shopspring/decimal"
    "net/http"
    "fmt"
    "encoding/json"
)

type BinanceKLine struct {
    Candle CandleStick `json:"k"`
}

type CandleStick struct {
    Open   decimal.Decimal `json:"o"`
    Close  decimal.Decimal `json:"c"`
    High   decimal.Decimal `json:"h"`
    Low    decimal.Decimal `json:"l"`
    Volume decimal.Decimal `json:"v"`
}

type BollingerBands struct {
    Upper decimal.Decimal
    Lower decimal.Decimal
    SMA   decimal.Decimal
}

/*
func RSIGenerator(cw *CandleWindow, n int) func(cw *CandleWindow, n int) decimal.Decimal {
    avgGain = GetAverageGain(cw, n)
    avgLoss = GetAverageLoss(cw, n)
    return func(cw *CandleWindow, n int) decimal.Decimal {
        i++
        return i
    }
}
*/
type KLineResponse [][]interface{};

type CandleWindow struct {
    Window []CandleStick
}

func (cw *CandleWindow) Add(c CandleStick) {
    if (len(cw.Window) < 20) {
      cw.Window = append(cw.Window, c)
    } else {
      cw.Window = append(cw.Window, c)[1:]
    }
}

func (bb *BollingerBands) Print() {
    fmt.Printf("Upper: %s Lower: %s Average: %s\n", bb.Upper, bb.Lower, bb.SMA)
}

func GetKLineHistory(symbol string, interval string) []CandleStick {
  var kLineData KLineResponse
  var res []CandleStick

  base := "https://api.binance.com/api/v1/klines?symbol=%s&interval=%s"
  url := fmt.Sprintf(base, symbol, interval)
  resp, err := http.Get(url)

  if err = json.NewDecoder(resp.Body).Decode(&kLineData); err != nil {
    panic("couldn't decode KLINE response")
}

  for _, kline := range kLineData {
    o, _ := decimal.NewFromString(kline[1].(string))
    h, _ := decimal.NewFromString(kline[2].(string))
    l, _ := decimal.NewFromString(kline[3].(string))
    c, _ := decimal.NewFromString(kline[4].(string))
    v, _ := decimal.NewFromString(kline[5].(string))

    res = append(res, CandleStick{ Open:o, High:h, Low:l, Close:c, Volume: v })
  }
  return res
}

func GetBollingerBands(cw CandleWindow, n int) BollingerBands {
    if (len(cw.Window) != n) {
        return BollingerBands{}
    }

    sma := SimpleMovingAverage(cw.Window, n)
    return BollingerBands{
        Upper: UpperBand(sma, cw.Window, n),
        Lower: LowerBand(sma, cw.Window, n),
        SMA  : sma }
}

func UpperBand(sma decimal.Decimal, vals []CandleStick, n int) decimal.Decimal {
    variance := Variance(sma, vals, n)
    return sma.Add(variance.Mul(decimal.NewFromFloat(2)))
}

func LowerBand(sma decimal.Decimal, vals []CandleStick, n int) decimal.Decimal {
    variance := Variance(sma, vals, n)
    return sma.Sub(variance.Mul(decimal.NewFromFloat(2)))
}

func TypicalPrice(c CandleStick) decimal.Decimal {
    //return c.High.Add(c.Low).Add(c.Close).Div(decimal.NewFromFloat(3))
    return c.Close
}

func Variance(sma decimal.Decimal, vals []CandleStick, n int) decimal.Decimal {
    var res decimal.Decimal
    for _, val := range vals {
        diff := TypicalPrice(val).Sub(sma)
        sq   := diff.Pow(decimal.NewFromFloat(2))
        res = res.Add(sq)
    }

    return res.Div(decimal.New(int64(n - 1), 1))
}

func SimpleMovingAverage(vals []CandleStick, n int) decimal.Decimal {
    result := decimal.NewFromFloat(0)
    for _, val := range vals {
        result = result.Add(TypicalPrice(val))
    }
    result = result.Div(decimal.NewFromFloat(float64(n)))
    return result
}
