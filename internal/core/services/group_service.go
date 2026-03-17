package services

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
)

type GroupService struct {
	groupRepository interfaces.GroupRepository
	userRepository  interfaces.UserRepository
}

func NewGroupService(repo interfaces.GroupRepository, userRepository interfaces.UserRepository) *GroupService {
	return &GroupService{
		groupRepository: repo,
		userRepository:  userRepository,
	}
}

func (s *GroupService) GetGroup(ctx context.Context, groupId string) (*domain.Group, error) {
	return s.groupRepository.GetGroup(ctx, groupId)
}

func (s *GroupService) GetGroups(ctx context.Context, pageNumber, pageSize int) ([]*domain.Group, utils.Pagination, error) {
	return s.groupRepository.GetGroups(ctx, pageNumber, pageSize)
}

func (s *GroupService) UpdateGroup(
	ctx context.Context,
	groupId string,
	name, description, rules, icon string,
	user int,
) error {

	group, err := s.groupRepository.GetGroup(ctx, groupId)
	if err != nil || group == nil {
		return domain_errors.GroupNotFoundError{GroupID: groupId}
	}

	if !group.IsModerator(user) {
		return domain_errors.UnauthorizedError{}
	}

	if name != "" {
		if err := group.UpdateName(name); err != nil {
			return err
		}
	}

	if description != "" {
		if err := group.UpdateDescription(description); err != nil {
			return err
		}
	}

	if rules != "" {
		if err := group.UpdateRules(rules); err != nil {
			return err
		}
	}

	if icon != "" {
		if err := group.UpdateIcon(icon); err != nil {
			return err
		}
	}

	return s.groupRepository.UpdateGroup(ctx, group)
}

func (s *GroupService) AddGroupModerator(
	ctx context.Context,
	groupId string,
	moderator int,
	user int,
) error {

	group, err := s.groupRepository.GetGroup(ctx, groupId)
	if err != nil || group == nil {
		return domain_errors.GroupNotFoundError{GroupID: groupId}
	}

	creator, err := s.userRepository.GetUserById(ctx, user)
	if err != nil {
		return err
	}

	if !group.IsModerator(user) || creator.HasRole(value.UserRoleAdmin) {
		return domain_errors.UnauthorizedError{}
	}

	if err := group.AddModerator(moderator); err != nil {
		return err
	}

	return s.groupRepository.UpdateGroup(ctx, group)
}

func (s *GroupService) RemoveGroupModerator(
	ctx context.Context,
	groupId string,
	moderator int,
	user int,
) error {

	group, err := s.groupRepository.GetGroup(ctx, groupId)
	if err != nil || group == nil {
		return domain_errors.GroupNotFoundError{GroupID: groupId}
	}

	creator, err := s.userRepository.GetUserById(ctx, user)
	if err != nil {
		return err
	}

	if !group.IsModerator(user) || creator.HasRole(value.UserRoleAdmin) {
		return domain_errors.UnauthorizedError{}
	}

	if err := group.RemoveModerator(moderator); err != nil {
		return err
	}

	return s.groupRepository.UpdateGroup(ctx, group)
}
