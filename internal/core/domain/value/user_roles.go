package value

type UserRole uint8

const (
	UserRoleAdmin UserRole = iota + 1
	UserRoleModerator
	UserRoleUser
)
