package dtos

type AnimeListItemDTO struct {
	// Anime Info
	AnimeID       uint32 `json:"animeId"`
	AnimeTitle    string `json:"animeTitle"`
	AnimeEpisodes uint32 `json:"animeEpisodes"`
	AnimeCoverURL string `json:"animeCoverUrl"`
	// Entry Info
	Status          uint8      `json:"status"`
	EpisodesWatched uint32     `json:"episodesWatched"`
	Rating          *RatingDTO `json:"rating,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	RewatchCount    uint8      `json:"rewatchCount"`
	CreatedAt       string     `json:"createdAt"`
	EditedAt        *string    `json:"editedAt,omitempty"`
}

type RatingDTO struct {
	Overall    uint8 `json:"overall"`
	Story      uint8 `json:"story"`
	Visuals    uint8 `json:"visuals"`
	Soundtrack uint8 `json:"soundtrack"`
	Enjoyment  uint8 `json:"enjoyment"`
}
