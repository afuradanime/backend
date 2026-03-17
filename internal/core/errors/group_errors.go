package domain_errors

type NoModeratorsLeftError struct{}

func (e NoModeratorsLeftError) Error() string {
	return "Can't remove the last moderator"
}

type NotModeratingError struct{}

func (e NotModeratingError) Error() string {
	return "This user is not moderating this group"
}

type AlreadyModeratingError struct{}

func (e AlreadyModeratingError) Error() string {
	return "This user is already moderating this group"
}

type GroupNotFoundError struct {
	GroupID string
}

func (e GroupNotFoundError) Error() string {
	return "User with ID " + e.GroupID + " not found"
}
