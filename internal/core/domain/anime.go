package domain

import "time"

// TODO: Deviamos por isto num helper.go? depois ve como queres fazer isso
func ParseISODate(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}

	if t, err := time.Parse(time.RFC3339, *s); err == nil {
		return &t, nil
	}

	// Oops
	// TODO: fazer alguma coisa ??
	t := time.Now()
	return &t, nil
}

type Anime struct {
	ID    uint32
	URL   string
	Title string

	Synonyms     []string
	Descriptions []Description

	Type     AnimeType
	Source   string
	Episodes uint32
	Status   AnimeStatus
	Airing   bool
	Duration string // "24 min per ep"

	StartDate *time.Time
	EndDate   *time.Time

	Season    Season
	Broadcast Broadcast

	ImageURL        string
	SmallImageURL   string
	LargeImageURL   string
	TrailerEmbedURL string

	Tags      []Tag
	Producers []Producer
	Licensors []Licensor
	Studios   []Studio
}

/**
Make a partial anime. Full anime fields are not filled in rn
*/
func NewAnime(
	id uint32,
	url string,
	title string,
	atype AnimeType,
	source string,
	episodes uint32,
	status AnimeStatus,
	airing bool,
	duration string,
	startDateISO string,
	endDateISO string,
	seasonType SeasonType,
	seasonYear uint16,
	broadcastDay string,
	broadcastTime string,
	broadcastTimezone string,
	imageURL string,
	smallImageURL string,
	largeImageURL string,
	trailerEmbedURL string,
) (*Anime, error) {

	var startDate *time.Time
	if startDateISO != "" {
		if t, err := time.Parse(time.RFC3339, startDateISO); err == nil {
			startDate = &t
		} else {
			t, err := time.Parse("2006-01-02", startDateISO)
			if err != nil {
				return nil, err
			}
			startDate = &t
		}
	}

	var endDate *time.Time
	if endDateISO != "" {
		if t, err := time.Parse(time.RFC3339, endDateISO); err == nil {
			endDate = &t
		} else {
			t, err := time.Parse("2006-01-02", endDateISO)
			if err != nil {
				return nil, err
			}
			endDate = &t
		}
	}

	anime := &Anime{
		ID:       id,
		URL:      url,
		Title:    title,
		Type:     atype,
		Source:   source,
		Episodes: episodes,
		Status:   status,
		Airing:   airing,
		Duration: duration,

		StartDate: startDate,
		EndDate:   endDate,

		Season: Season{
			Season: seasonType,
			Year:   seasonYear,
		},

		Broadcast: Broadcast{
			Day:      broadcastDay,
			Time:     broadcastTime,
			Timezone: broadcastTimezone,
		},

		ImageURL:        imageURL,
		SmallImageURL:   smallImageURL,
		LargeImageURL:   largeImageURL,
		TrailerEmbedURL: trailerEmbedURL,

		// Full anime only fields initialized but empty
		Synonyms:     []string{},
		Descriptions: []Description{},
		Tags:         []Tag{},
		Producers:    []Producer{},
		Licensors:    []Licensor{},
		Studios:      []Studio{},
	}

	return anime, nil
}
