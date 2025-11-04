package indices

import (
	"time"

	"github.com/dydxprotocol/slinky/oracle/config"
)

const (
	// Name is the name of the Polygon Indices API provider.
	Name = "polygon_indices_api"

	// URL is the root URL for the Polygon Indices API.
	URL = "https://api.massive.com/v3/snapshot"
)

// DefaultAPIConfig is the default configuration for querying indices
// on the Polygon Indices API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           false,
	Enabled:          true,
	Timeout:          10 * time.Second,
	Interval:         20 * time.Second,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints:        []config.Endpoint{{URL: URL}},
}

type (
	// PolygonIndicesResponse is the response from the Polygon Indices API.
	PolygonIndicesResponse struct { //nolint
		Status    string                 `json:"status"`
		Results   []PolygonIndicesResult `json:"results"`
		RequestID string                 `json:"request_id"`
	}

	// PolygonIndicesResult is a single index result in the response.
	PolygonIndicesResult struct { //nolint
		Ticker       string                    `json:"ticker"`
		Value        float64                   `json:"value"`
		LastUpdated  int64                     `json:"last_updated"`
		Timeframe    string                    `json:"timeframe"`
		Name         string                    `json:"name"`
		MarketStatus string                    `json:"market_status"`
		Type         string                    `json:"type"`
		Session      PolygonIndicesSessionData `json:"session"`
	}

	// PolygonIndicesSessionData represents the session data in the response.
	PolygonIndicesSessionData struct { //nolint
		Change        float64 `json:"change"`
		ChangePercent float64 `json:"change_percent"`
		Close         float64 `json:"close"`
		High          float64 `json:"high"`
		Low           float64 `json:"low"`
		Open          float64 `json:"open"`
		PreviousClose float64 `json:"previous_close"`
	}
)
