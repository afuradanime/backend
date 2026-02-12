package domain

import "time"

// TODO: Deviamos por isto num helper.go? depois ve como queres fazer isso
func ParseISODate(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}

	if t, err := time.Parse(time.RFC3339, *s); err == nil {
		return &t
	}

	// Oops
	// TODO: fazer alguma coisa ??
	t := time.Now()
	return &t
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

	var startDate *time.Time = ParseISODate(&startDateISO)
	var endDate *time.Time = ParseISODate(&endDateISO)

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

// Builder methods for the full anime fields.
// We can use these to fill in the full anime fields after we create the partial anime with the NewAnime constructor
// Perhaps we should have methods to add as list
func (anime *Anime) AddDescription(desc Description) {
	anime.Descriptions = append(anime.Descriptions, desc)
}

func (anime *Anime) AddSynonym(synonym string) {
	anime.Synonyms = append(anime.Synonyms, synonym)
}

func (anime *Anime) AddTag(tag Tag) {
	anime.Tags = append(anime.Tags, tag)
}

func (anime *Anime) AddProducer(producer Producer) {
	anime.Producers = append(anime.Producers, producer)
}

func (anime *Anime) AddLicensor(licensor Licensor) {
	anime.Licensors = append(anime.Licensors, licensor)
}

func (anime *Anime) AddStudio(studio Studio) {
	anime.Studios = append(anime.Studios, studio)
}
