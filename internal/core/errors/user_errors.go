package domain_errors

type UserNotFoundError struct {
	UserID string
}

func (e UserNotFoundError) Error() string {
	return "User with ID " + e.UserID + " not found"
}

type UserCantTranslate struct{}

func (e UserCantTranslate) Error() string {
	return "You are not allowed to translate"
}

type TooManySocials struct{}

func (e TooManySocials) Error() string {
	return "You can only have 5 social links"
}
