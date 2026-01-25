package value

import (
	"errors"
	"regexp"
)

type Email string

const emailRegexPattern = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

func NewEmail(e string) (*Email, error) {
	if len(e) == 0 {
		return nil, errors.New("Email cannot be empty")
	}

	if regexp.MustCompile(emailRegexPattern).MatchString(e) == false {
		return nil, errors.New("Invalid email format")
	}

	email := Email(e)
	return &email, nil
}
