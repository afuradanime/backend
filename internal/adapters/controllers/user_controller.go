package controllers

import (
	"strconv"
	"time"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-fuego/fuego"
)

type UserController struct {
	userService interfaces.UserService
}

func NewUserController(s interfaces.UserService) *UserController {
	return &UserController{userService: s}
}

type UserListResponse struct {
	Data       []*domain.User   `json:"data"`
	Pagination utils.Pagination `json:"pagination"`
}

func (uc *UserController) GetUsers(ctx fuego.ContextNoBody) (UserListResponse, error) {
	pageNumber, pageSize := utils.GetPaginationParams(ctx, 20)

	users, pagination, err := uc.userService.GetUsers(ctx.Context(), pageNumber, pageSize)
	if err != nil {
		return UserListResponse{}, fuego.InternalServerError{Detail: "Failed to retrieve users"}
	}

	return UserListResponse{
		Data:       users,
		Pagination: pagination,
	}, nil
}

func (uc *UserController) SearchByUsername(ctx fuego.ContextNoBody) (UserListResponse, error) {
	username := ctx.QueryParam("q")
	if username == "" {
		return UserListResponse{}, fuego.BadRequestError{Detail: "Missing search query"}
	}

	pageNumber, pageSize := utils.GetPaginationParams(ctx, 20)

	users, pagination, err := uc.userService.SearchByUsername(ctx.Context(), username, pageNumber, pageSize)
	if err != nil {
		return UserListResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return UserListResponse{
		Data:       users,
		Pagination: pagination,
	}, nil
}

func (uc *UserController) GetUserByID(ctx fuego.ContextNoBody) (*domain.User, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid user ID"}
	}

	user, err := uc.userService.GetUserByID(ctx.Context(), id)
	if err != nil {
		return nil, fuego.NotFoundError{Detail: err.Error()}
	}

	return user, nil
}

type UpdateUserInfoBody struct {
	Email                 	*string   `json:"Email"`
	Username              	*string   `json:"Username"`
	Location              	*string   `json:"Location"`
	Pronouns              	*string   `json:"Pronouns"`
	Socials               	*[]string `json:"Socials"`
	Birthday              	*string   `json:"Birthday"`
	AllowsFriendRequests  	*bool     `json:"AllowsFriendRequests"`
	AllowsRecommendations 	*bool     `json:"AllowsRecommendations"`
	ListPrivate			 	*bool     `json:"ListPrivate"`
	AvatarURL 			  	*string	  `json:"AvatarURL"`
	AcceptedTermsOfService  *bool	  `json:"AcceptedTermsOfService"`
}

func (uc *UserController) UpdateUserInfo(ctx fuego.ContextWithBody[UpdateUserInfoBody]) (any, error) {
	id, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	updateData, err := ctx.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid request body"}
	}

	var birthday *time.Time
	if updateData.Birthday != nil {
		t, err := time.Parse("2006-01-02", *updateData.Birthday)
		if err != nil {
			return nil, fuego.BadRequestError{Detail: "Invalid birthday format, expected YYYY-MM-DD"}
		}
		birthday = &t
	}

	err = uc.userService.UpdatePersonalInfo(
		ctx.Context(), id,
		updateData.Email,
		updateData.Username,
		updateData.Location,
		updateData.Pronouns,
		updateData.Socials,
		birthday,
		updateData.AllowsFriendRequests,
		updateData.AllowsRecommendations,
		updateData.ListPrivate,
		updateData.AvatarURL,
		updateData.AcceptedTermsOfService,
	)
	if err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}

type RestrictAccountBody struct {
	CanPost      bool `json:"CanPost"`
	CanTranslate bool `json:"CanTranslate"`
}

func (uc *UserController) RestrictAccount(ctx fuego.ContextWithBody[RestrictAccountBody]) (any, error) {
	if !middlewares.IsLoggedUserOfRole(ctx.Context(), value.UserRoleModerator) {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	targetID, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid user ID"}
	}

	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid request body"}
	}

	if err := uc.userService.RestrictAccount(ctx.Context(), targetID, body.CanPost, body.CanTranslate); err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}
