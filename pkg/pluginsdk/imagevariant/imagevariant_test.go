package imagevariant

import "testing"

func TestCanonicalValues(t *testing.T) {
	for _, tc := range []struct {
		name string
		got  string
		want string
	}{
		{"Card", Card, "card"},
		{"Featured", Featured, "featured"},
		{"Large", Large, "large"},
		{"Full", Full, "full"},
		{"Original", Original, "original"},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}
