package util

import "fmt"

func PadStringTo(v string, n int) string {
	if len(v) >= n {
		return v[:n]
	}
	return fmt.Sprintf("%-*s", n, v)
}
