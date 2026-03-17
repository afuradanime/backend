package value

type PostParentType uint8

const (
	ParentTypeUser   PostParentType = iota // profile
	ParentTypeThread                       // anime specific posts
	ParentTypeGroup                        // group specific posts
	ParentTypePost                         // another post (nested replies)
)
