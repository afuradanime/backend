package domain_errors

type UnauthorizedError struct {
}

func (e UnauthorizedError) Error() string {
	return "You cannot do that shit bruh"
}
