package domain_errors

type AnimeNotFoundError struct {
	AnimeID string
}

func (e AnimeNotFoundError) Error() string {
	return "Anime " + e.AnimeID + " not found"
}
