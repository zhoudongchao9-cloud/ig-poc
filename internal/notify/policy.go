package notify

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Policy decides which callback targets this deployment may deliver to.
// The allowed domain suffix comes from IG_CALLBACK_SUFFIX so operators
// can point callbacks at a private GitHub Enterprise deployment.
type Policy struct {
	allowedSuffix string
}

// NewPolicy reads the configured callback domain suffix.
func NewPolicy() Policy {
	return Policy{allowedSuffix: configuredCallbackHostSuffix()}
}

// Target is an approved callback destination. It only exists for URLs
// that passed policy checks, so handlers can hand it straight to the
// dispatcher without re-validating.
type Target struct {
	destination string
}

// Approve checks a raw callback URL against policy and returns the
// delivery target.
func (p Policy) Approve(raw string) (Target, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return Target{}, fmt.Errorf("invalid callback url: %w", err)
	}
	if u.Scheme != "https" {
		return Target{}, fmt.Errorf("callback url must use https")
	}
	if u.User != nil {
		return Target{}, fmt.Errorf("callback url must not embed credentials")
	}
	if !strings.HasSuffix(strings.ToLower(u.Hostname()), p.allowedSuffix) {
		return Target{}, fmt.Errorf("callback host not in the allowed set")
	}
	return Target{destination: raw}, nil
}

// Deliverable is the URL the dispatcher posts events to.
func (t Target) Deliverable() string {
	return t.destination
}

// configuredCallbackHostSuffix resolves the configured callback domain
// suffix; the public default is github.com.
func configuredCallbackHostSuffix() string {
	if s := strings.TrimSpace(os.Getenv("IG_CALLBACK_SUFFIX")); s != "" {
		return strings.ToLower(s)
	}
	return "github.com"
}
