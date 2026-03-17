package domain

import (
	"slices"
	"time"

	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
)

const NOTES_MAX_LEN = 64

// Represents an entry in a user's anime list
type UserAnimeList struct {
	UserID        int            `json:"userId" bson:"user_id"`
	UserListItems []UserListItem `json:"userListItems" bson:"user_list_items"`
}

/*
* MongoDB has a limit of 16MB per document, so we need to be careful with the size of this field.
* Here are some tests I ran (10/3/2026):
  - 100 -> 50kb
  - 1000 -> 500kb
  - 10k -> 5mb
  - 24k -> 12mb

* We set a max length for notes to prevent users from adding excessively long notes that could cause issues when saving the list to the database.
* We also set the anime id to be an unsigned short (65,535) to take less space, we currently have 24k anime in the database, so in 2060 when
* there are no more frontends and APIs will ship their own swagger schemas and users will generate their own UI we can migrate to a larger type if needed
* Episode count is also a uint16 for the same reason no shot we ever get past 65k episodes in an anime
*
* After making AnimeID and EpisodesWatched uint16 instead of uint32: 12MBs -> 10MBs
* We are also wasting space by storing 4 ratings as chars (4 bytes) as they are natural numbers that go from 0 to 10, so that's only 4 bits per rating,
* we have 4 ratings, so 16 bits total, we can pack them into a single uint16 and cut memory usage by half, that also removes struct overhead prob
*
* After packing ratings into a single uint16: 10MBs -> 7MBs
* Another idea is to pack the dates as uint32 unix timestamps, but that would make it impossible to represent dates past 2106,
* I'll be dead by then, let's keep going.
* The bson names were also shortened to reduce mongodb field name overhead
* Current results: 5MBs for 24k anime in the list
*/
type UserListItem struct {
	AnimeID         uint16                    `json:"animeId" bson:"a"`
	Status          value.AnimeListItemStatus `json:"status" bson:"s"`
	EpisodesWatched uint16                    `json:"episodesWatched" bson:"e"`
	Rating          *uint16                   `json:"rating,omitempty" bson:"r,omitempty"`
	Notes           *string                   `json:"notes,omitempty" bson:"n,omitempty"`
	RewatchCount    uint8                     `json:"rewatchCount" bson:"w,omitempty"`
	CreatedAt       uint32                    `json:"createdAt" bson:"c"`
	EditedAt        uint32                    `json:"editedAt,omitempty" bson:"t,omitempty"`
}

// Represents the rating a user gives to an anime in their list, with an overall rating and optional breakdown by category
type Rating struct {
	Overall    uint8 `json:"overall" bson:"overall"`
	Story      uint8 `json:"story" bson:"story"`
	Visuals    uint8 `json:"visuals" bson:"visuals"`
	Soundtrack uint8 `json:"soundtrack" bson:"soundtrack"`
}

func NewPersonalAnimeList(userID int) *UserAnimeList {
	return &UserAnimeList{
		UserID:        userID,
		UserListItems: []UserListItem{},
	}
}

func NewAnimeListItem(userID int, animeID uint32, status value.AnimeListItemStatus) *UserListItem {
	now := time.Now()
	return &UserListItem{
		AnimeID:         uint16(animeID),
		Status:          status,
		EpisodesWatched: 0,
		Rating:          nil,
		Notes:           nil,
		RewatchCount:    0,
		CreatedAt:       uint32(now.Unix()),
	}
}

func (al *UserAnimeList) AddListItem(item UserListItem) {
	al.UserListItems = append(al.UserListItems, item)
}

func (al *UserAnimeList) RemoveListItem(animeID uint32) {

	al.UserListItems = slices.DeleteFunc(al.UserListItems, func(item UserListItem) bool {
		return item.AnimeID == uint16(animeID)
	})
}

func (al *UserAnimeList) RemoveListItemByIndex(index int) {
	al.UserListItems = slices.Delete(al.UserListItems, index, index+1)
}

func (al *UserAnimeList) GetListItem(animeID uint32) (*UserListItem, bool) {
	for i, item := range al.UserListItems {
		if item.AnimeID == uint16(animeID) {
			return &al.UserListItems[i], true
		}
	}
	return nil, false
}

func (al *UserListItem) UpdateRewatchCount(newCount uint8) {
	al.RewatchCount = newCount
	now := time.Now()
	al.EditedAt = uint32(now.Unix())
}

func (al *UserListItem) UpdateNotes(notes string) error {
	if len(notes) > NOTES_MAX_LEN {
		return &domain_errors.NotesLengthTooLong{
			MaxLength: NOTES_MAX_LEN,
		}
	}
	al.Notes = &notes
	return nil
}

func (al *UserListItem) UpdateStatus(newStatus value.AnimeListItemStatus) {
	al.Status = newStatus
	now := time.Now()
	al.EditedAt = uint32(now.Unix())
}

func (al *UserListItem) UpdateProgress(episodesWatched uint32, totalEpisodes uint32) error {
	// cannot watch more episodes than the anime has (unless it's still airing, in which case totalEpisodes would be 0 or nil)
	if episodesWatched > totalEpisodes && al.Status == value.AnimeListItemStatusCompleted {
		return &domain_errors.InvalidEpisodeCountErr{}
	}

	al.EpisodesWatched = uint16(episodesWatched)
	now := time.Now()
	al.EditedAt = uint32(now.Unix())

	// if the user has watched all episodes, mark the anime as completed
	if episodesWatched >= totalEpisodes {
		al.Status = value.AnimeListItemStatusCompleted
	}

	return nil
}

func (al *UserListItem) AddRating(story, visuals, soundtrack uint8) error {
	if story > 10 || visuals > 10 || soundtrack > 10 {
		return &domain_errors.InvalidRatingErr{}
	}

	overall := uint8((uint16(story) + uint16(visuals) + uint16(soundtrack)) / 3)
	r := &Rating{
		Overall:    overall,
		Story:      story,
		Visuals:    visuals,
		Soundtrack: soundtrack,
	}

	al.Rating = new(uint16)
	*al.Rating = RatingToUint16(r)

	now := time.Now()
	al.EditedAt = uint32(now.Unix())

	return nil
}

func (al *UserListItem) RemoveRating() {
	al.Rating = nil
	now := time.Now()
	al.EditedAt = uint32(now.Unix())
}

func RatingToUint16(r *Rating) uint16 {
	if r == nil {
		return 0
	}

	return uint16(r.Overall) | (uint16(r.Story) << 4) | (uint16(r.Visuals) << 8) | (uint16(r.Soundtrack) << 12)
}

func Uint16ToRating(r uint16) *Rating {
	if r == 0 {
		return nil
	}

	return &Rating{
		Overall:    uint8(r & 0xF),
		Story:      uint8((r >> 4) & 0xF),
		Visuals:    uint8((r >> 8) & 0xF),
		Soundtrack: uint8((r >> 12) & 0xF),
	}
}
