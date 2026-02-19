package services

import (
	"context"
	"errors"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/interfaces"
)

type ThreadService struct {
	threadsRepository    interfaces.ThreadsRepository
	userRepository       interfaces.UserRepository
	friendshipRepository interfaces.FriendshipRepository
}

func NewThreadService(repo interfaces.ThreadsRepository, userRepo interfaces.UserRepository, friendshipRepo interfaces.FriendshipRepository) *ThreadService {
	return &ThreadService{threadsRepository: repo, userRepository: userRepo, friendshipRepository: friendshipRepo}
}

func (s *ThreadService) CreateThreadPost(ctx context.Context, contextId int, contextType string, posterId int, content string) (*domain.ThreadPost, error) {
	// Get the thread context, ensuring it exists
	thContext, err := s.threadsRepository.GetThreadContextByID(ctx, contextId)
	if err != nil {
		return nil, errors.New("Trying to post on a thread that doesn't exist: " + err.Error())
	}

	// Get the poster user, ensuring they exist
	poster, err := s.userRepository.GetUserById(ctx, posterId)
	if err != nil {
		return nil, errors.New("Error retrieving poster user: " + err.Error())
	} else if poster == nil {
		return nil, errors.New("Poster user not found.")
	}

	// If this is a profile thread, check if the poster is allowed to post on this profile
	if thContext.ContextType == domain.ContextTypeProfile {
		thContextProfileUser, err := s.userRepository.GetUserById(ctx, thContext.ContextId)
		if err != nil {
			return nil, errors.New("Error retrieving thread context profile user: " + err.Error())
		} else if thContextProfileUser == nil {
			return nil, errors.New("Thread context profile user not found.")
		}

		// Skip friendship check if the poster is the profile owner
		if poster.ID != thContextProfileUser.ID {
			// Check if the poster is allowed to post on this profile thread
			friendship, err := s.friendshipRepository.GetFriendship(ctx, poster.ID, thContextProfileUser.ID)
			if err != nil {
				return nil, errors.New("Error retrieving friendship: " + err.Error())
			}
			// Cannot post on profiles that blocked me
			if friendship.GetStatus() == value.FriendshipStatusBlocked {
				return nil, errors.New("Poster is not allowed to post on this profile thread.")
			}
		}
	}

	// Later add a ContentMiddlewareService for content moderation

	// Passed all checks, create the thread post
	post := domain.NewThreadPost(thContext.ContextId, poster.ID, content)
	return s.threadsRepository.CreateThreadPost(ctx, post)
}

func (s *ThreadService) GetThreadPostsByContext(ctx context.Context, contextId int) ([]*domain.ThreadPost, error) {
	// Verify the thread context exists
	_, err := s.threadsRepository.GetThreadContextByID(ctx, contextId)
	if err != nil {
		return nil, errors.New("Thread context not found: " + err.Error())
	}

	return s.threadsRepository.GetThreadPostsByContext(ctx, contextId)
}
