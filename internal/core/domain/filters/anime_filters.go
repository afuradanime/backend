package filters

type AnimeFilter struct {
	Name        *string
	Type        *uint32
	Status      *uint32
	StartDate   *int64
	EndDate     *int64
	MinEpisodes *uint32
	MaxEpisodes *uint32
}
