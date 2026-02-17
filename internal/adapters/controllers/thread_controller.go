package controllers

import (
	"net/http"

	"github.com/afuradanime/backend/internal/core/interfaces"
)

type ThreadController struct {
	threadService interfaces.ThreadsService
}

func NewThreadController(threadService interfaces.ThreadsService) *ThreadController {
	return &ThreadController{
		threadService: threadService,
	}
}

func (c *ThreadController) CreateThreadPost(w http.ResponseWriter, r *http.Request) {
	//TODO
}
