package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
)

type PostService struct {
	postRepo          interfaces.PostRepository
	userRepo          interfaces.UserRepository
	friendshipService interfaces.FriendshipService
}

func NewPostService(postRepo interfaces.PostRepository, userRepo interfaces.UserRepository, friendshipService interfaces.FriendshipService) *PostService {
	return &PostService{
		postRepo:          postRepo,
		userRepo:          userRepo,
		friendshipService: friendshipService,
	}
}

func (s *PostService) GetPostById(ctx context.Context, postId string) (*domain.Post, error) {
	// no domain rules, anyone can fetch any post and see it, even if it's deleted
	// the only thing we're limiting are interactions like posting and replying
	return s.postRepo.GetPostById(ctx, postId)
}

func (s *PostService) GetPostReplies(ctx context.Context, parentID string) ([]*domain.Post, error) {
	return s.postRepo.GetPostReplies(ctx, parentID)
}

// Create a top most post, which is a post that is not a reply to another post
func (s *PostService) CreatePost(ctx context.Context, parentId string, parentType value.PostParentType, text string, posterId int) (*domain.Post, error) {
	// There's a ton of rules here that will also apply to replies so let's break them down:

	// Checkpoint 0 - Poster Exists ?
	poster, err := s.userRepo.GetUserById(ctx, posterId)
	if err != nil {
		return nil, errors.New("Failed to fetch poster's user data: " + err.Error())
	} else if poster == nil {
		return nil, domain_errors.UserNotFoundError{UserID: strconv.Itoa(posterId)}
	}

	// Checkpoint 1 - System Context Blockage
	if !poster.CanPost {
		return nil, errors.New("User is not allowed to post")
	}

	// Checkpoint 2 - Profile Context Blockage
	if parentType == value.ParentTypeUser {
		userOfProfileId, err := strconv.Atoi(parentId)
		if err != nil {
			return nil, errors.New("Invalid parent id: " + err.Error())
		}
		// TODO: we could also validate that the user of the profile exists but meh
		if userOfProfileId != poster.ID {
			friendship, err := s.friendshipService.FetchFriendshipStatus(ctx, poster.ID, userOfProfileId)
			if err != nil {
				return nil, errors.New("failed to fetch friendship status: " + err.Error())
			} else if friendship.Status == value.FriendshipStatusBlocked {
				return nil, domain_errors.UserBlockedError{
					Initiator: strconv.Itoa(friendship.Initiator),
					Receiver:  strconv.Itoa(friendship.Receiver),
				}
			}
		}
	}

	// Checkpoint X - Anything else, we could add group blockage or forum blockage, etc..
	// Maybe even a content filter checkpoint should be added here to prevent certain words, etc...

	// All checkpoints cleared, we can create the post
	newPost := domain.NewPost(
		parentId, parentType, text, poster.ID,
	)

	return s.postRepo.CreatePost(ctx, newPost)
}

func (s *PostService) CreateReply(ctx context.Context, replyToPostID string, text string, createdBy int) (*domain.Post, error) {
	// Shares some of the validations of a top most post (createPost) but also has some of its own:

	// Checkpoint 0 - The post being replied to must exist
	postBeingRepliedTo, err := s.postRepo.GetPostById(ctx, replyToPostID)
	if err != nil {
		return nil, errors.New("Failed to fetch post being replied to: " + err.Error())
	} else if postBeingRepliedTo == nil {
		return nil, domain_errors.PostNotFoundError{PostID: replyToPostID}
	}

	// Checkpoint 1 - The post being replied to must not be deleted
	if postBeingRepliedTo.IsDeleted() {
		return nil, domain_errors.PostDeletedError{PostID: replyToPostID}
	}

	// Checkpoint 2 - Owner of the post being replied to must not have blocked the replier
	ownerOfPostBeingRepliedTo, err := s.userRepo.GetUserById(ctx, *postBeingRepliedTo.CreatedBy)
	if err != nil {
		return nil, errors.New("Failed to fetch owner of post being replied to: " + err.Error())
	} else if ownerOfPostBeingRepliedTo == nil {
		return nil, domain_errors.UserNotFoundError{UserID: strconv.Itoa(*postBeingRepliedTo.CreatedBy)}
	}
	friendship, err := s.friendshipService.FetchFriendshipStatus(ctx, createdBy, ownerOfPostBeingRepliedTo.ID)
	if err != nil {
		return nil, errors.New("failed to fetch friendship status: " + err.Error())
	} else if friendship.Status == value.FriendshipStatusBlocked {
		return nil, domain_errors.UserBlockedError{
			Initiator: strconv.Itoa(friendship.Initiator),
			Receiver:  strconv.Itoa(friendship.Receiver),
		}
	}

	// Checkpoint 3 - The same checkpoints as createPost apply to replies as well
	reply, err := s.CreatePost(ctx, replyToPostID, value.ParentTypePost, text, createdBy)
	if err != nil {
		return nil, err
	}

	// Update the parent post's replies list
	if err := s.postRepo.AddReplyToPost(ctx, replyToPostID, reply.ID); err != nil {
		return nil, errors.New("failed to update parent post replies: " + err.Error())
	}

	return reply, nil
}

func (s *PostService) DeletePost(ctx context.Context, postID string, deleterId int) error {
	post, err := s.postRepo.GetPostById(ctx, postID)
	if err != nil {
		return errors.New("Failed to fetch post: " + err.Error())
	} else if post == nil {
		return domain_errors.PostNotFoundError{PostID: postID}
	}

	if post.IsDeleted() {
		return domain_errors.PostDeletedError{PostID: postID}
	}

	if *post.CreatedBy != deleterId {
		return domain_errors.NotPostOwnerError{UserID: strconv.Itoa(deleterId), PostID: postID}
	}

	post.Delete()
	return s.postRepo.UpdatePost(ctx, post)
}
