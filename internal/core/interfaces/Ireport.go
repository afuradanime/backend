package interfaces

import (
	"context"

	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/utils"
)

type UserReportService interface {
	SubmitReport(ctx context.Context, reason value.ReportReason, targetUserID, reporterID int) error
	GetReports(ctx context.Context, pageNumber, pageSize int) ([]repositories.ReportResult, utils.Pagination, error)
	GetReportsByTarget(ctx context.Context, targetUserID int, pageNumber, pageSize int) ([]domain.UserReport, utils.Pagination, error)
	DeleteReport(ctx context.Context, id int, moderatorID int) error
}

type UserReportRepository interface {
	CreateReport(ctx context.Context, report *domain.UserReport) error
	GetReportByID(ctx context.Context, id int) (*domain.UserReport, error)
	DeleteReport(ctx context.Context, id int) error
	GetReports(ctx context.Context, pageNumber, pageSize int) ([]repositories.ReportResult, utils.Pagination, error)
	GetReportsByTarget(ctx context.Context, targetUserID int, pageNumber, pageSize int) ([]domain.UserReport, utils.Pagination, error)
	CountReportsByTarget(ctx context.Context, targetUserID int) (int, error)
	HasReported(ctx context.Context, reporterID, targetUserID int) (bool, error)
}
