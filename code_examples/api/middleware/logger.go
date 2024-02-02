package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/samverrall/gopherjobs/api/contextutil"
)

func AttachLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		logger, ok := c.Locals(contextutil.LoggerCtx{}).(*slog.Logger)

		if !ok || logger == nil {
			logger = slog.Default()
		}

		logger = logger.With(
			"request_id", uuid.NewString(),
			"method", c.Method(),
			"remoteAddr", ctx.RemoteAddr().String(),
			"url", c.OriginalURL(),
		)

		c.Locals(contextutil.LoggerCtx{}, logger)
		return c.Next()
	}
}
