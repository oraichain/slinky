package polygon

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dydxprotocol/slinky/oracle/config"
	"github.com/dydxprotocol/slinky/oracle/types"
	"github.com/dydxprotocol/slinky/providers/base/websocket/handlers"
	"github.com/massive-com/client-go/v2/websocket/models"
	"go.uber.org/zap"
)

var _ types.PriceWebSocketDataHandler = (*WebSocketHandler)(nil)

// WebSocketHandler implements the WebSocketDataHandler interface. This is used to handle
// messages received from the Binance websocket API.
type WebSocketHandler struct {
	logger *zap.Logger

	// ws is the config for the Binance websocket.
	ws config.WebSocketConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewWebSocketDataHandler returns a new Binance PriceWebSocketDataHandler.
func NewWebSocketDataHandler(
	logger *zap.Logger,
	ws config.WebSocketConfig,
) (types.PriceWebSocketDataHandler, error) {
	if ws.Name != PolygonIndicesName {
		return nil, fmt.Errorf("expected websocket config name %s, got %s", PolygonIndicesName, ws.Name)
	}

	if !ws.Enabled {
		return nil, fmt.Errorf("websocket config for %s is not enabled", PolygonIndicesName)
	}

	if err := ws.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid websocket config for %s: %w", PolygonIndicesName, err)
	}

	return &WebSocketHandler{
		logger: logger,
		ws:     ws,
		cache:  types.NewProviderTickers(),
	}, nil
}

// HandleMessage is used to handle a message received from the data provider. The Polygon websocket
// API is expected to handle the following types of messages:
//  1. SubscribeMessageResponse: This is a response to a subscription request. If the subscription
//     was successful, the response will contain a nil result. If the subscription failed, a
//     re-subscription message will be returned.
//  2. StreamMessageResponse: This is a response to a stream message. The stream message contains
//     the latest price of a ticker - either received when a trade is made or an automated price
//     update is received.
//
// Heartbeat messages are handled by default by the gorilla websocket library. The Polygon websocket
// API does not require any additional heartbeat messages to be sent. The pong frames are sent
// automatically by the gorilla websocket library.
func (h *WebSocketHandler) HandleMessage(
	message []byte,
) (types.PriceResponse, []handlers.WebsocketEncodedMessage, error) {
	var (
		resp             types.PriceResponse
		indicesResponses []IndicesResponse
	)

	var msgs []json.RawMessage
	if err := json.Unmarshal(message, &msgs); err != nil {
		return resp, nil, fmt.Errorf("failed to process raw messages: %w", err)
	}

	for _, msg := range msgs {
		var ev models.EventType
		err := json.Unmarshal(msg, &ev)
		if err != nil {
			return resp, nil, fmt.Errorf("failed to unmarshal event: %w", err)
		}
		switch ev.EventType {
		case "status":
			continue
		case "V": // this is indices response so this event is for index value
			var indicesResponse IndicesResponse
			if err := json.Unmarshal(msg, &indicesResponse); err != nil {
				return resp, nil, fmt.Errorf("failed to unmarshal index value: %w", err)
			}
			indicesResponses = append(indicesResponses, indicesResponse)
		}
	}

	resp, err := h.parseIndicesResponse(indicesResponses)

	return resp, nil, err
}

// CreateMessages is used to create a message to send to Binance. This is used to subscribe to
// the given tickers. This is called when the connection to the data provider is first established.
// Notably, the tickers have a unique identifier that is used to identify the messages going back
// and forth. This unique identifier is the same one sent in the initial subscription.
func (h *WebSocketHandler) CreateMessages(
	tickers []types.ProviderTicker,
) ([]handlers.WebsocketEncodedMessage, error) {
	var params []string
	for _, ticker := range tickers {
		params = append(params, TopicPrefix+"."+ticker.GetOffChainTicker())
		h.cache.Add(ticker)
	}

	msg, err := json.Marshal(&models.ControlMessage{
		Action: "subscribe",
		Params: strings.Join(params, ","),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to marshal control message: %w", err)
	}

	return []handlers.WebsocketEncodedMessage{msg}, nil
}

// HeartBeatMessages is not used for Binance. Heartbeats are handled on an ad-hoc basis when
// messages are received from the Binance websocket API.
//
// ref: https://developers.binance.com/docs/binance-spot-api-docs/web-socket-streams#aggregate-trade-streams
func (h *WebSocketHandler) HeartBeatMessages() ([]handlers.WebsocketEncodedMessage, error) {
	return nil, nil
}

// Copy is used to create a copy of the WebSocketHandler.
func (h *WebSocketHandler) Copy() types.PriceWebSocketDataHandler {
	return &WebSocketHandler{
		logger: h.logger,
		ws:     h.ws,
		cache:  types.NewProviderTickers(),
	}
}
