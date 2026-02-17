package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
)

type FriendshipService interface {
	SendFriendRequest(ctx context.Context, initiator int, receiver int) error
	AcceptFriendRequest(ctx context.Context, initiator int, receiver int) error
	DeclineFriendRequest(ctx context.Context, initiator int, receiver int) error
	BlockUser(ctx context.Context, initiator int, receiver int) error
	GetFriendList(ctx context.Context, userId int) ([]domain.User, error)
	GetPendingFriendRequests(ctx context.Context, userId int) ([]domain.User, error)
}

type FriendshipRepository interface {
	CreateFriendship(ctx context.Context, friendship *domain.Friendship) error
	GetFriendship(ctx context.Context, initiator int, receiver int) (*domain.Friendship, error)
	UpdateFriendshipStatus(ctx context.Context, initiator int, receiver int, status value.FriendshipStatus) error
	GetFriends(ctx context.Context, userId int) ([]int, error)
	GetPendingFriendRequests(ctx context.Context, userId int) ([]int, error)
}
