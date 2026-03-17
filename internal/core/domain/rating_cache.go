package domain

import (
	"github.com/afuradanime/backend/internal/core/utils"
)

const RECENT_EVALUATION_RING_SIZE = 5

/*
* The rating cache is a denormalized structure that stores the user's ratings for an anime,
* allowing for quick retrieval without needing to calculate the overall rating from individual components every time.
* It is updated whenever the user changes their ratings for story, visuals, or soundtrack,
* and it provides an overall rating that can be used for sorting and displaying in the UI.

* This is shadow data, it is managed automatically in the background, all external perations should be
* fetching the values only
 */
type RatingCache struct {
	AnimeID int `json:"animeId" bson:"anime_id"`

	TotalOverall float32 `json:"overall" bson:"overall"`

	TotalStory      uint32 `json:"story" bson:"story"`
	TotalVisuals    uint32 `json:"visuals" bson:"visuals"`
	TotalSoundtrack uint32 `json:"soundtrack" bson:"soundtrack"`

	RecentEvaluation *utils.RingBuffer[UserRating] `json:"recentEvals" bson:"recent_evals"`

	// Number of users who have rated this anime
	UserCounter int `json:"user_counter" bson:"user_counter"`
}

type UserRating struct {
	user   int
	rating Rating
}

func NewRatingCache(animeID int) *RatingCache {
	return &RatingCache{
		AnimeID:          animeID,
		RecentEvaluation: utils.NewRingBuffer[UserRating](RECENT_EVALUATION_RING_SIZE),
	}
}

func (r *RatingCache) UpdateCache(userId, story, visuals, soundtrack uint32) {
	// Update totals
	r.TotalStory += story
	r.TotalVisuals += visuals
	r.TotalSoundtrack += soundtrack

	// Increment user counter
	r.UserCounter++

	// Recalculate overall rating (simple average of the three components)
	r.TotalOverall = float32((r.TotalStory + r.TotalVisuals + r.TotalSoundtrack) / uint32(3.0*r.UserCounter))

	// Add recent rating
	// TODO: update and remove
	// r.RecentEvaluation.Add(UserRating{
	// 	user: int(userId),
	// 	rating: Rating{
	// 		Story:      uint8(story),
	// 		Visuals:    uint8(visuals),
	// 		Soundtrack: uint8(soundtrack),
	// 		Overall:    uint8((story + visuals + soundtrack) / 3.0),
	// 	},
	// })
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
		r.TotalOverall = float32((r.TotalStory + r.TotalVisuals + r.TotalSoundtrack) / uint32(3.0*r.UserCounter))
	} else {
		r.TotalOverall = 0
	}
}

func (r *RatingCache) UpdateExistingRating(userId, oldStory, oldVisuals, oldSoundtrack, newStory, newVisuals, newSoundtrack uint32) {
	// Update totals by removing old ratings and adding new ratings
	r.TotalStory = r.TotalStory - oldStory + newStory
	r.TotalVisuals = r.TotalVisuals - oldVisuals + newVisuals
	r.TotalSoundtrack = r.TotalSoundtrack - oldSoundtrack + newSoundtrack

	// Recalculate overall rating
	if r.UserCounter > 0 {
		r.TotalOverall = float32((r.TotalStory + r.TotalVisuals + r.TotalSoundtrack) / uint32(3.0*r.UserCounter))
	} else {
		r.TotalOverall = 0
	}
}

func (r *RatingCache) GetOverallRating() float32 {
	return r.TotalOverall
}
