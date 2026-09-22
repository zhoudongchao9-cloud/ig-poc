package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"time"

	gin "github.com/gin-gonic/gin"
)

var usageSitePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`)

var usageClient = &http.Client{Timeout: 5 * time.Second}

// UsageBeacon posts a usage summary to the caller's endpoint after a
// request completes. Opt-in per request via the X-Usage-Endpoint header.
func UsageBeacon(c *gin.Context) {
	endpoint := c.GetHeader("X-Usage-Endpoint")
	if endpoint == "" {
		c.Next()
		return
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid endpoint"})
		return
	}
	if !usageSitePattern.MatchString(u.Hostname()) {
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid site"})
		return
	}
	start := time.Now()
	c.Next()
	if c.Writer.Status() >= 400 {
		return
	}
	body, _ := json.Marshal(map[string]any{
		"request_id": c.Writer.Header().Get("X-Request-ID"),
		"status":     c.Writer.Status(),
		"latency_ms": time.Since(start).Milliseconds(),
	})
	go func() {
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := usageClient.Do(req)
		if err == nil {
			resp.Body.Close()
		}
	}()
}
