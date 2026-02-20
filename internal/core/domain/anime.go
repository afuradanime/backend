package domain

import (
	"time"

	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/utils"
)

type Anime struct {
	ID    uint32
	URL   string
	Title string

	Synonyms     []string
	Descriptions []value.Description

	Type     value.AnimeType
	Source   string
	Episodes uint32
	Status   value.AnimeStatus
	Airing   bool
	Duration string // "24 min per ep"

	StartDate *time.Time
	EndDate   *time.Time

	Season    value.Season
	Broadcast value.Broadcast

	ImageURL        string
	SmallImageURL   string
	LargeImageURL   string
	TrailerEmbedURL string

	Tags      []value.Tag
	Producers []value.Producer
	Licensors []value.Licensor
	Studios   []value.Studio
}

/*
*
Make a partial anime. Full anime fields are not filled in rn
*/
func NewAnime(
	id uint32,
	url string,
	title string,
	atype value.AnimeType,
	source string,
	episodes uint32,
	status value.AnimeStatus,
	airing bool,
	duration string,
	startDateISO string,
	endDateISO string,
	seasonType value.SeasonType,
	seasonYear uint16,
	broadcastDay string,
	broadcastTime string,
	broadcastTimezone string,
	imageURL string,
	smallImageURL string,
	largeImageURL string,
	trailerEmbedURL string,
) (*Anime, error) {

	var startDate *time.Time = utils.ParseISODate(&startDateISO)
	var endDate *time.Time = utils.ParseISODate(&endDateISO)

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

		Season: value.Season{
			Season: seasonType,
			Year:   seasonYear,
		},

		Broadcast: value.Broadcast{
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
		Descriptions: []value.Description{},
		Tags:         []value.Tag{},
		Producers:    []value.Producer{},
		Licensors:    []value.Licensor{},
		Studios:      []value.Studio{},
	}

	return anime, nil
}

// Builder methods for the full anime fields.
// We can use these to fill in the full anime fields after we create the partial anime with the NewAnime constructor
// Perhaps we should have methods to add as list
func (anime *Anime) AddDescription(desc value.Description) {
	anime.Descriptions = append(anime.Descriptions, desc)
}

func (anime *Anime) AddSynonym(synonym string) {
	anime.Synonyms = append(anime.Synonyms, synonym)
}

func (anime *Anime) AddTag(tag value.Tag) {
	anime.Tags = append(anime.Tags, tag)
}

func (anime *Anime) AddProducer(producer value.Producer) {
	anime.Producers = append(anime.Producers, producer)
}

func (anime *Anime) AddLicensor(licensor value.Licensor) {
	anime.Licensors = append(anime.Licensors, licensor)
}

func (anime *Anime) AddStudio(studio value.Studio) {
	anime.Studios = append(anime.Studios, studio)
}
