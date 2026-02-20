package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/afuradanime/backend/internal/adapters/middlewares"
	"github.com/afuradanime/backend/internal/core/domain/value"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
	"github.com/go-chi/chi/v5"
)

type UserReportController struct {
	reportService interfaces.UserReportService
}

func NewUserReportController(reportService interfaces.UserReportService) *UserReportController {
	return &UserReportController{reportService: reportService}
}

func (c *UserReportController) SubmitReport(w http.ResponseWriter, r *http.Request) {
	reporterID, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var body struct {
		Reason value.ReportReason `json:"Reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.reportService.SubmitReport(r.Context(), body.Reason, targetID, reporterID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *UserReportController) GetReports(w http.ResponseWriter, r *http.Request) {
	if !middlewares.IsLoggedUserOfRole(r, value.UserRoleModerator) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pageNumber, pageSize := utils.GetPaginationParams(r, 20)

	results, pagination, err := c.reportService.GetReports(r.Context(), pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       results,
		"pagination": pagination,
	})
}

func (c *UserReportController) GetReportsByTarget(w http.ResponseWriter, r *http.Request) {
	if !middlewares.IsLoggedUserOfRole(r, value.UserRoleModerator) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	pageNumber, pageSize := utils.GetPaginationParams(r, 20)

	reports, pagination, err := c.reportService.GetReportsByTarget(r.Context(), targetID, pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":       reports,
		"pagination": pagination,
	})
}

func (c *UserReportController) DeleteReport(w http.ResponseWriter, r *http.Request) {
	if !middlewares.IsLoggedUserOfRole(r, value.UserRoleModerator) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	modID, ok := middlewares.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	if err := c.reportService.DeleteReport(r.Context(), id, modID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
