package services

import (
	"context"

	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/interfaces"
)

type ThreadService struct {
	threadsRepository interfaces.ThreadsRepository
}

func NewThreadService(repo interfaces.ThreadsRepository) *ThreadService {
	return &ThreadService{threadsRepository: repo}
}

func (s *ThreadService) CreateThreadPost(ctx context.Context, context int, userId int, content string) (*domain.ThreadPost, error) {
	post := domain.NewThreadPost(context, userId, content)
	return s.threadsRepository.CreateThreadPost(ctx, post)
}
