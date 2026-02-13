package domain_errors

type UserNotFoundError struct {
	UserID string
}

func (e UserNotFoundError) Error() string {
	return "User with ID " + e.UserID + " not found"
}
