package expires

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/mholt/caddy"
)

func TestParseBlock(t *testing.T) {
	conf := `expires {
		match .*\.jpg 1y
		match .*\.png 1m
		match .*\.bmp 1d
		match .*\.css 1h
		match .*\.js 1i
		match .*\.txt 1s
		match_header Content-Type .*json 1d
	}`

	c := caddy.NewTestController("http", conf)
	rules, err := parseRules(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(rules) != 7 {
		t.Fatalf("Wrong rules count %d expected 7", len(rules))
	}

	for i, rule := range rules {
		if rule.Duration() == 0 {
			t.Fatalf("Incorrect duration for %d", i)
		}
	}
}

func TestMatchPath(t *testing.T) {
	rule := matchDef{}
	rule.Parse([]string{".*\\.jpg", "1y"})

	header := http.Header{}
	header.Set("Content-Type", "image/jpeg")

	request := &http.Request{}
	request.URL, _ = url.Parse("http://www.example.com/image.jpg")

	if !rule.Match(header, request) {
		t.Fatalf("Expected %s to match %v", request.URL.Path, rule.re)
	}

	request.URL, _ = url.Parse("http://www.example.com/config.json")

	if rule.Match(header, request) {
		t.Fatalf("Expected %s to NOT match %v", request.URL.Path, rule.re)
	}
}

func TestMatchHeader(t *testing.T) {
	rule := headerMatchDef{}
	rule.Parse([]string{"Content-Type", ".+/.*json", "1y"})

	request := &http.Request{}
	request.URL, _ = url.Parse("http://www.example.com/config.json")

	header := http.Header{}

	if rule.Match(header, request) {
		t.Fatalf("Unset header should'nt match %v", rule.re)
	}

	header.Set("Content-Type", "application/json")

	if !rule.Match(header, request) {
		t.Fatalf("Expected %s to match %v", header.Get("Content-Type"), rule.re)
	}

	header.Set("Content-Type", "application/javascript")

	if rule.Match(header, request) {
		t.Fatalf("Expected %s to NOT match %v", header.Get("Content-Type"), rule.re)
	}
}

func TestParseDuration(t *testing.T) {
	duration := parseDuration("")
	if duration != 0 {
		t.Fatalf("Empty duration not 0: %d", duration)
	}
	duration = parseDuration("1y")
	if duration != 365*24*time.Hour {
		t.Fatalf("Wrong duration for year: %d", duration)
	}
	duration = parseDuration("1m")
	if duration != 30*24*time.Hour {
		t.Fatalf("Wrong duration for month: %d", duration)
	}
	duration = parseDuration("1d")
	if duration != 24*time.Hour {
		t.Fatalf("Wrong duration for day: %d", duration)
	}
	duration = parseDuration("1h")
	if duration != time.Hour {
		t.Fatalf("Wrong duration for hour: %d", duration)
	}
	duration = parseDuration("1i")
	if duration != time.Minute {
		t.Fatalf("Wrong duration for minute: %d", duration)
	}
	duration = parseDuration("1s")
	if duration != time.Second {
		t.Fatalf("Wrong duration for second: %d", duration)
	}
	duration = parseDuration("1y1m1d1h1i1s")
	if duration != ((396 * 24 * time.Hour) + time.Hour + time.Minute + time.Second) {
		t.Fatalf("Wrong duration for second: %d", duration)
	}
}
