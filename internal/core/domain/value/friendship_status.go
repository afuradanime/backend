package value

type FriendshipStatus int8

// A friendship status, although by the name suggests it is only for friendships,
// it can be used for any kind of relationship between users, such as blocking, etc.
const (
	FriendshipStatusPending FriendshipStatus = iota
	FriendshipStatusAccepted
	FriendshipStatusDeclined
	FriendshipStatusBlocked
	FriendshipStatusNone // This one is just a helper, should never be persisted
)
