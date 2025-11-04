package polygon

import (
	"fmt"
	"time"

	"github.com/dydxprotocol/slinky/oracle/types"
	"github.com/dydxprotocol/slinky/pkg/math"
)

func (h *WebSocketHandler) parseIndicesResponse(resp []IndicesResponse) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	for _, r := range resp {
		// Convert the price to a big.Float.
		price := math.Float64ToBigFloat(r.Value)

		ticker, ok := h.cache.FromOffChainTicker(r.Ticker)
		if !ok {
			return types.NewPriceResponse(resolved, unresolved), fmt.Errorf("unknown ticker %s", r.Ticker)
		}

		timestamp := time.Unix(r.TimeStamp, 0)
		resolved[ticker] = types.NewPriceResult(price, timestamp)
	}

	return types.NewPriceResponse(resolved, unresolved), nil
}
