package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

func Logging(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Логируем начало обработки запроса.
		logger.Info("Incoming request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
		)

		// Можно сохранить logger в контексте, если потребуется в обработчике.
		ctx := context.WithValue(r.Context(), "logger", logger)
		next.ServeHTTP(w, r.WithContext(ctx))

		// После обработки запроса логируем время выполнения.
		duration := time.Since(start)
		logger.Info("Request processed",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
			slog.Duration("duration", duration),
		)
	})
}
