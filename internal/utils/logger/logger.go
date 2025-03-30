package logger

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

type responseCapture struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (r *responseCapture) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (l *Logger) LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		req := c.Request()
		res := c.Response()

		requestBody := readRequestBody(c)

		rec := &responseCapture{
			ResponseWriter: res.Writer,
			body:           new(bytes.Buffer),
		}
		res.Writer = rec

		err := next(c)

		duration := time.Since(start)

		responseBody := rec.body.String()

		slog.Info("HTTP Request",
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
			slog.Any("user_id", c.Get("user_id")),
			slog.String("request_body", requestBody),
			slog.Int("status", res.Status),
			slog.String("response_body", responseBody),
			slog.Duration("duration", duration),
			slog.Any("error", err),
		)

		return err
	}
}

func readRequestBody(c echo.Context) string {
	req := c.Request()

	if req.Body == nil {
		return ""
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return ""
	}

	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes)
}
