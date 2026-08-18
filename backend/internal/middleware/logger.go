package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

func Logger() func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		start := time.Now()

		// 1. Process the request down the chain
		next(ctx)

		// 2. Post-execution metrics extraction
		duration := time.Since(start)
		status := ctx.Status()
		if status == 0 {
			status = http.StatusOK
		}

		opID := ""
		if op := ctx.Operation(); op != nil {
			opID = op.OperationID
		}

		clientIP := getClientIP(ctx)

		// 3. Construct structured attributes
		attrs := []slog.Attr{
			slog.String("method", ctx.Method()),
			slog.String("path", ctx.URL().Path),
			slog.Int("status", status),
			slog.Int64("duration_ms", duration.Milliseconds()),
			slog.String("operation_id", opID),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", ctx.Header("User-Agent")),
		}

		// 4. Determine log level based on HTTP status code
		level := slog.LevelInfo
		switch {
		case status >= http.StatusInternalServerError:
			level = slog.LevelError
		case status >= http.StatusBadRequest:
			level = slog.LevelWarn
		}

		slog.LogAttrs(ctx.Context(), level, "http_request", attrs...)
	}
}

func getClientIP(ctx huma.Context) string {
	if xff := ctx.Header("X-Forwarded-For"); xff != "" {
		// CloudFront/ALB appends client IP as the first comma-separated item
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xrip := ctx.Header("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}
	return ""
}
