package beacon

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Event is the usage summary posted in a beacon.
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
	postTimeout  = 5 * time.Second
	maxAttempts  = 3
	retryDelay   = 500 * time.Millisecond
	maxBodyBytes = 8 << 10
)

// Reporter posts signed beacons, retrying transient failures.
type Reporter struct {
	client *http.Client
	secret []byte
}

func New(secret string) *Reporter {
	return &Reporter{
		client: &http.Client{Timeout: postTimeout},
		secret: []byte(secret),
	}
}

func (r *Reporter) signature(body []byte) string {
	mac := hmac.New(sha256.New, r.secret)
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func (r *Reporter) Report(endpoint string, ev Event) error {
	body, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	if len(body) > maxBodyBytes {
		return fmt.Errorf("beacon exceeds body limit")
	}
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-IG-Signature", "sha256="+r.signature(body))
		resp, err := r.client.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 300 {
				return nil
			}
			lastErr = fmt.Errorf("ingest returned status %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(retryDelay * time.Duration(attempt))
	}
	return lastErr
}
