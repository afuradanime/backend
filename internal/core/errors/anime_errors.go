package domain_errors

type AnimeNotFoundError struct {
	AnimeID string
}

func (e AnimeNotFoundError) Error() string {
	return "Anime " + e.AnimeID + " not found"
}

type StudioNotFoundError struct {
	StudioID string
}

func (e StudioNotFoundError) Error() string {
	return "Studio " + e.StudioID + " not found"
}

type ProducerNotFoundError struct {
	ProducerID string
}

func (e ProducerNotFoundError) Error() string {
	return "Producer " + e.ProducerID + " not found"
}

type LicensorNotFoundError struct {
	LicensorID string
}

func (e LicensorNotFoundError) Error() string {
	return "Licensor " + e.LicensorID + " not found"
}

type AnimeFetchFailedError struct{}

func (e AnimeFetchFailedError) Error() string {
	return "Could not process query"
}
