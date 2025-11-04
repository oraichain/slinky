package polygon

import (
	"github.com/dydxprotocol/slinky/oracle/config"
)

const (
	// Name is the name of the Polygon provider.
	PolygonIndicesName = "polygon_indices_ws"

	// WSS is the websocket URL for Polygon.
	RealTimeIndicesWSS = "wss://socket.massive.com/indices"

	// MarketTypeIndices is the market type for the Polygon indices.
	MarketTypeIndices = "indices"

	TopicPrefix = "V"
)

// DefaultWebSocketConfig is the default configuration for the OKX Websocket.
var DefaultPolygonIndicesWebSocketConfig = config.WebSocketConfig{
	Name:                          PolygonIndicesName,
	Enabled:                       true,
	MaxBufferSize:                 config.DefaultMaxBufferSize,
	ReconnectionTimeout:           config.DefaultReconnectionTimeout,
	PostConnectionTimeout:         config.DefaultPostConnectionTimeout,
	Endpoints:                     []config.Endpoint{{URL: RealTimeIndicesWSS}},
	ReadBufferSize:                config.DefaultReadBufferSize,
	WriteBufferSize:               config.DefaultWriteBufferSize,
	HandshakeTimeout:              config.DefaultHandshakeTimeout,
	EnableCompression:             config.DefaultEnableCompression,
	ReadTimeout:                   config.DefaultReadTimeout,
	WriteTimeout:                  config.DefaultWriteTimeout,
	PingInterval:                  config.DefaultPingInterval,
	WriteInterval:                 config.DefaultWriteInterval,
	MaxReadErrorCount:             config.DefaultMaxReadErrorCount,
	MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
	MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
	MarketType:                    MarketTypeIndices,
}
