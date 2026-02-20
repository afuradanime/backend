package utils

import "time"

func Clamp(a, bottom, top int) int {

	if a < bottom {
		a = bottom
	}

	if a > top {
		a = top
	}

	return a
}

func ClampBottom(a, bottom int) int {

	if a < bottom {
		a = bottom
	}

	return a
}

func ClampTop(a, top int) int {

	if a > top {
		a = top
	}

	return a
}

func ParseISODate(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}

	if t, err := time.Parse(time.RFC3339, *s); err == nil {
		return &t
	}

	t := time.Now()
	return &t
}
