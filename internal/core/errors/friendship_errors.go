package domain_errors

type CannotFriendYourselfError struct{}

func (e CannotFriendYourselfError) Error() string {
	return "You cannot friend yourself"
}

type NotFriendsError struct {
	Initiator string
	Receiver  string
}

func (e NotFriendsError) Error() string {
	return "Users " + e.Initiator + " and " + e.Receiver + " are not friends"
}

type UserBlockedError struct {
	Initiator string
	Receiver  string
}

func (e UserBlockedError) Error() string {
	return "User " + e.Receiver + " has blocked user " + e.Initiator
}

type AlreadyFriendsError struct {
	Initiator string
	Receiver  string
}

func (e AlreadyFriendsError) Error() string {
	return "Users " + e.Initiator + " and " + e.Receiver + " are already friends"
}

type FriendRequestAlreadySentError struct {
	Initiator string
	Receiver  string
}

func (e FriendRequestAlreadySentError) Error() string {
	return "User " + e.Initiator + " has already sent a friend request to user " + e.Receiver
}

type CantOperateOnNonPendingRequestError struct {
}

func (e CantOperateOnNonPendingRequestError) Error() string {
	return "You can only accept or decline pending friend requests"
}

type CannotBlockYourselfError struct{}

func (e CannotBlockYourselfError) Error() string {
	return "You cannot block yourself"
}

type CannotAcceptAlienRequest struct{}

func (e CannotAcceptAlienRequest) Error() string {
	return "You cannot act on this request"
}
