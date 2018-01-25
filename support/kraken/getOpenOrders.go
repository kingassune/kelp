package kraken

import (
	"fmt"

	"github.com/lightyeario/kelp/support/exchange/assets"
	"github.com/lightyeario/kelp/support/exchange/dates"
	"github.com/lightyeario/kelp/support/exchange/number"
	"github.com/lightyeario/kelp/support/exchange/orderbook"
)

// GetOpenOrders impl.
func (k krakenExchange) GetOpenOrders() (map[assets.TradingPair][]orderbook.OpenOrder, error) {
	openOrdersResponse, e := k.api.OpenOrders(map[string]string{})
	if e != nil {
		return nil, e
	}

	m := map[assets.TradingPair][]orderbook.OpenOrder{}
	for _, o := range openOrdersResponse.Open {
		pair, e := assets.FromString(k.assetConverter, o.Description.AssetPair)
		if e != nil {
			return nil, e
		}
		if _, ok := m[*pair]; !ok {
			m[*pair] = []orderbook.OpenOrder{}
		}
		if _, ok := m[assets.TradingPair{Base: pair.Quote, Quote: pair.Base}]; ok {
			return nil, fmt.Errorf("open orders are listed with repeated base/quote pairs for %s", *pair)
		}

		m[*pair] = append(m[*pair], orderbook.OpenOrder{
			Order: orderbook.Order{
				Pair:        pair,
				OrderAction: orderbook.OrderActionFromString(o.Description.Type),
				OrderType:   orderbook.OrderTypeFromString(o.Description.OrderType),
				Price:       number.FromFloat(o.Price),
				Volume:      number.MustFromString(o.Volume),
				Timestamp:   dates.MakeTimestamp(int64(o.OpenTime)),
			},
			ID:             o.ReferenceID,
			StartTime:      dates.MakeTimestamp(int64(o.StartTime)),
			ExpireTime:     dates.MakeTimestamp(int64(o.ExpireTime)),
			VolumeExecuted: number.FromFloat(o.VolumeExecuted),
		})
	}
	return m, nil
}