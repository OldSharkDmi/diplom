package ws

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"train-backend/internal/domain/repository"
)

// Handler — WS-endpoint: PUSH-уведомления о поездах.
//   - optional query-param uid=… — фильтр только для нужного поезда.
func Handler(repo repository.TrainRepository) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		uid := c.Query("uid")

		ctx := context.Background()  // базовый контекст
		lis, err := repo.Listen(ctx) // LISTEN train_updates
		if err != nil {
			_ = c.WriteMessage(websocket.TextMessage,
				[]byte("listen error: "+err.Error()))
			return
		}
		defer lis.Close()

		for {
			payload, err := lis.Wait(ctx)
			if err != nil { // клиент закрыл соединение или ctx отменён
				return
			}
			if uid != "" && payload != uid {
				continue
			}
			_ = c.WriteMessage(websocket.TextMessage, []byte(payload))
		}
	})
}
