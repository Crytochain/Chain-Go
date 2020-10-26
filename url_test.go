package accounts
import (
	"testing"
)
func TestURLParsing(t *testing.T) {
	url, err := parseURL("https:
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if url.Scheme != "https" {
		t.Errorf("expected: %v, got: %v", "https", url.Scheme)
	}
	if url.Path != "ethereum.org" {
		t.Errorf("expected: %v, got: %v", "ethereum.org", url.Path)
	}
	_, err = parseURL("ethereum.org")
	if err == nil {
		t.Error("expected err, got: nil")
	}
}
func TestURLString(t *testing.T) {
	url := URL{Scheme: "https", Path: "ethereum.org"}
	if url.String() != "https:
		t.Errorf("expected: %v, got: %v", "https:
	}
	url = URL{Scheme: "", Path: "ethereum.org"}
	if url.String() != "ethereum.org" {
		t.Errorf("expected: %v, got: %v", "ethereum.org", url.String())
	}
}
func TestURLMarshalJSON(t *testing.T) {
	url := URL{Scheme: "https", Path: "ethereum.org"}
	json, err := url.MarshalJSON()
	if err != nil {
		t.Errorf("unexpcted error: %v", err)
	}
	if string(json) != "\"https:
		t.Errorf("expected: %v, got: %v", "\"https:
	}
}
func TestURLUnmarshalJSON(t *testing.T) {
	url := &URL{}
	err := url.UnmarshalJSON([]byte("\"https:
	if err != nil {
		t.Errorf("unexpcted error: %v", err)
	}
	if url.Scheme != "https" {
		t.Errorf("expected: %v, got: %v", "https", url.Scheme)
	}
	if url.Path != "ethereum.org" {
		t.Errorf("expected: %v, got: %v", "https", url.Path)
	}
}
func TestURLComparison(t *testing.T) {
	tests := []struct {
		urlA   URL
		urlB   URL
		expect int
	}{
		{URL{"https", "ethereum.org"}, URL{"https", "ethereum.org"}, 0},
		{URL{"http", "ethereum.org"}, URL{"https", "ethereum.org"}, -1},
		{URL{"https", "ethereum.org/a"}, URL{"https", "ethereum.org"}, 1},
		{URL{"https", "abc.org"}, URL{"https", "ethereum.org"}, -1},
	}
	for i, tt := range tests {
		result := tt.urlA.Cmp(tt.urlB)
		if result != tt.expect {
			t.Errorf("test %d: cmp mismatch: expected: %d, got: %d", i, tt.expect, result)
		}
	}
}
