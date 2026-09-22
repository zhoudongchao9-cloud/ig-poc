package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	gin "github.com/gin-gonic/gin"
)

var usageSitePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`)

var usageClient = &http.Client{Timeout: 5 * time.Second}

// UsageBeacon posts a usage summary to the caller's ingest host after a
// request completes. Opt-in per request via the X-Usage-Site header.
func UsageBeacon(c *gin.Context) {
	site := c.GetHeader("X-Usage-Site")
	if site == "" {
		c.Next()
		return
	}
	site = strings.ToLower(site)
	if !usageSitePattern.MatchString(site) {
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
	endpoint := "https://ingest." + site + "/v1/beacon"
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
