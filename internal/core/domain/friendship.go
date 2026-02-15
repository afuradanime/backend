package domain

import (
	"time"

	"github.com/afuradanime/backend/internal/core/domain/value"
)

type Friendship struct {
	Initiator string
	Receiver  string
	Status    value.FriendshipStatus

	CreatedAt string
}

func NewFriendRequest(initiator string, receiver string) *Friendship {
	return &Friendship{
		Initiator: initiator,
		Receiver:  receiver,
		Status:    value.FriendshipStatusPending,
	}
}

func NewBlockedUser(initiator string, receiver string) *Friendship {

	currentTime := time.Now().Format(time.RFC3339)

	return &Friendship{
		Initiator: initiator,
		Receiver:  receiver,
		Status:    value.FriendshipStatusBlocked,
		CreatedAt: currentTime,
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
