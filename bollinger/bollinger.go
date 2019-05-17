package orderbook
import (
    "github.com/shopspring/decimal"
//    "fmt"
)

type BinanceKLine struct {
    Candle CandleStick `json:"k"`
}

type CandleStick struct {
    Open  decimal.Decimal  `json:"o"`
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

func GetBollingerBands(vals []decimal.Decimal, n int) BollingerBands {
    if (len(vals) != n) {
        return BollingerBands{}
    }

    sma := SimpleMovingAverage(vals, n)
    return BollingerBands{
        Upper: UpperBand(sma, vals, n),
        Lower: LowerBand(sma, vals, n),
        SMA  : sma }
}

func UpperBand(sma decimal.Decimal, vals []decimal.Decimal, n int) decimal.Decimal {
    variance := Variance(sma, vals, n)
    return sma.Add(variance.Mul(decimal.NewFromFloat(2)))
}

func LowerBand(sma decimal.Decimal, vals []decimal.Decimal, n int) decimal.Decimal {
    variance := Variance(sma, vals, n)
    return sma.Sub(variance.Mul(decimal.NewFromFloat(2)))
}

func TypicalPrice(c CandleStick) decimal.Decimal {
    return c.High.Add(c.Low).Add(c.Close).Div(decimal.NewFromFloat(3))
}

func Variance(sma decimal.Decimal, vals []decimal.Decimal, n int) decimal.Decimal {
    var res decimal.Decimal
    for _, val := range vals {
        diff := val.Sub(sma)
        sq   := diff.Pow(decimal.NewFromFloat(2))
        res = res.Add(sq)
    }

    return res.Div(decimal.New(int64(n - 1), 1))
}

func SimpleMovingAverage(vals []decimal.Decimal, n int) decimal.Decimal {
    result := decimal.NewFromFloat(0)
    for _, val := range vals {
        result = result.Add(val)
    }
    result = result.Div(decimal.NewFromFloat(float64(n)))
    return result
}
