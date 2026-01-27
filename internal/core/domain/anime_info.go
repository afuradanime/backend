package domain

type Season struct {
	Season SeasonType
	Year   uint16
}

type Broadcast struct {
	Day      string // "Saturdays"
	Time     string // "01:00"
	Timezone string // "Asia/Tokyo"
}

type Description struct {
	Language    Language
	Description string
}

type Producer struct {
	ID   uint32
	Name string
	Type string
	URL  string
}

type Licensor struct {
	ID   uint32
	Name string
	Type string
	URL  string
}

type Studio struct {
	ID   uint32
	Name string
	URL  string
}

type Tag struct {
	ID   uint32
	Name string
	Type TagType
	URL  string
}
