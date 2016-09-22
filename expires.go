package expires

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type matchDef struct {
	Re       *regexp.Regexp
	Duration time.Duration
}

func init() {
	caddy.RegisterPlugin("expires", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	rules, err := parseRules(c)
	if err != nil {
		return err
	}

	cfg := httpserver.GetConfig(c)
	mid := func(next httpserver.Handler) httpserver.Handler {
		return expiresHandler{Next: next, Rules: rules}
	}
	cfg.AddMiddleware(mid)

	return nil
}

func parseRules(c *caddy.Controller) ([]*matchDef, error) {
	rules := []*matchDef{}

	for c.Next() {
		for c.NextBlock() {
			if c.Val() != "match" {
				return nil, c.SyntaxErr("match")
			}
			args := c.RemainingArgs()
			if len(args) != 2 {
				return nil, c.ArgErr()
			}
			re, err := regexp.Compile(args[0])
			if err != nil {
				return nil, err
			}
			duration := parseDuration(args[1])
			rule := &matchDef{Re: re, Duration: duration}

			rules = append(rules, rule)
		}
	}
	return rules, nil
}

func parseDuration(str string) time.Duration {
	durationRegex := regexp.MustCompile(`(?P<years>\d+y)?(?P<months>\d+m)?(?P<days>\d+d)?T?(?P<hours>\d+h)?(?P<minutes>\d+i)?(?P<seconds>\d+s)?`)
	matches := durationRegex.FindStringSubmatch(str)

	years := parseInt64(matches[1])
	months := parseInt64(matches[2])
	days := parseInt64(matches[3])
	hours := parseInt64(matches[4])
	minutes := parseInt64(matches[5])
	seconds := parseInt64(matches[6])

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	return time.Duration(years*24*365*hour + months*30*24*hour + days*24*hour + hours*hour + minutes*minute + seconds*second)
}

func parseInt64(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return 0
	}
	return int64(parsed)
}

type expiresHandler struct {
	Next  httpserver.Handler
	Rules []*matchDef
}

func (h expiresHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	for _, rule := range h.Rules {
		if rule.Re.MatchString(r.URL.Path) {
			w.Header().Set("Expires", time.Now().Add(rule.Duration).Format(time.RFC1123))
			break
		}
	}
	return h.Next.ServeHTTP(w, r)
}
