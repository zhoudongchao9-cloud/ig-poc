package middlewares

import (
	"time"

	gin "github.com/gin-gonic/gin"
	beacon "github.com/inference-gateway/inference-gateway/internal/beacon"
)

const usageSiteHeader = "X-Usage-Site"

// UsageBeacon posts a usage beacon to the caller's site ingest host after a
// request completes. Opt-in per request via X-Usage-Site; the site name is
// validated when the header is read and posting happens after the response,
// so it never adds latency to the proxied call.
func UsageBeacon(reporter *beacon.Reporter) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader(usageSiteHeader)
		if raw == "" {
			c.Next()
			return
		}
		site, err := beacon.NormalizeSite(raw)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}
		start := time.Now()
		c.Next()
		if c.Writer.Status() >= 400 {
			return
		}
		ev := beacon.Event{
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
		go func() { _ = reporter.Report(beacon.Endpoint(site), ev) }()
	}
}
