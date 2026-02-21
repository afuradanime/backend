package domain_errors

type PostNotFoundError struct {
	PostID string
}

func (e PostNotFoundError) Error() string {
	return "Post " + e.PostID + " not found"
}

type PostDeletedError struct {
	PostID string
}

func (e PostDeletedError) Error() string {
	return "Post " + e.PostID + " is deleted"
}

type NotPostOwnerError struct {
	UserID string
	PostID string
}

func (e NotPostOwnerError) Error() string {
	return "User " + e.UserID + " is not the owner of post " + e.PostID
}
