package api

import (
	"time"

	gin "github.com/gin-gonic/gin"
	notify "github.com/inference-gateway/inference-gateway/internal/notify"
)

// CallerKey identifies the registering party. The gateway fronts services
// with their own API keys, so the Authorization header is the caller
// identity; requests without one fall back to the client IP.
func CallerKey(c *gin.Context) string {
	if key := c.Request.Header.Get("Authorization"); key != "" {
		return key
	}
	return c.ClientIP()
}

type CallbackRegistration struct {
	URL string `json:"url"`
}

// RegisterCallbackHandler stores a caller's webhook target after policy
// approval. Callers register once; their subsequent requests get usage
// events delivered to the stored target.
func RegisterCallbackHandler(registry *notify.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CallbackRegistration
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, ErrorResponse{Error: "invalid registration"})
			return
		}
		if err := registry.Register(CallerKey(c), req.URL); err != nil {
			c.JSON(400, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(200, ResponseJSON{Message: "callback registered"})
	}
}

// CallbackDelivery delivers usage events to callers that registered a
// webhook. The target was approved at registration time; delivery only
// resolves the stored target, so it never adds latency to the proxied
// call.
func CallbackDelivery(registry *notify.Registry, dispatcher *notify.Dispatcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		if c.Writer.Status() >= 400 {
			return
		}
		target, ok := registry.Target(CallerKey(c))
		if !ok {
			return
		}
		ev := notify.Event{
			RequestID: c.Writer.Header().Get("X-Request-ID"),
			Status:    c.Writer.Status(),
			LatencyMS: time.Since(start).Milliseconds(),
		}
		if v, ok := c.Get("provider"); ok {
			ev.Provider, _ = v.(string)
		}
		if v, ok := c.Get("model"); ok {
			ev.Model, _ = v.(string)
		}
		// TODO: wire token counts from the usage middleware.
		go func() { _ = dispatcher.Deliver(target, ev) }()
	}
}
