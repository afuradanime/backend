package value

import (
	"errors"
	"strings"
)

type TinyStr string

func NewTinyStr(s string) (*TinyStr, error) {

	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return nil, errors.New("cannot be empty")
	}

	if len(s) > 32 {
		return nil, errors.New("too long, max 32 characters")
	}

	tinyStr := TinyStr(s)
	return &tinyStr, nil
}
