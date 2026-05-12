package commands

import "testing"

func TestCompareSemver(t *testing.T) {
	tests := []struct {
		a    string
		b    string
		want int
	}{
		{"0.1.0", "0.1.1", -1},
		{"0.1.1", "0.1.1", 0},
		{"1.0.0", "0.9.9", 1},
		{"0.2", "0.2.0", 0},
		{"0.2.1-beta", "0.2.1", 0},
	}
	for _, tc := range tests {
		got, err := compareSemver(tc.a, tc.b)
		if err != nil {
			t.Fatalf("compareSemver(%q, %q) unexpected error: %v", tc.a, tc.b, err)
		}
		if got != tc.want {
			t.Fatalf("compareSemver(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestParseUpdateArgs(t *testing.T) {
	opts, err := parseUpdateArgs([]string{"--check"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.checkOnly {
		t.Fatal("expected checkOnly=true")
	}
}

func TestParseUpdateArgsRejectsRepoFlag(t *testing.T) {
	if _, err := parseUpdateArgs([]string{"--repo", "owner/repo"}); err == nil {
		t.Fatal("expected error for unsupported --repo flag")
	}
}
