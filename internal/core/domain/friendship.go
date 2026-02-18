package domain

import (
	"time"

	"github.com/afuradanime/backend/internal/core/domain/value"
)

type Friendship struct {
	Initiator int                    `json:"initiator" bson:"initiator"`
	Receiver  int                    `json:"receiver" bson:"receiver"`
	Status    value.FriendshipStatus `json:"status" bson:"status"`

	CreatedAt time.Time `json:"CreatedAt" bson:"created_at"`
}

func NewFriendRequest(initiator int, receiver int) *Friendship {
	return &Friendship{
		Initiator: initiator,
		Receiver:  receiver,
		Status:    value.FriendshipStatusPending,
		CreatedAt: time.Now(),
	}
}

func NewBlockedUser(initiator int, receiver int) *Friendship {

	return &Friendship{
		Initiator: initiator,
		Receiver:  receiver,
		Status:    value.FriendshipStatusBlocked,
		CreatedAt: time.Now(),
	}
}

func (f *Friendship) Accept() {
	f.Status = value.FriendshipStatusAccepted
}

func (f *Friendship) Decline() {
	f.Status = value.FriendshipStatusDeclined
}

func (f *Friendship) Block() {
	f.Status = value.FriendshipStatusBlocked
}

func (f *Friendship) GetStatus() value.FriendshipStatus {
	return f.Status
}
