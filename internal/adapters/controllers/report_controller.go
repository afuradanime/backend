package controllers

import (
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-fuego/fuego"
)

type UserReportController struct {
	reportService interfaces.UserReportService
}

func NewUserReportController(reportService interfaces.UserReportService) *UserReportController {
	return &UserReportController{reportService: reportService}
}

type SubmitReportBody struct {
	Reason value.ReportReason `json:"Reason"`
}

func (c *UserReportController) SubmitReport(ctx fuego.ContextWithBody[SubmitReportBody]) (any, error) {
	reporterID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	targetID, err := strconv.Atoi(ctx.PathParam("userID"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid user ID"}
	}

	body, err := ctx.Body()
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid request body"}
	}

	if err := c.reportService.SubmitReport(ctx.Context(), body.Reason, targetID, reporterID); err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}

type ReportListResponse struct {
	Data       []repositories.ReportResult `json:"data"`
	Pagination utils.Pagination            `json:"pagination"`
}

func (c *UserReportController) GetReports(ctx fuego.ContextNoBody) (ReportListResponse, error) {
	if !middlewares.IsLoggedUserOfRole(ctx.Context(), value.UserRoleModerator) {
		return ReportListResponse{}, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	pageNumber, pageSize := utils.GetPaginationParams(ctx, 20)

	results, pagination, err := c.reportService.GetReports(ctx.Context(), pageNumber, pageSize)
	if err != nil {
		return ReportListResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return ReportListResponse{
		Data:       results,
		Pagination: pagination,
	}, nil
}

type ReportByTargetListResponse struct {
	Data       []domain.UserReport `json:"data"`
	Pagination utils.Pagination    `json:"pagination"`
}

func (c *UserReportController) GetReportsByTarget(ctx fuego.ContextNoBody) (ReportByTargetListResponse, error) {
	if !middlewares.IsLoggedUserOfRole(ctx.Context(), value.UserRoleModerator) {
		return ReportByTargetListResponse{}, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	targetID, err := strconv.Atoi(ctx.PathParam("userID"))
	if err != nil {
		return ReportByTargetListResponse{}, fuego.BadRequestError{Detail: "Invalid user ID"}
	}

	pageNumber, pageSize := utils.GetPaginationParams(ctx, 20)

	reports, pagination, err := c.reportService.GetReportsByTarget(ctx.Context(), targetID, pageNumber, pageSize)
	if err != nil {
		return ReportByTargetListResponse{}, fuego.InternalServerError{Detail: err.Error()}
	}

	return ReportByTargetListResponse{
		Data:       reports,
		Pagination: pagination,
	}, nil
}

func (c *UserReportController) DeleteReport(ctx fuego.ContextNoBody) (any, error) {
	if !middlewares.IsLoggedUserOfRole(ctx.Context(), value.UserRoleModerator) {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	modID, ok := middlewares.GetUserIDFromContext(ctx.Context())
	if !ok {
		return nil, fuego.UnauthorizedError{Detail: "Unauthorized"}
	}

	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, fuego.BadRequestError{Detail: "Invalid report ID"}
	}

	if err := c.reportService.DeleteReport(ctx.Context(), id, modID); err != nil {
		return nil, fuego.InternalServerError{Detail: err.Error()}
	}

	return nil, nil
}
