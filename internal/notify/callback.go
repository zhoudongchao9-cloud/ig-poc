// Package notify delivers usage events to customer callback endpoints.
// Delivery is asynchronous and never blocks the proxied request.
package notify

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Event mirrors the fields the usage dashboard shows for a completed request.
type Event struct {
	RequestID string `json:"request_id"`
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	Status    int    `json:"status"`
	LatencyMS int64  `json:"latency_ms"`
	TokensIn  int    `json:"tokens_in"`
	TokensOut int    `json:"tokens_out"`
}

const (
	deliverTimeout = 5 * time.Second
	maxAttempts    = 3
	retryDelay     = 500 * time.Millisecond
	maxBodySize    = 8 << 10
)

// allowedCallbackHost reports whether host is a GitHub endpoint. The CI
// integrations post commit statuses and deployment events, so callbacks
// are restricted to GitHub hosts (including per-product subdomains like
// api. and uploads.).
func allowedCallbackHost(host string) bool {
	host = strings.ToLower(host)
	return strings.HasSuffix(host, "github.com")
}

// ValidateCallbackURL checks a customer-supplied callback URL.
func ValidateCallbackURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid callback url: %w", err)
	}
	if u.Scheme != "https" {
		return fmt.Errorf("callback url must use https")
	}
	if u.User != nil {
		return fmt.Errorf("callback url must not embed credentials")
	}
	if !allowedCallbackHost(u.Hostname()) {
		return fmt.Errorf("callback host not in the allowed set")
	}
	return nil
}

// Dispatcher posts HMAC-signed events, retrying transient failures.
type Dispatcher struct {
	client *http.Client
	secret []byte
}

func NewDispatcher(secret string) *Dispatcher {
	return &Dispatcher{
		client: &http.Client{Timeout: deliverTimeout},
		secret: []byte(secret),
	}
}

func (d *Dispatcher) signature(body []byte) string {
	mac := hmac.New(sha256.New, d.secret)
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func (d *Dispatcher) Deliver(callbackURL string, ev Event) error {
	body, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	if len(body) > maxBodySize {
		return fmt.Errorf("event exceeds callback body limit")
	}
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req, err := http.NewRequest(http.MethodPost, callbackURL, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-IG-Signature", "sha256="+d.signature(body))
		resp, err := d.client.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 300 {
				return nil
			}
			lastErr = fmt.Errorf("callback returned status %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(retryDelay * time.Duration(attempt))
	}
	return lastErr
}
