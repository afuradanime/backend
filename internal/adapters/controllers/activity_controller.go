package controllers

import (
	"strconv"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/go-fuego/fuego"
)

type ActivityController struct {
	tracker *domain.ActivityTracker
}

func NewActivityController(tracker *domain.ActivityTracker) *ActivityController {
	return &ActivityController{tracker: tracker}
}

type UserActivityResponse struct {
	UserID   int `json:"user_id"`
	IsOnline int `json:"is_online"`
}

type ActivityStatsResponse struct {
	OnlineCount int   `json:"online_count"`
	OnlineUsers []int `json:"online_users"`
}

func (c *ActivityController) IsUserOnline(ctx fuego.ContextNoBody) (UserActivityResponse, error) {
	userID, err := strconv.Atoi(ctx.PathParam("userID"))
	if err != nil {
		return UserActivityResponse{}, fuego.BadRequestError{Detail: "Invalid user ID"}
	}

	return UserActivityResponse{
		UserID:   userID,
		IsOnline: c.tracker.IsActive(userID),
	}, nil
}

func (c *ActivityController) GetActivityStats(ctx fuego.ContextNoBody) (ActivityStatsResponse, error) {
	online := c.tracker.GetActiveUsers()
	return ActivityStatsResponse{
		OnlineCount: len(online),
		OnlineUsers: online,
	}, nil
}
