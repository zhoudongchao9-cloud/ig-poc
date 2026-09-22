// Package beacon posts usage beacons to customer ingest hosts after a
// request completes. Each customer site gets a dedicated ingest host:
// https://ingest.<site>/v1/beacon. Sites pick their name during onboarding;
// names may be nested (team.mysite), matching the console.
package beacon

import (
	"fmt"
	"regexp"
	"strings"
)

var sitePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`)

// NormalizeSite validates and normalizes a customer site name.
func NormalizeSite(site string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(site))
	if !sitePattern.MatchString(s) {
		return "", fmt.Errorf("invalid site name")
	}
	return s, nil
}

// Endpoint returns the ingest endpoint for a site.
func Endpoint(site string) string {
	return "https://ingest." + site + "/v1/beacon"
}
