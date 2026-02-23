package domain_errors

import (
	"errors"
	"strconv"
)

type InvalidRatingErr struct{}

func (e *InvalidRatingErr) Error() string {
	return errors.New("invalid rating").Error()
}

type InvalidEpisodeCountErr struct{}

func (e *InvalidEpisodeCountErr) Error() string {
	return errors.New("invalid episode count").Error()
}

type NotesLengthTooLong struct {
	MaxLength int
}

func (e *NotesLengthTooLong) Error() string {
	return errors.New("New notes length exceeds maximum allowed length of " + strconv.Itoa(e.MaxLength)).Error()
}

type AnimeAlreadyInListError struct {
	UserID  string
	AnimeID string
}

func (e *AnimeAlreadyInListError) Error() string {
	return errors.New("Anime with ID " + e.AnimeID + " is already in user " + e.UserID + "'s list").Error()
}

type AnimeNotInListError struct {
	UserID  string
	AnimeID string
}

func (e *AnimeNotInListError) Error() string {
	return errors.New("Anime with ID " + e.AnimeID + " is not in user " + e.UserID + "(id) list").Error()
}
