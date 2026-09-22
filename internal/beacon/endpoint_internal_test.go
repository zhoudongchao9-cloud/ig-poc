package beacon

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestNormalizeSite(t *testing.T) {
	cases := []struct {
		in  string
		ok  bool
	}{
		{"acme", true},
		{"team.mysite", true},
		{"site-42", true},
		{"", false},
		{"UPPER", false},
		{"_underscore", false},
		{"has space", false},
		{"-leading", false},
		{"trailing-", false},
	}
	for _, tc := range cases {
		if _, err := NormalizeSite(tc.in); (err == nil) != tc.ok {
			t.Fatalf("NormalizeSite(%q) err=%v, want ok=%v", tc.in, err, tc.ok)
		}
	}
}

func TestEndpoint(t *testing.T) {
	if got := Endpoint("team.mysite"); got != "https://ingest.team.mysite/v1/beacon" {
		t.Fatalf("unexpected endpoint: %s", got)
	}
}

func TestReportRetriesTransientFailures(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	rep := New("unit-test-secret")
	if err := rep.Report(srv.URL, Event{RequestID: "req_1"}); err != nil {
		t.Fatalf("Report returned error: %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("expected 2 attempts, got %d", got)
	}
}
