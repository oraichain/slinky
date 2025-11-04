package indices

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	providertypes "github.com/dydxprotocol/slinky/providers/types"

	"github.com/dydxprotocol/slinky/oracle/config"
	"github.com/dydxprotocol/slinky/oracle/types"
)

var _ types.PriceAPIDataHandler = (*APIHandler)(nil)

// APIHandler implements the PriceAPIDataHandler interface for Polygon Indices.
type APIHandler struct {
	// api is the config for the Polygon Indices API.
	api config.APIConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewAPIHandler returns a new Polygon Indices PriceAPIDataHandler.
func NewAPIHandler(
	api config.APIConfig,
) (types.PriceAPIDataHandler, error) {
	if api.Name != Name {
		return nil, fmt.Errorf("expected api config name %s, got %s", Name, api.Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api config for %s is not enabled", Name)
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config for %s: %w", Name, err)
	}

	return &APIHandler{
		api:   api,
		cache: types.NewProviderTickers(),
	}, nil
}

// CreateURL returns the URL that is used to fetch data from the Polygon Indices API for the
// given tickers. The URL format is: ${baseURL}/indices?ticker.any_of=I:SPX,I:DJX,I:NQTHLMT
func (h *APIHandler) CreateURL(
	tickers []types.ProviderTicker,
) (string, error) {
	indicesIds := make([]string, len(tickers))
	for i, ticker := range tickers {
		indicesIds[i] = ticker.GetOffChainTicker()
		h.cache.Add(ticker)
	}

	baseURL := strings.TrimSuffix(h.api.Endpoints[0].URL, "/")
	return fmt.Sprintf("%s/indices?ticker.any_of=%s&%s=%s",
		baseURL,
		strings.Join(indicesIds, ","),
		h.api.Endpoints[0].Authentication.APIKeyHeader,
		h.api.Endpoints[0].Authentication.APIKey,
	), nil
}

// ParseResponse parses the response from the Polygon Indices API. The response is expected
// to contain multiple index prices.
func (h *APIHandler) ParseResponse(
	tickers []types.ProviderTicker,
	resp *http.Response,
) types.PriceResponse {
	// Parse the response.
	var result PolygonIndicesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToDecode),
		)
	}

	var (
		resolved   = make(types.ResolvedPrices)
		unresolved = make(types.UnResolvedPrices)
	)

	// Check if the response status is OK.
	if result.Status != "OK" {
		return types.NewPriceResponseWithErr(
			tickers,
			providertypes.NewErrorWithCode(
				fmt.Errorf("unexpected status: %s", result.Status),
				providertypes.ErrorAPIGeneral,
			),
		)
	}

	// Build a map of ticker -> result for quick lookup.
	resultMap := make(map[string]PolygonIndicesResult)
	for _, r := range result.Results {
		resultMap[r.Ticker] = r
	}

	// Process each requested ticker.
	for _, ticker := range tickers {
		offChainTicker := ticker.GetOffChainTicker()

		indexResult, ok := resultMap[offChainTicker]
		if !ok {
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(
					fmt.Errorf("no response for ticker %s", offChainTicker),
					providertypes.ErrorNoResponse,
				),
			}
			continue
		}

		resolved[ticker] = types.NewPriceResult(
			big.NewFloat(indexResult.Value),
			time.Now().UTC(),
		)
	}

	return types.NewPriceResponse(resolved, unresolved)
}
