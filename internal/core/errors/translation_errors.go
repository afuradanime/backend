package domain_errors

type TranslationNotFoundError struct {
	AnimeID string
}

func (e TranslationNotFoundError) Error() string {
	return "Translation for anime " + e.AnimeID + " not found"
}

type TranslationNotPendingError struct{}

func (e TranslationNotPendingError) Error() string {
	return "Translation is not pending"
}

type AlreadyTranslatedError struct{}

func (e AlreadyTranslatedError) Error() string {
	return "This anime has already been translated"
}

type AlreadySubmittedTranslation struct{}

func (e AlreadySubmittedTranslation) Error() string {
	return "You've already submitted a translation for this anime"
}
