package utils

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
	"unicode"
)

const joinCodeChars = "ABCDEFGHJKLMNPQRSTUVWXYZ1234567890"

func GenerateJoinCode(title string) string {
	prefix := ""
	for _, r := range title {
		if unicode.IsLetter(r) {
			prefix += strings.ToUpper(string(r))
			if len(prefix) == 4 {
				break
			}
		}
	}
	for len(prefix) < 4 {
		prefix += "X"
	}
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	suffix := make([]byte, 4)
	for i, v := range b {
		suffix[i] = joinCodeChars[int(v)%len(joinCodeChars)]
	}
	return fmt.Sprintf("%s-%d-%s", prefix, time.Now().Year(), string(suffix))
}
