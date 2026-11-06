package launcher

import (
	"strings"
	"testing"
)

func TestLevenshteinDistance(t *testing.T) {
	if dist := LevenshteinDistance("nginx", "ngnx"); dist != 1 {
		t.Errorf("Expected distance 1, got %d", dist)
	}
	if dist := LevenshteinDistance("clean", "clean"); dist != 0 {
		t.Errorf("Expected distance 0, got %d", dist)
	}
	if dist := LevenshteinDistance("service", "srvice"); dist != 1 {
		t.Errorf("Expected distance 1, got %d", dist)
	}
}

func TestSpotlightFilteringAndGhostText(t *testing.T) {
	st := NewSpotlightState()

	// Initial state with empty query
	if len(st.Results) == 0 {
		t.Fatalf("Expected default items, got empty")
	}

	// Filter with prefix "clea"
	st.Query = "clea"
	st.Filter()

	if len(st.Results) == 0 {
		t.Fatalf("Expected match for 'clea', got none")
	}

	if st.GhostText == "" {
		t.Errorf("Expected non-empty ghost text for 'clea'")
	}

	// Test AcceptGhostText
	st.AcceptGhostText()
	if !strings.HasPrefix(strings.ToLower(st.Query), "clean") {
		t.Errorf("Expected query to start with 'clean', got %q", st.Query)
	}
}

func TestTypoToleranceInSpotlight(t *testing.T) {
	st := NewSpotlightState()

	// Typo "ngnx" for nginx
	st.Query = "ngnx"
	st.Filter()

	foundNginx := false
	for _, res := range st.Results {
		if strings.Contains(strings.ToLower(res.Title), "nginx") {
			foundNginx = true
			break
		}
	}

	if !foundNginx {
		t.Errorf("Expected typo-tolerant match for 'ngnx' to find nginx, got %d results", len(st.Results))
	}
}
