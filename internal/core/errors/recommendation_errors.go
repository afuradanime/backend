package domain_errors

type CannotRecommendYourselfError struct{}

func (e CannotRecommendYourselfError) Error() string {
	return "Cannot recommend an anime to yourself"
}

type RecommendationsDisabled struct{}

func (e RecommendationsDisabled) Error() string {
	return "This user does not allow recommendations"
}

type RecommendationStackFull struct{}

func (e RecommendationStackFull) Error() string {
	return "This user's recommendation stack is already full"
}

type AlreadyRecommended struct{}

func (e AlreadyRecommended) Error() string {
	return "This user has already been recommended this anime"
}
