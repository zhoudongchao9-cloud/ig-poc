package middlewares

import (
	"time"

	gin "github.com/gin-gonic/gin"
	notify "github.com/inference-gateway/inference-gateway/internal/notify"
)

const callbackURLHeader = "X-Callback-URL"

// Callback lets customers receive a usage event when their request
// completes. The X-Callback-URL header is approved by policy as soon as
// it is read; delivery happens after the response using the approved
// target, so it never adds latency to the proxied call.
func Callback(dispatcher *notify.Dispatcher, policy notify.Policy) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader(callbackURLHeader)
		if raw == "" {
			c.Next()
			return
		}
		target, err := policy.Approve(raw)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}
		start := time.Now()
		c.Next()
		if c.Writer.Status() >= 400 {
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
