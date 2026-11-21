package x402origin

import "testing"

// Cases are taken from the pre-fix behaviour recorded on 2026-11-19.
// Every expectation is the canonical (fixed) form, so this table fails
// when normalize points at normalizeBuggy.
func TestCanonicalOrigin(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{"default https port dropped", "https://api.example.com:443", "api.example.com"},
		{"host lowercased", "https://API.Example.com", "api.example.com"},
		{"default http port dropped", "http://api.example.com:80", "api.example.com"},
		{"non-default port kept", "https://api.example.com:8443", "api.example.com:8443"},
		{"ipv6 default port", "https://[::1]:443", "[::1]"},
		{"ipv6 non-default port", "https://[2001:db8::1]:8443", "[2001:db8::1]:8443"},
		{"no port, no change", "https://api.example.com", "api.example.com"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := normalize(c.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.want {
				t.Errorf("normalize(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestRejectsEmptyHost(t *testing.T) {
	if _, err := normalize("https:///path"); err == nil {
		t.Fatal("expected an error for an empty host")
	}
}
