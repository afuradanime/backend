package utils

import (
	"net/http"
	"strconv"
)

const MAX_PAGE_SIZE = 50

type Pagination struct {
	PageNumber int
	PageSize   int
	TotalPages int
}

func GetPaginationParams(r *http.Request, defaultSize int) (pageNumber, pageSize int) {
	pageNumber = 0
	pageSize = defaultSize

	if pageStr := r.URL.Query().Get("pageNumber"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
			pageNumber = p
		}
	}

	if sizeStr := r.URL.Query().Get("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	pageSize = Clamp(pageSize, 1, MAX_PAGE_SIZE)

	return
}
