package value

type AnimeListItemStatus uint8

const (
	AnimeListItemStatusWatching AnimeListItemStatus = iota
	AnimeListItemStatusCompleted
	AnimeListItemStatusPaused
	AnimeListItemStatusDropped
	AnimeListItemStatusPlanning
)
