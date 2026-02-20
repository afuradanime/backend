package services

import (
	"context"
	"strconv"
	"time"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/interfaces"
	"github.com/afuradanime/backend/internal/core/utils"
)

type UserService struct {
	userRepository interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) *UserService {
	return &UserService{userRepository: repo}
}

func (s *UserService) GetUsers(ctx context.Context, pageNumber, pageSize int) ([]*domain.User, utils.Pagination, error) {
	return s.userRepository.GetUsers(ctx, pageNumber, pageSize)
}

func (s *UserService) SearchByUsername(ctx context.Context, username string, pageNumber, pageSize int) ([]*domain.User, utils.Pagination, error) {
	return s.userRepository.SearchByUsername(ctx, username, pageNumber, pageSize)
}
func (s *UserService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	return s.userRepository.GetUserById(ctx, id)
}

func (s *UserService) GetUserByProvider(ctx context.Context, provider string, providerID string) (*domain.User, error) {
	return s.userRepository.GetUserByProvider(ctx, provider, providerID)
}

func (s *UserService) RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// TODO: Check if user with same email or username already exists before creating a new one
	added_user, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return added_user, nil
}

func (s *UserService) UpdatePersonalInfo(ctx context.Context, id int, email *string, username *string, location *string, pronouns *string, socials *[]string, birthday *time.Time, allowsFR, allowsRec *bool) error {
	user, err := s.GetUserByID(ctx, id)
	if err != nil || user == nil {
		return err
	}

	if email != nil {
		if err := user.UpdateEmail(*email); err != nil {
			return err
		}
	}
	if username != nil {
		if err := user.UpdateUsername(*username); err != nil {
			return err
		}
	}
	if location != nil {
		user.UpdateLocation(*location)
	}
	if pronouns != nil {
		user.UpdatePronouns(*pronouns)
	}
	if socials != nil {
		if err := user.UpdateSocials(*socials); err != nil {
			return err
		}
	}
	if birthday != nil {
		user.UpdateBirthday(*birthday)
	}
	if allowsFR != nil {
		user.UpdateAllowsFriendRequests(*allowsFR)
	}
	if allowsRec != nil {
		user.UpdateAllowsRecommendations(*allowsRec)
	}

	return s.userRepository.UpdateUser(ctx, user)
}

func (s *UserService) RestrictAccount(ctx context.Context, id int, canPost, canTranslate bool) error {
	user, err := s.GetUserByID(ctx, id)
	if err != nil || user == nil {
		return err
	}

	if user.HasRole(value.UserRoleAdmin) {
		// Népia
		return domain_errors.CantRestrictAnAdmin{}
	}

	user.RestrictAccesses(canPost, canTranslate)
	return s.userRepository.UpdateUser(ctx, user)
}

func (s *UserService) UpdateLastLogin(ctx context.Context, id int) error {
	user, err := s.GetUserByID(ctx, id)
	if err != nil || user == nil {
		return err
	}
	user.UpdateLastLogin()
	return s.userRepository.UpdateUser(ctx, user)
}

func (s *UserService) RewardBadge(ctx context.Context, moderatorID int, targetUserID int, badge value.UserBadges) error {
	user, err := s.GetUserByID(ctx, targetUserID)
	if err != nil || user == nil {
		return domain_errors.UserNotFoundError{UserID: strconv.Itoa(targetUserID)}
	}
	user.RewardBadge(badge)
	return s.userRepository.UpdateUser(ctx, user)
}
