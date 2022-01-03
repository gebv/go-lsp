package internal

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
)

func hashFromStrSlice(in []string) string {
	sort.Strings(in)
	hsha2 := sha1.Sum([]byte(strings.Join(in, ",")))
	if len(in) == 0 {
		return ""
	}
	return fmt.Sprintf("%x", hsha2)
}
