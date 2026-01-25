package domain

type Anime struct {
	ID            uint32
	Sources       string
	Title         string
	Type          uint32
	Episodes      uint32
	Status        uint32
	Picture       string
	Thumbnail     string
	DurationValue *float32
}

func NewAnime(
	id uint32,
	sources string,
	title string,
	atype uint32,
	episodes uint32,
	status uint32,
	picture string,
	thumbnail string,
	durationValue *float32,
) *Anime {
	return &Anime{
		ID:            id,
		Sources:       sources,
		Title:         title,
		Type:          atype,
		Episodes:      episodes,
		Status:        status,
		Picture:       picture,
		Thumbnail:     thumbnail,
		DurationValue: durationValue,
	}
}
