package polygon

import (
	"fmt"
	"time"

	"github.com/dydxprotocol/slinky/oracle/types"
	"github.com/dydxprotocol/slinky/pkg/math"
)

func (h *WebSocketHandler) parseIndicesResponse(resp IndicesResponse) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	// Convert the price to a big.Float.
	price := math.Float64ToBigFloat(resp.Value)

	ticker, ok := h.cache.FromOffChainTicker(resp.Ticker)
	if !ok {
		return types.NewPriceResponse(resolved, unresolved), fmt.Errorf("unknown ticker %s", resp.Ticker)
	}

	timestamp := time.Unix(resp.TimeStamp, 0)

	resolved[ticker] = types.NewPriceResult(price, timestamp)
	return types.NewPriceResponse(resolved, unresolved), nil
}
