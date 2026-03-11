package utils

import (
	"strconv"

	"github.com/go-fuego/fuego"
)

const MAX_PAGE_SIZE = 50

type Pagination struct {
	PageNumber int
	PageSize   int
	TotalPages int
}

func GetPaginationParams(ctx fuego.ContextNoBody, defaultSize int) (pageNumber, pageSize int) {
	pageNumber = 1
	pageSize = defaultSize

	if pageStr := ctx.QueryParam("pageNumber"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 1 {
			pageNumber = p
		}
	}

	if sizeStr := ctx.QueryParam("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	pageSize = Clamp(pageSize, 1, MAX_PAGE_SIZE)

	return
}
