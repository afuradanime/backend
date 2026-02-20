package services

import (
	"context"

	"github.com/afuradanime/backend/internal/adapters/repositories"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
)

type UserReportService struct {
	reportRepository interfaces.UserReportRepository
	userRepository   interfaces.UserRepository
}

func NewUserReportService(
	reportRepo interfaces.UserReportRepository,
	userRepo interfaces.UserRepository,
) *UserReportService {
	return &UserReportService{
		reportRepository: reportRepo,
		userRepository:   userRepo,
	}
}

func (s *UserReportService) SubmitReport(ctx context.Context, reason value.ReportReason, targetUserID, reporterID int) error {
	if targetUserID == reporterID {
		return domain_errors.CannotReportYourselfError{}
	}

	// Check target exists
	target, err := s.userRepository.GetUserById(ctx, targetUserID)
	if err != nil || target == nil {
		return domain_errors.UserNotFoundError{}
	}

	// Check not already reported
	already, err := s.reportRepository.HasReported(ctx, reporterID, targetUserID)
	if err != nil {
		return err
	}
	if already {
		return domain_errors.AlreadyReportedError{}
	}

	report := domain.NewUserReport(reason, targetUserID, reporterID)
	return s.reportRepository.CreateReport(ctx, report)
}

func (s *UserReportService) GetReports(ctx context.Context, pageNumber, pageSize int) ([]repositories.ReportResult, utils.Pagination, error) {
	return s.reportRepository.GetReports(ctx, pageNumber, pageSize)
}

func (s *UserReportService) GetReportsByTarget(ctx context.Context, targetUserID int, pageNumber, pageSize int) ([]domain.UserReport, utils.Pagination, error) {
	return s.reportRepository.GetReportsByTarget(ctx, targetUserID, pageNumber, pageSize)
}

func (s *UserReportService) DeleteReport(ctx context.Context, id int, moderatorID int) error {
	report, err := s.reportRepository.GetReportByID(ctx, id)
	if err != nil || report == nil {
		return domain_errors.ReportNotFoundError{}
	}
	return s.reportRepository.DeleteReport(ctx, id)
}
