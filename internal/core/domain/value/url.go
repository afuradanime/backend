package value

import (
	"errors"
	"regexp"
	"strings"
)

// https://stackoverflow.com/questions/6168962/typical-url-lengths-for-storage-calculation-purposes-url-shortener
// 95% confidence
const MAX_URL_SIZE = 157

var URL_REGEX = regexp.MustCompile(`/((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+(:[0-9]+)?|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)/`)

type URL string

func NewURL(s string) (*URL, error) {

	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return nil, errors.New("cannot be empty")
	}

	if len(s) > MAX_URL_SIZE {
		return nil, errors.New("too long, max 157 characters")
	}

	// if !URL_REGEX.MatchString(s) {
	// 	return nil, errors.New("not a valid URL")
	// }

	url := URL(s)
	return &url, nil
}
