package domain

type AnimeType uint8

const (
	AnimeTypeTV AnimeType = iota + 1
	AnimeTypeOVA
	AnimeTypeMovie
	AnimeTypeSpecial
	AnimeTypeONA
	AnimeTypeMusic
	AnimeTypeUnknown
)

type AnimeStatus uint8

const (
	StatusFinishedAiring AnimeStatus = iota + 1
	StatusCurrentlyAiring
	StatusNotYetAired
	StatusUnknown
)

type SeasonType uint8

const (
	SeasonSpring SeasonType = iota
	SeasonSummer
	SeasonFall
	SeasonWinter
	SeasonUndefined
)

type Language uint8

const (
	LanguageEnglish Language = iota + 1
	LanguagePortuguese
)

type TagType uint8

const (
	TagGenre TagType = iota
	TagTheme
	TagDemographic
	TagExplicitGenre
)
