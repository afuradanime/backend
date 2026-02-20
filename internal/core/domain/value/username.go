package value

import (
	"errors"
	"regexp"
	"strings"
)

type Username string

var USERNAME_REGEX_PATTERN = regexp.MustCompile(`^[\p{L}0-9 _\-]+$`)

func NewUsername(s string) (*Username, error) {

	s = strings.TrimSpace(s)

	if len([]int32(s)) < 2 {
		return nil, errors.New("too short, min 2 characters")
	}

	if len([]int32(s)) > 64 {
		return nil, errors.New("too long, max 120 characters")
	}

	if !USERNAME_REGEX_PATTERN.MatchString(s) {
		return nil, errors.New("the username can only contain letters from the Portuguese alphabet, numbers, spaces, underscores and hyphens")
	}

	tinyStr := Username(s)
	return &tinyStr, nil
}
