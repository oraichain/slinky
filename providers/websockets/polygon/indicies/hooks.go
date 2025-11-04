package polygon

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dydxprotocol/slinky/oracle/config"
	wshandlers "github.com/dydxprotocol/slinky/providers/base/websocket/handlers"
	"github.com/massive-com/client-go/v2/websocket/models"
)

func PostDialHook(cfg config.ProviderConfig) wshandlers.PostDialHook {
	return func(handler *wshandlers.WebSocketConnHandlerImpl) error {
		handler.Conn.SetReadDeadline(time.Now().Add(cfg.WebSocket.ReadTimeout))
		// Authenticate the connection.
		msg, err := json.Marshal(&models.ControlMessage{
			Action: models.Auth,
			Params: cfg.WebSocket.Endpoints[0].Authentication.APIKey,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal auth message: %w", err)
		}
		handler.Write(msg)

		return nil
	}
}
