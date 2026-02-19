package domain

import (
	"time"

	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
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

func (f *Friendship) Block() error {
	if f.Status == value.FriendshipStatusBlocked {
		return domain_errors.AlreadyBlocked{}
	}
	f.Status = value.FriendshipStatusBlocked
	return nil
}

func (f *Friendship) Accept() error {
	if f.Status != value.FriendshipStatusPending {
		return domain_errors.CantOperateOnNonPendingRequestError{}
	}
	f.Status = value.FriendshipStatusAccepted
	return nil
}

func (f *Friendship) Decline() error {
	if f.Status != value.FriendshipStatusPending {
		return domain_errors.CantOperateOnNonPendingRequestError{}
	}
	f.Status = value.FriendshipStatusDeclined
	return nil
}
func (f *Friendship) GetStatus() value.FriendshipStatus {
	return f.Status
}
