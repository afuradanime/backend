package controllers

import (
	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-fuego/fuego"
)

type GroupController struct {
	groupService interfaces.GroupService
}

func NewGroupController(s interfaces.GroupService) *GroupController {
	return &GroupController{groupService: s}
}

type GroupListResponse struct {
	Data       []*domain.Group  `json:"data"`
	Pagination utils.Pagination `json:"pagination"`
}

func (gc *GroupController) GetGroups(ctx fuego.ContextNoBody) (GroupListResponse, error) {

	pageNumber, pageSize := utils.GetPaginationParams(ctx, 20)

	groups, pagination, err := gc.groupService.GetGroups(ctx.Context(), pageNumber, pageSize)
	if err != nil {
		return GroupListResponse{}, fuego.InternalServerError{Detail: "Failed to retrieve groups"}
	}

	return GroupListResponse{
		Data:       groups,
		Pagination: pagination,
	}, nil
}

func (gc *GroupController) GetGroupByID(ctx fuego.ContextNoBody) (*domain.Group, error) {
	groupID := ctx.PathParam("id")
	if groupID == "" {
		return nil, fuego.BadRequestError{Detail: "Invalid group ID"}
	}

	group, err := gc.groupService.GetGroup(ctx.Context(), groupID)
	if err != nil {
		return nil, fuego.NotFoundError{Detail: err.Error()}
	}

	return group, nil
}

type UpdateGroupBody struct {
	Name        *string `json:"Name"`
	Description *string `json:"Description"`
	Rules       *string `json:"Rules"`
	Icon        *string `json:"Icon"`
}

func (gc *GroupController) UpdateGroup(ctx fuego.ContextWithBody[UpdateGroupBody]) (any, error) {
	userID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	groupID := ctx.PathParam("id")
	if groupID == "" {
		return nil, fuego.BadRequestError{Detail: "Invalid group ID"}
	}

	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid request body"}
	}

	// since service expects strings, unwrap pointers safely
	var name, description, rules, icon string

	if body.Name != nil {
		name = *body.Name
	}
	if body.Description != nil {
		description = *body.Description
	}
	if body.Rules != nil {
		rules = *body.Rules
	}
	if body.Icon != nil {
		icon = *body.Icon
	}

	err = gc.groupService.UpdateGroup(
		ctx.Context(),
		groupID,
		name,
		description,
		rules,
		icon,
		userID,
	)

	if err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}

type ModifyModeratorBody struct {
	ModeratorID int `json:"ModeratorID"`
}

func (gc *GroupController) AddGroupModerator(ctx fuego.ContextWithBody[ModifyModeratorBody]) (any, error) {
	userID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	groupID := ctx.PathParam("id")
	if groupID == "" {
		return nil, fuego.BadRequestError{Detail: "Invalid group ID"}
	}

	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid request body"}
	}

	err = gc.groupService.AddGroupModerator(
		ctx.Context(),
		groupID,
		body.ModeratorID,
		userID,
	)

	if err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}

func (gc *GroupController) RemoveGroupModerator(ctx fuego.ContextWithBody[ModifyModeratorBody]) (any, error) {
	userID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	groupID := ctx.PathParam("id")
	if groupID == "" {
		return nil, fuego.BadRequestError{Detail: "Invalid group ID"}
	}

	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid request body"}
	}

	err = gc.groupService.RemoveGroupModerator(
		ctx.Context(),
		groupID,
		body.ModeratorID,
		userID,
	)

	if err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}
