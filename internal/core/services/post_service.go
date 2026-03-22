package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
)

type PostService struct {
	postRepo          interfaces.PostRepository
	userRepo          interfaces.UserRepository
	friendshipService interfaces.FriendshipService
	groupService      interfaces.GroupService
	animeService      interfaces.AnimeService
}

func NewPostService(
	postRepo interfaces.PostRepository,
	userRepo interfaces.UserRepository,
	friendshipService interfaces.FriendshipService,
	animeService interfaces.AnimeService,
	groupService interfaces.GroupService,
) *PostService {
	return &PostService{
		postRepo:          postRepo,
		userRepo:          userRepo,
		friendshipService: friendshipService,
		animeService:      animeService,
		groupService:      groupService,
	}
}

func (s *PostService) GetPostById(ctx context.Context, postId string) (*domain.Post, error) {
	// no domain rules, anyone can fetch any post and see it, even if it's deleted
	// the only thing we're limiting are interactions like posting and replying
	post, err := s.postRepo.GetPostById(ctx, postId)
	if err != nil {
		return nil, err
	}

	// Parse markdown
	// We do this server side because today user 1 might be called Makoto naegi
	// But tomorrow he may be called Nagito komaeda or something we never know
	middlewares.ParsePost(post, ctx, s.animeService, s.userRepo)

	return post, err
}

func (s *PostService) GetPostReplies(ctx context.Context, parentID string, parentType value.PostParentType) ([]*domain.Post, error) {
	posts, err := s.postRepo.GetPostReplies(ctx, parentID, parentType)

	if err != nil {
		return nil, err
	}

	for _, element := range posts {
		middlewares.ParsePost(element, ctx, s.animeService, s.userRepo)
	}

	return posts, err
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

		if userOfProfileId != poster.ID {
			friendship, err := s.friendshipService.FetchFriendshipStatus(ctx, poster.ID, userOfProfileId)
			if err != nil {
				return nil, errors.New("failed to fetch friendship status: " + err.Error())
			}
			if friendship != nil && friendship.Status == value.FriendshipStatusBlocked {
				return nil, domain_errors.UserBlockedError{
					Initiator: strconv.Itoa(friendship.Initiator),
					Receiver:  strconv.Itoa(friendship.Receiver),
				}
			}
		}
	} else if parentType == value.ParentTypeThread {

		animeId, err := strconv.Atoi(parentId)
		if err != nil {
			return nil, errors.New("Invalid parent id: " + err.Error())
		}

		_, err = s.animeService.FetchAnimeByID(uint32(animeId))
		if err != nil {
			return nil, errors.New("Invalid parent id: " + err.Error())
		}
	} else if parentType == value.ParentTypeGroup {
		groupId, err := strconv.Atoi(parentId)
		if err != nil {
			return nil, errors.New("Invalid parent id: " + err.Error())
		}

		group, err := s.groupService.GetGroup(ctx, groupId)
		if err != nil {
			return nil, errors.New("Invalid parent id: " + err.Error())
		}

		// Check if group private
		if !group.Public {
			if !group.IsModerator(posterId) && !poster.HasRole(value.UserRoleAdmin) {
				return nil, domain_errors.UnauthorizedError{}
			}
		}

	} else {
		return nil, errors.New("Unsupported thread context")
	}

	// Checkpoint 3 - Sanitize input
	cleanText := utils.SanitizeText(text)
	if len(cleanText) == 0 {
		return nil, errors.New("post content cannot be empty after sanitization")
	}

	// Checkpoint X - Anything else, we could add group blockage or forum blockage, etc..
	// Maybe even a content filter checkpoint should be added here to prevent certain words, etc...

	// All checkpoints cleared, we can create the post
	newPost := domain.NewPost(
		parentId, parentType, cleanText, poster.ID,
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
	}
	if friendship != nil && friendship.Status == value.FriendshipStatusBlocked {
		return nil, domain_errors.UserBlockedError{
			Initiator: strconv.Itoa(friendship.Initiator),
			Receiver:  strconv.Itoa(friendship.Receiver),
		}
	}

	cleanText := utils.SanitizeText(text)
	if len(cleanText) == 0 {
		return nil, errors.New("post content cannot be empty after sanitization")
	}

	// Checkpoint 3 - The same checkpoints as createPost apply to replies as well
	reply, err := s.CreatePost(ctx, replyToPostID, value.ParentTypePost, cleanText, createdBy)
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
