package commands

import (
	"fmt"
	"strconv"
	"strings"
)

func normalizeSemver(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	return v
}

func compareSemver(a, b string) (int, error) {
	ap, err := parseSemver(a)
	if err != nil {
		return 0, fmt.Errorf("invalid current version %q: %w", a, err)
	}
	bp, err := parseSemver(b)
	if err != nil {
		return 0, fmt.Errorf("invalid latest version %q: %w", b, err)
	}
	for i := 0; i < 3; i++ {
		if ap[i] < bp[i] {
			return -1, nil
		}
		if ap[i] > bp[i] {
			return 1, nil
		}
	}
	return 0, nil
}

func parseSemver(v string) ([3]int, error) {
	var out [3]int
	v = strings.SplitN(v, "-", 2)[0]
	parts := strings.Split(v, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return out, fmt.Errorf("expected x.y or x.y.z")
	}
	for i := 0; i < 3; i++ {
		if i >= len(parts) {
			out[i] = 0
			continue
		}
		n, err := strconv.Atoi(parts[i])
		if err != nil {
			return out, err
		}
		out[i] = n
	}
	return out, nil
}
