package services

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
)

type FriendshipService struct {
	userRepository       interfaces.UserRepository
	friendshipRepository interfaces.FriendshipRepository
}

func NewFriendshipService(userRepo interfaces.UserRepository, friendshipRepo interfaces.FriendshipRepository) *FriendshipService {
	return &FriendshipService{
		userRepository:       userRepo,
		friendshipRepository: friendshipRepo,
	}
}

func (s *FriendshipService) SendFriendRequest(ctx context.Context, initiator string, receiver string) error {

	// Check if relationship already exists
	f, err := s.friendshipRepository.GetFriendship(ctx, initiator, receiver)
	if err == nil && f != nil {
		// Check if blocked or already friends
		if f.GetStatus() == value.FriendshipStatusBlocked {
			return domain_errors.UserBlockedError{
				Initiator: initiator,
				Receiver:  receiver,
			}
		} else if f.GetStatus() == value.FriendshipStatusAccepted {
			return domain_errors.AlreadyFriendsError{
				Initiator: initiator,
				Receiver:  receiver,
			}
		} else if f.GetStatus() == value.FriendshipStatusPending {
			return domain_errors.FriendRequestAlreadySentError{
				Initiator: initiator,
				Receiver:  receiver,
			}
		}
	}

	// Check user validity
	r, err := s.userRepository.GetUserById(ctx, receiver)
	if err != nil {
		return err
	}

	i, err := s.userRepository.GetUserById(ctx, initiator)
	if err != nil {
		return err
	}

	if r == nil {
		return domain_errors.UserNotFoundError{
			UserID: receiver,
		}
	}

	if i == nil {
		return domain_errors.UserNotFoundError{
			UserID: initiator,
		}
	}

	if receiver == initiator {
		return domain_errors.CannotFriendYourselfError{}
	}

	friendship := domain.NewFriendRequest(initiator, receiver)
	return s.friendshipRepository.CreateFriendship(ctx, friendship)
}

// TODO: For accepting and declining, we should also check if the user is the receiver of the request,
// otherwise we might have a security issue where a user can accept or decline friend requests that are not meant for them
// but we need auth for that
func (s *FriendshipService) AcceptFriendRequest(ctx context.Context, initiator string, receiver string) error {

	// Check if friend request exists
	f, err := s.friendshipRepository.GetFriendship(ctx, initiator, receiver)
	if err != nil {
		return domain_errors.NotFriendsError{
			Initiator: initiator,
			Receiver:  receiver,
		}
	}

	if f.Status != value.FriendshipStatusPending {
		return domain_errors.CantOperateOnNonPendingRequestError{}
	}

	return s.friendshipRepository.UpdateFriendshipStatus(ctx, initiator, receiver, value.FriendshipStatusAccepted)
}

func (s *FriendshipService) DeclineFriendRequest(ctx context.Context, initiator string, receiver string) error {

	f, err := s.friendshipRepository.GetFriendship(ctx, initiator, receiver)
	if err != nil {
		return domain_errors.NotFriendsError{
			Initiator: initiator,
			Receiver:  receiver,
		}
	}

	if f.Status != value.FriendshipStatusPending {
		return domain_errors.CantOperateOnNonPendingRequestError{}
	}

	return s.friendshipRepository.UpdateFriendshipStatus(ctx, initiator, receiver, value.FriendshipStatusDeclined)
}

func (s *FriendshipService) BlockUser(ctx context.Context, initiator string, receiver string) error {

	// If not friends, create a new blocked relationship
	if _, err := s.friendshipRepository.GetFriendship(ctx, initiator, receiver); err != nil {
		friendship := domain.NewBlockedUser(initiator, receiver)
		return s.friendshipRepository.CreateFriendship(ctx, friendship)
	}

	if initiator == receiver {
		return domain_errors.CannotBlockYourselfError{}
	}

	return s.friendshipRepository.UpdateFriendshipStatus(ctx, initiator, receiver, value.FriendshipStatusBlocked)
}
func (s *FriendshipService) GetFriendList(ctx context.Context, userId string) ([]domain.User, error) {
	friends, err := s.friendshipRepository.GetFriends(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Get friend details
	friendDetails := make([]domain.User, len(friends))
	for i, friendId := range friends {

		var f *domain.User
		f, err = s.userRepository.GetUserById(ctx, friendId)
		if err != nil {
			return nil, err
		}

		friendDetails[i] = *f
	}

	return friendDetails, nil
}

func (s *FriendshipService) GetPendingFriendRequests(ctx context.Context, userId string) ([]domain.User, error) {
	requests, err := s.friendshipRepository.GetPendingFriendRequests(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Get request details
	requestDetails := make([]domain.User, len(requests))
	for i, requestId := range requests {

		var r *domain.User
		r, err = s.userRepository.GetUserById(ctx, requestId)
		if err != nil {
			return nil, err
		}

		requestDetails[i] = *r
	}

	return requestDetails, nil
}
