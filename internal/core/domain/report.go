package domain

import (
	"time"

	"github.com/afuradanime/backend/internal/core/domain/value"
)

type UserReport struct {
	ID         int                `json:"ID" bson:"_id"`
	Reason     value.ReportReason `json:"Reason" bson:"reason"`
	TargetUser int                `json:"TargetUser" bson:"target_user"`

	CreatedAt time.Time `json:"CreatedAt" bson:"created_at"`
	CreatedBy int       `json:"CreatedBy" bson:"created_by"`
}

func NewUserReport(reason value.ReportReason, targetUser, reporter int) *UserReport {
	return &UserReport{
		Reason:     reason,
		TargetUser: targetUser,
		CreatedAt:  time.Now(),
		CreatedBy:  reporter,
	}
}
