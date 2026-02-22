package domain

import "time"

const MAX_RECOMMENDATION_STACK = 20

type Recommendation struct {
	Initiator int `json:"initiator" bson:"initiator"`
	Receiver  int `json:"receiver" bson:"receiver"`
	Anime     int `json:"Anime" bson:"anime"`

	CreatedAt time.Time `json:"CreatedAt" bson:"created_at"`
}

func NewRecommendation(initiator int, receiver int, anime int) *Recommendation {
	return &Recommendation{
		Initiator: initiator,
		Receiver:  receiver,
		Anime:     anime,
		CreatedAt: time.Now(),
	}
}
