package value

import (
	"errors"
	"strings"
)

// Max sized ddescription in the database is "Chika Gentou Gekiga: Shoujo Tsubaki" at
// 4982 characters in length, we'll add a cushion by upping that by 1000 characters
// For those who care, AVG(length(description)) is 401.864908788039
const MAX_SIZE = 6000

type LongStr string

func NewLongStr(s string) (*LongStr, error) {

	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return nil, errors.New("cannot be empty")
	}

	if len(s) > MAX_SIZE {
		return nil, errors.New("too long, max 6000 characters")
	}

	longStr := LongStr(s)
	return &longStr, nil
}
