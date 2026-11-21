package x402origin

import (
	"fmt"
	"net/url"
	"strings"
)

// normalizeBuggy reproduces the pre-fix behaviour verbatim:
// validate.go returned parsed.Host, so the default port, the case of the
// host and the absence of an explicit canonical form all leaked through.
func normalizeBuggy(raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("empty host in %q", raw)
	}
	return parsed.Host, nil
}

// normalizeFixed is the post-fix behaviour: host lowercased, default port
// dropped, non-default port kept, IPv6 re-bracketed.
func normalizeFixed(raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("empty host in %q", raw)
	}
	host := strings.ToLower(parsed.Hostname())
	port := parsed.Port()
	scheme := strings.ToLower(parsed.Scheme)
	if (scheme == "http" && port == "80") || (scheme == "https" && port == "443") {
		port = ""
	}
	if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	if port != "" {
		host = host + ":" + port
	}
	return host, nil
}

// normalize is the seam that must_fail.sh flips to normalizeBuggy.
var normalize = normalizeFixed
