package ws

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/websocket/v2"
	"train-backend/internal/domain/repository"
)

func Handler(repo repository.TrainRepository) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		ctx := context.Background()
		lis, err := repo.Listen(ctx)
		if err != nil {
			_ = c.Close()
			return
		}
		defer lis.Close()

		for {
			payload, err := lis.WaitForNotification(ctx)
			if err != nil {
				return
			}
			_ = c.WriteMessage(websocket.TextMessage, []byte(payload))
		}
	})
}
