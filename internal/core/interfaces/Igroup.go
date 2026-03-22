package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/utils"
)

type GroupService interface {
	// Not yet defined how to do this
	// CreateGroup(ctx context.Context, name, description, rules, icon string) error
	GetGroup(ctx context.Context, groupId int) (*domain.Group, error)
	GetGroups(ctx context.Context, pageNumber, pageSize int) ([]*domain.Group, utils.Pagination, error)

	UpdateGroup(ctx context.Context, groupId int, name, description, rules, icon string, user int) error
	AddGroupModerator(ctx context.Context, groupId int, moderator, user int) error
	RemoveGroupModerator(ctx context.Context, groupId int, moderator, user int) error
}

type GroupRepository interface {
	CreateGroup(ctx context.Context, group *domain.Group) error
	GetGroup(ctx context.Context, groupId int) (*domain.Group, error)
	GetGroups(ctx context.Context, pageNumber, pageSize int) ([]*domain.Group, utils.Pagination, error)
	UpdateGroup(ctx context.Context, group *domain.Group) error
}
