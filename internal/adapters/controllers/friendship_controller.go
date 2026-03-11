package controllers

import (
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-fuego/fuego"
)

type FriendshipController struct {
	friendshipService interfaces.FriendshipService
}

func NewFriendshipController(friendshipService interfaces.FriendshipService) *FriendshipController {
	return &FriendshipController{
		friendshipService: friendshipService,
	}
}

func (c *FriendshipController) SendFriendRequest(ctx fuego.ContextNoBody) (any, error) {

	initiator, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	receiver, err := strconv.Atoi(ctx.PathParam("receiver"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid receiver ID"}
	}

	if err := c.friendshipService.SendFriendRequest(ctx.Context(), initiator, receiver); err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}

func (c *FriendshipController) AcceptFriendRequest(ctx fuego.ContextNoBody) (any, error) {
	receiver, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	initiator, err := strconv.Atoi(ctx.PathParam("initiator"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid initiator ID"}
	}

	if err := c.friendshipService.AcceptFriendRequest(ctx.Context(), initiator, receiver); err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}

func (c *FriendshipController) DeclineFriendRequest(ctx fuego.ContextNoBody) (any, error) {
	receiver, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	initiator, err := strconv.Atoi(ctx.PathParam("initiator"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid initiator ID"}
	}

	if err := c.friendshipService.DeclineFriendRequest(ctx.Context(), initiator, receiver); err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}

func (c *FriendshipController) BlockUser(ctx fuego.ContextNoBody) (any, error) {
	initiator, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	receiver, err := strconv.Atoi(ctx.PathParam("receiver"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid receiver ID"}
	}

	if err := c.friendshipService.BlockUser(ctx.Context(), initiator, receiver); err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}

type ListFriendsResponse struct {
	Friends    []domain.User    `json:"data"`
	Pagination utils.Pagination `json:"pagination"`
}

func (c *FriendshipController) ListFriends(ctx fuego.ContextNoBody) (ListFriendsResponse, error) {
	targetUser, err := strconv.Atoi(ctx.PathParam("userID"))
	if err != nil {
		return ListFriendsResponse{}, fuego.BadRequestError{Detail: "Invalid user ID"}
	}

	pageNumber, pageSize := utils.GetPaginationParams(ctx, 50)

	friends, pagination, err := c.friendshipService.GetFriendList(ctx.Context(), targetUser, pageNumber, pageSize)
	if err != nil {
		return ListFriendsResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return ListFriendsResponse{
		Friends: friends, Pagination: pagination,
	}, nil
}

func (c *FriendshipController) ListPendingFriendRequests(ctx fuego.ContextNoBody) (ListFriendsResponse, error) {
	userID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return ListFriendsResponse{}, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	pageNumber, pageSize := utils.GetPaginationParams(ctx, 50)

	requests, pagination, err := c.friendshipService.GetPendingFriendRequests(ctx.Context(), userID, pageNumber, pageSize)
	if err != nil {
		return ListFriendsResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return ListFriendsResponse{
		Friends: requests, Pagination: pagination,
	}, nil
}

type FriendshipStatusResponse struct {
	Initiator int                    `json:"initiator"`
	Receiver  int                    `json:"receiver"`
	Status    value.FriendshipStatus `json:"status"`
}

func (c *FriendshipController) FetchFriendshipStatus(ctx fuego.ContextNoBody) (FriendshipStatusResponse, error) {
	userA, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return FriendshipStatusResponse{}, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	userB, err := strconv.Atoi(ctx.PathParam("receiver"))
	if err != nil {
		return FriendshipStatusResponse{}, fuego.BadRequestError{Detail: "Invalid receiver ID"}
	}

	friendshipStatus, err := c.friendshipService.FetchFriendshipStatus(ctx.Context(), userA, userB)
	if err != nil {
		return FriendshipStatusResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	status := value.FriendshipStatusNone
	if friendshipStatus != nil {
		status = friendshipStatus.Status
	}

	return FriendshipStatusResponse{
		Initiator: userA,
		Receiver:  userB,
		Status:    status,
	}, nil

}
