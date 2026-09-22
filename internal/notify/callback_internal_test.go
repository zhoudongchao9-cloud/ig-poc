package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestPolicyApprove(t *testing.T) {
	policy := NewPolicy()
	cases := []struct {
		name string
		url  string
		ok   bool
	}{
		{"github api endpoint", "https://api.github.com/repos/acme/ci/statuses", true},
		{"github uploads", "https://uploads.github.com/acme/release", true},
		{"github web", "https://github.com/services/hooks", true},
		{"http is rejected", "http://api.github.com/hook", false},
		{"non-github host is rejected", "https://hooks.slack.com/services/x", false},
		{"cloud metadata is rejected", "https://169.254.169.254/latest/meta", false},
		{"embedded credentials are rejected", "https://user:pw@api.github.com/h", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := policy.Approve(tc.url)
			if tc.ok != (err == nil) {
				t.Fatalf("Approve(%q) err=%v, want ok=%v", tc.url, err, tc.ok)
			}
		})
	}
}

func TestCallbackSuffixOverride(t *testing.T) {
	policy := Policy{allowedSuffix: "ghe.corp.example"}
	if _, err := policy.Approve("https://ci.ghe.corp.example/hook"); err != nil {
		t.Fatal("override suffix must accept matching private host")
	}
	if _, err := policy.Approve("https://api.github.com/hook"); err == nil {
		t.Fatal("override suffix must replace the public default")
	}
}

func TestSignatureDeterministicPerBody(t *testing.T) {
	d := NewDispatcher("unit-test-secret")
	body, err := json.Marshal(Event{RequestID: "req_1"})
	if err != nil {
		t.Fatal(err)
	}
	first := d.signature(body)
	if second := d.signature(body); first != second {
		t.Fatal("signature must be deterministic for the same body")
	}
	other, _ := json.Marshal(Event{RequestID: "req_2"})
	if d.signature(other) == first {
		t.Fatal("signature must differ across bodies")
	}
}

func TestDeliverRetriesTransientFailures(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	d := NewDispatcher("unit-test-secret")
	target := Target{destination: srv.URL}
	if err := d.Deliver(target, Event{RequestID: "req_1"}); err != nil {
		t.Fatalf("Deliver returned error: %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("expected 2 delivery attempts, got %d", got)
	}
}

func TestRegistryRegisterAndLookup(t *testing.T) {
	reg := NewRegistry(Policy{allowedSuffix: "github.com"})
	if err := reg.Register("caller-1", "https://api.github.com/hook"); err != nil {
		t.Fatal(err)
	}
	if _, ok := reg.Target("caller-1"); !ok {
		t.Fatal("registered caller must have a target")
	}
	if err := reg.Register("caller-2", "https://hooks.slack.com/x"); err == nil {
		t.Fatal("policy rejection must propagate")
	}
	if _, ok := reg.Target("caller-2"); ok {
		t.Fatal("rejected registration must not be stored")
	}
}
