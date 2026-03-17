package domain

import (
	"github.com/afuradanime/backend/internal/core/utils"
)

const RECENT_EVALUATION_RING_SIZE = 5

type UserRating struct {
	User   int
	Rating Rating
}

/*
* The rating cache is a denormalized structure that stores the user's ratings for an anime,
* allowing for quick retrieval without needing to calculate the overall rating from individual components every time.
* It is updated whenever the user changes their ratings for story, visuals, or soundtrack,
* and it provides an overall rating that can be used for sorting and displaying in the UI.

* This is shadow data, it is managed automatically in the background, all external perations should be
* fetching the values only
 */

type RatingCache struct {
	AnimeID          int                           `json:"animeId" bson:"anime_id"`
	TotalOverall     float32                       `json:"overall" bson:"overall"`
	TotalStory       uint32                        `json:"story" bson:"story"`
	TotalVisuals     uint32                        `json:"visuals" bson:"visuals"`
	TotalSoundtrack  uint32                        `json:"soundtrack" bson:"soundtrack"`
	RecentEvaluation *utils.RingBuffer[UserRating] `json:"recentEvals" bson:"recent_evals"`
	UserCounter      int                           `json:"user_counter" bson:"user_counter"`
}

func NewRatingCache(animeID int) *RatingCache {
	return &RatingCache{
		AnimeID:          animeID,
		RecentEvaluation: utils.NewRingBuffer[UserRating](RECENT_EVALUATION_RING_SIZE),
	}
}

// Aux method to calculate average
func (r *RatingCache) CalculateAverage() {
	r.TotalOverall = float32(r.TotalStory+r.TotalVisuals+r.TotalSoundtrack) / float32(3*r.UserCounter)
}

func (r *RatingCache) UpdateCache(userId uint32, story, visuals, soundtrack uint32) {
	r.TotalStory += story
	r.TotalVisuals += visuals
	r.TotalSoundtrack += soundtrack
	r.UserCounter++

	// Recalculate overall
	r.TotalOverall = float32(r.TotalStory+r.TotalVisuals+r.TotalSoundtrack) / float32(3*r.UserCounter)

	// Add to snapshot
	r.addToRecent(userId, story, visuals, soundtrack)
}

func (r *RatingCache) RemoveRating(userId, story, visuals, soundtrack uint32) {
	// Update totals
	r.TotalStory -= story
	r.TotalVisuals -= visuals
	r.TotalSoundtrack -= soundtrack

	// Decrement user counter
	if r.UserCounter > 0 {
		r.UserCounter--
	}

	// Recalculate overall rating if there are still ratings left
	if r.UserCounter > 0 {
		r.CalculateAverage()
	} else {
		r.TotalOverall = 0
	}

	// Note: We don't remove from the ring buffer
}

func (r *RatingCache) UpdateExistingRating(userId, oldStory, oldVisuals, oldSoundtrack, newStory, newVisuals, newSoundtrack uint32) {
	// Update totals by removing old ratings and adding new ratings
	r.TotalStory = r.TotalStory - oldStory + newStory
	r.TotalVisuals = r.TotalVisuals - oldVisuals + newVisuals
	r.TotalSoundtrack = r.TotalSoundtrack - oldSoundtrack + newSoundtrack

	// Recalculate overall rating
	if r.UserCounter > 0 {
		r.CalculateAverage()
	} else {
		r.TotalOverall = 0
	}

	// We treat an update as a new rating to avoid getting phylosophical about
	// the choice of a ring buffer
	// We treat an update as a "New Recent Event"
	r.addToRecent(userId, newStory, newVisuals, newSoundtrack)
}

func (r *RatingCache) GetOverallRating() float32 {
	return r.TotalOverall
}

func (r *RatingCache) addToRecent(userId uint32, s, v, st uint32) {

	if r.RecentEvaluation == nil {
		r.RecentEvaluation = utils.NewRingBuffer[UserRating](RECENT_EVALUATION_RING_SIZE)
	}

	r.RecentEvaluation.Add(UserRating{
		User: int(userId),
		Rating: Rating{
			Story:      uint8(s),
			Visuals:    uint8(v),
			Soundtrack: uint8(st),
			Overall:    uint8((s + v + st) / 3),
		},
	})
}
