package auth

import (
	"fmt"
	"strings"
)

// refreshTokenMathcer
type refreshTokenMathcer struct {
	prefix       string
	randomLength int
}

func (r refreshTokenMathcer) Matches(x interface{}) bool {
	key, ok := x.(string)
	if !ok {
		return false
	}
	if len(key) != r.randomLength+len(r.prefix) {
		return false
	}
	randomPart := key[len(r.prefix):]
	return len(randomPart) == r.randomLength && strings.HasPrefix(key, r.prefix)
}

func (r refreshTokenMathcer) String() string {
	return fmt.Sprintf("matches prefix %q and has length %d", r.prefix, r.randomLength)
}

// codeMatcher
type codeMatcher struct {
	length int
}

func (c codeMatcher) Matches(x interface{}) bool {
	key, ok := x.(string)
	if !ok {
		return false
	}
	return len(key) == c.length
}

func (c codeMatcher) String() string {
	return fmt.Sprintf("has length %d", c.length)
}
