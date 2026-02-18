package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
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

func (s *FriendshipService) SendFriendRequest(ctx context.Context, initiator int, receiver int) error {

	if receiver == initiator {
		return domain_errors.CannotFriendYourselfError{}
	}

	// Check if relationship already exists
	f, err := s.friendshipRepository.GetFriendship(ctx, initiator, receiver)
	if err == nil && f != nil {
		// Check if blocked or already friends
		if f.GetStatus() == value.FriendshipStatusBlocked {
			return domain_errors.UserBlockedError{
				Initiator: strconv.Itoa(initiator),
				Receiver:  strconv.Itoa(receiver),
			}
		} else if f.GetStatus() == value.FriendshipStatusAccepted {
			return domain_errors.AlreadyFriendsError{
				Initiator: strconv.Itoa(initiator),
				Receiver:  strconv.Itoa(receiver),
			}
			// Here there are two possibilities:
			// 1- The user sent the request twice (initiator = f.initiator)
			// 2- The receiver has already sent a request (initiator = f.receiver), in which case, we just want to accept it
		} else if f.GetStatus() == value.FriendshipStatusPending {

			if initiator == f.Receiver {

				// Accept request
				return s.AcceptFriendRequest(ctx, initiator, receiver)

			} else {

				return domain_errors.FriendRequestAlreadySentError{
					Initiator: strconv.Itoa(initiator),
					Receiver:  strconv.Itoa(receiver),
				}
			}

		} else if f.GetStatus() == value.FriendshipStatusDeclined {
			// If the request was declined, we can allow sending a new friend request
			// but we "delete" the old declined request
			return s.friendshipRepository.UpdateFriendshipStatus(ctx, initiator, receiver, value.FriendshipStatusPending)
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

	// TODO: Verificar com um check "allowsFriendRequests"

	if r == nil {
		return domain_errors.UserNotFoundError{
			UserID: strconv.Itoa(receiver),
		}
	}

	if i == nil {
		return domain_errors.UserNotFoundError{
			UserID: strconv.Itoa(initiator),
		}
	}

	friendship := domain.NewFriendRequest(initiator, receiver)
	return s.friendshipRepository.CreateFriendship(ctx, friendship)
}

func (s *FriendshipService) AcceptFriendRequest(ctx context.Context, initiator int, receiver int) error {

	// Check if friend request exists
	f, err := s.friendshipRepository.GetFriendship(ctx, initiator, receiver)
	if err != nil {
		return domain_errors.NotFriendsError{
			Initiator: strconv.Itoa(initiator),
			Receiver:  strconv.Itoa(receiver),
		}
	}

	if f.Status != value.FriendshipStatusPending {
		return domain_errors.CantOperateOnNonPendingRequestError{}
	}

	if receiver != f.Receiver {
		return domain_errors.CannotAcceptAlienRequest{}
	}

	return s.friendshipRepository.UpdateFriendshipStatus(ctx, initiator, receiver, value.FriendshipStatusAccepted)
}

func (s *FriendshipService) DeclineFriendRequest(ctx context.Context, initiator int, receiver int) error {

	f, err := s.friendshipRepository.GetFriendship(ctx, initiator, receiver)
	if err != nil {
		return domain_errors.NotFriendsError{
			Initiator: strconv.Itoa(initiator),
			Receiver:  strconv.Itoa(receiver),
		}
	}

	if f.Status != value.FriendshipStatusPending {
		return domain_errors.CantOperateOnNonPendingRequestError{}
	}

	if receiver != f.Receiver {
		return domain_errors.CannotAcceptAlienRequest{}
	}

	return s.friendshipRepository.UpdateFriendshipStatus(ctx, initiator, receiver, value.FriendshipStatusDeclined)
}

func (s *FriendshipService) BlockUser(ctx context.Context, initiator int, receiver int) error {

	f, err := s.friendshipRepository.GetFriendship(ctx, initiator, receiver)
	if err != nil && !errors.Is(err, domain_errors.UserNotFoundError{}) {
		return err
	}

	// If not friends, create a new blocked relationship
	if f == nil {
		friendship := domain.NewBlockedUser(initiator, receiver)
		return s.friendshipRepository.CreateFriendship(ctx, friendship)
	}

	if initiator == receiver {
		return domain_errors.CannotBlockYourselfError{}
	}

	return s.friendshipRepository.UpdateFriendshipStatus(ctx, initiator, receiver, value.FriendshipStatusBlocked)
}
func (s *FriendshipService) GetFriendList(ctx context.Context, userId int, pageNumber, pageSize int) ([]domain.User, utils.Pagination, error) {
	friends, pagination, err := s.friendshipRepository.GetFriends(ctx, userId, pageNumber, pageSize)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	// Get friend details
	friendDetails := make([]domain.User, len(friends))
	for i, friendId := range friends {

		var f *domain.User
		f, err = s.userRepository.GetUserById(ctx, friendId)
		if err != nil {
			return nil, utils.Pagination{}, err
		}

		friendDetails[i] = *f
	}

	return friendDetails, pagination, nil
}

func (s *FriendshipService) GetPendingFriendRequests(ctx context.Context, userId int, pageNumber, pageSize int) ([]domain.User, utils.Pagination, error) {
	requests, pagination, err := s.friendshipRepository.GetPendingFriendRequests(ctx, userId, pageNumber, pageSize)
	if err != nil {
		return nil, utils.Pagination{}, err
	}

	// Get request details
	requestDetails := make([]domain.User, len(requests))
	for i, requestId := range requests {

		var r *domain.User
		r, err = s.userRepository.GetUserById(ctx, requestId)
		if err != nil {
			return nil, utils.Pagination{}, err
		}

		requestDetails[i] = *r
	}

	return requestDetails, pagination, nil
}

func (s *FriendshipService) AreFriends(ctx context.Context, userA int, userB int) (bool, error) {
	f, err := s.friendshipRepository.GetFriendship(ctx, userA, userB)
	if err == nil && f != nil && f.GetStatus() == value.FriendshipStatusAccepted {
		return true, nil
	}

	return false, nil
}
