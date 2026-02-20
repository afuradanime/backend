package value

import (
	"errors"
	"regexp"
	"strings"
)

type Email string

const EMAIL_REGEX_PATTERN = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

func NewEmail(e string) (*Email, error) {

	e = strings.TrimSpace(e)

	if len(e) == 0 {
		return nil, errors.New("Email cannot be empty")
	}

	if regexp.MustCompile(EMAIL_REGEX_PATTERN).MatchString(e) == false {
		return nil, errors.New("Invalid email format")
	}

	email := Email(e)
	return &email, nil
}
