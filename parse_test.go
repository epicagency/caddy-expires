package expires

import (
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
	}`

	c := caddy.NewTestController("http", conf)
	rules, err := parseRules(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(rules) != 6 {
		t.Fatalf("Wrong rules count %d expected 6", len(rules))
	}

	for i, rule := range rules {
		if rule.Re == nil {
			t.Fatalf("Empty rule %d", i)
		}
		if rule.Duration == 0 {
			t.Fatalf("Incorrect duration for %d: %s", i, rule.Re.String())
		}
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
