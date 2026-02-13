package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
)

type FriendshipService interface {
	SendFriendRequest(ctx context.Context, initiator string, receiver string) error
	AcceptFriendRequest(ctx context.Context, initiator string, receiver string) error
	DeclineFriendRequest(ctx context.Context, initiator string, receiver string) error
	BlockUser(ctx context.Context, initiator string, receiver string) error
	GetFriendList(ctx context.Context, userId string) ([]domain.User, error)
	GetPendingFriendRequests(ctx context.Context, userId string) ([]domain.User, error)
}

type FriendshipRepository interface {
	CreateFriendship(ctx context.Context, friendship *domain.Friendship) error
	GetFriendship(ctx context.Context, initiator string, receiver string) (*domain.Friendship, error)
	UpdateFriendshipStatus(ctx context.Context, initiator string, receiver string, status value.FriendshipStatus) error
	GetFriends(ctx context.Context, userId string) ([]string, error)
	GetPendingFriendRequests(ctx context.Context, userId string) ([]string, error)
}
