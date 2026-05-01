package value

type ActivityStatus uint8

const (
	Offline ActivityStatus = iota + 1
	Online
	Idle
)
