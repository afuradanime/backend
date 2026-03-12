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
		if p, err := strconv.Atoi(pageStr); err == nil {
			pageNumber = p
		}
	}

	if sizeStr := ctx.QueryParam("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil {
			pageSize = s
		}
	}

	pageNumber = ClampBottom(pageNumber, 0)
	pageSize = Clamp(pageSize, 0, MAX_PAGE_SIZE)

	return
}
