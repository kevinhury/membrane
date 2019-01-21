package urlutils

import "strings"

// MatchPath func
func MatchPath(path string, paths []string) (string, bool) {
	for _, p := range paths {
		c := strings.Split(p, "/")
		comps := strings.Split(path, "/")

		for i := 0; i < len(comps) && i < len(c); i++ {
			if c[i] == "" && i != 0 {
				return p, true
			}
			if c[i] != comps[i] {
				return "", false
			}
		}
	}

	return "", false
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
