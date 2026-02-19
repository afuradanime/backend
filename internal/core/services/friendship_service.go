package services

import (
	"context"
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
	f, err := s.FetchFriendshipStatus(ctx, initiator, receiver)
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

	// Check user rights
	if !r.AllowsFriendRequests {
		return domain_errors.UserDoesntAllowFriendRequests{}
	}

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
	f, err := s.friendshipRepository.GetFriendship(ctx, initiator, receiver)
	if err != nil {
		return domain_errors.NotFriendsError{
			Initiator: strconv.Itoa(initiator),
			Receiver:  strconv.Itoa(receiver),
		}
	}

	if receiver != f.Receiver {
		return domain_errors.CannotAcceptAlienRequest{}
	}

	if err := f.Accept(); err != nil {
		return err
	}

	return s.friendshipRepository.UpdateFriendship(ctx, f)
}

func (s *FriendshipService) DeclineFriendRequest(ctx context.Context, initiator int, receiver int) error {
	f, err := s.friendshipRepository.GetFriendship(ctx, initiator, receiver)
	if err != nil {
		return domain_errors.NotFriendsError{
			Initiator: strconv.Itoa(initiator),
			Receiver:  strconv.Itoa(receiver),
		}
	}

	if receiver != f.Receiver {
		return domain_errors.CannotAcceptAlienRequest{}
	}

	if err := f.Decline(); err != nil {
		return err
	}

	return s.friendshipRepository.DeleteFriendship(ctx, initiator, receiver)
}

func (s *FriendshipService) BlockUser(ctx context.Context, initiator int, receiver int) error {
	if initiator == receiver {
		return domain_errors.CannotBlockYourselfError{}
	}

	f, _ := s.FetchFriendshipStatus(ctx, initiator, receiver)

	if f == nil {
		friendship := domain.NewBlockedUser(initiator, receiver)
		return s.friendshipRepository.CreateFriendship(ctx, friendship)
	}

	if err := f.Block(); err != nil {
		return err
	}

	return s.friendshipRepository.UpdateFriendship(ctx, f)
}

func (s *FriendshipService) GetFriendList(ctx context.Context, userId int, pageNumber, pageSize int) ([]domain.User, utils.Pagination, error) {
	return s.friendshipRepository.GetFriends(ctx, userId, pageNumber, pageSize)
}

func (s *FriendshipService) GetPendingFriendRequests(ctx context.Context, userId int, pageNumber, pageSize int) ([]domain.User, utils.Pagination, error) {
	return s.friendshipRepository.GetPendingFriendRequests(ctx, userId, pageNumber, pageSize)
}

func (s *FriendshipService) FetchFriendshipStatus(ctx context.Context, userA int, userB int) (*domain.Friendship, error) {
	// Check if A sent a request to B
	f, err := s.friendshipRepository.GetFriendship(ctx, userA, userB)
	if err == nil && f != nil {
		return f, nil
	}

	// Check if B sent a request to A
	f2, err := s.friendshipRepository.GetFriendship(ctx, userB, userA)
	if err == nil && f2 != nil {
		return f2, nil
	}

	return nil, nil
}
