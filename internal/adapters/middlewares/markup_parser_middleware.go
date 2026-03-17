package middlewares

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/interfaces"
)

var (
	userMentionRegex  = regexp.MustCompile(`@(\d+)`)     // Matches @ followed by digits
	animeMentionRegex = regexp.MustCompile(`\[#(\d+)\]`) // Matches [# followed by digits and ]
)

// Replace custom markup with real data
func ParsePost(
	post *domain.Post,
	ctx context.Context,
	animeRepo interfaces.AnimeService,
	userRepo interfaces.UserRepository,
) {
	if post == nil || post.Text == nil {
		return
	}

	userIds, animeIds := ExtractMentions(*post.Text)
	for _, u := range userIds {
		user, err := userRepo.GetUserById(ctx, u)
		if err != nil || user == nil {
			continue // skip this mention
		}
		replaced := strings.ReplaceAll(
			*post.Text, "@"+strconv.Itoa(u),
			"<a href=\"/profile/"+strconv.Itoa(u)+"\">@"+string(user.Username)+"</a>",
		)
		post.Text = &replaced
	}
	for _, a := range animeIds {
		anime, err := animeRepo.FetchAnimeByID(uint32(a))
		if err != nil || anime == nil { // <-- guard anime
			continue
		}
		replaced := strings.ReplaceAll(
			*post.Text, "[#"+strconv.Itoa(a)+"]",
			"<a href=\"/anime/"+strconv.Itoa(a)+"\">"+anime.Title+"</a>",
		)
		post.Text = &replaced
	}
}

func ExtractMentions(text string) ([]int, []int) {
	userIDs := []int{}
	animeIDs := []int{}

	// Find all @user matches
	userMatches := userMentionRegex.FindAllStringSubmatch(text, -1)
	for _, match := range userMatches {
		if id, err := strconv.Atoi(match[1]); err == nil {
			userIDs = append(userIDs, id)
		}
	}

	// Find all [#anime] matches
	animeMatches := animeMentionRegex.FindAllStringSubmatch(text, -1)
	for _, match := range animeMatches {
		if id, err := strconv.Atoi(match[1]); err == nil {
			animeIDs = append(animeIDs, id)
		}
	}

	return uniqueInts(userIDs), uniqueInts(animeIDs)
}

// Helper to avoid duplicate IDs if a user tags the same anime twice
func uniqueInts(input []int) []int {
	u := make([]int, 0, len(input))
	m := make(map[int]bool)
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}
