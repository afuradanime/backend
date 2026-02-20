package domain

import (
	"time"

	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
)

type DescriptionTranslation struct {
	ID int `json:"ID" bson:"_id"`

	Anime int `json:"Anime" bson:"anime"`

	TranslatedDescription value.LongStr                      `json:"TranslatedDescription" bson:"translated_description"`
	TranslationStatus     value.DescriptionTranslationStatus `json:"status" bson:"status"`

	CreatedBy int       `json:"CreatedBy" bson:"created_by"`
	CreatedAt time.Time `json:"CreatedAt" bson:"created_at"`

	AcceptedBy *int       `json:"AcceptedBy" bson:"accepted_by,omitempty"`
	AcceptedAt *time.Time `json:"AcceptedAt" bson:"accepted_at,omitempty"`
}

// create a pending translation submission
func NewDescriptionTranslation(animeID int, translatedDescription string, createdBy int) (*DescriptionTranslation, error) {

	newDescription, err := value.NewLongStr(translatedDescription)
	if err != nil {
		return nil, err
	}

	return &DescriptionTranslation{
		// ID will be set by mongo auto-increment
		Anime:                 animeID,
		TranslatedDescription: *newDescription,
		TranslationStatus:     value.DescriptionTranslationPending,
		CreatedBy:             createdBy,
		CreatedAt:             time.Now(),
		AcceptedBy:            nil,
		AcceptedAt:            nil,
	}, nil
}

func (t *DescriptionTranslation) Accept(moderatorID int) error {
	if !t.IsPending() {
		return domain_errors.TranslationNotPendingError{}
	}
	now := time.Now()
	t.AcceptedBy = &moderatorID
	t.AcceptedAt = &now
	t.TranslationStatus = value.DescriptionTranslationApproved
	return nil
}

func (t *DescriptionTranslation) IsPending() bool {
	return t.TranslationStatus == value.DescriptionTranslationPending
}

func (t *DescriptionTranslation) IsAccepted() bool {
	return t.TranslationStatus == value.DescriptionTranslationApproved
}

func (t *DescriptionTranslation) BelongsTo(userID int) bool {
	return t.CreatedBy == userID
}
