package repositories

/*
#cgo LDFLAGS: -L${SRCDIR}/../../../drivers -lanime_facts
#cgo CFLAGS: -I${SRCDIR}/../../../../anime-facts-core/include

#include "anime_facts_api.h"
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"strconv"

	"github.com/afuradanime/backend/internal/core/domain"
)

type AnimeRepository struct{}

func NewAnimeRepository() *AnimeRepository {
	return &AnimeRepository{}
}

func (r *AnimeRepository) FetchAnimeByID(animeID uint32) (*domain.Anime, error) {
	var animePtr C.anime_t
	rc := C.fetch_anime_by_id(C.uint(animeID), &animePtr)
	if rc != 0 {
		return nil, errors.New("There's no Anime with id " + strconv.Itoa(int(animeID)))
	}

	var duration *float32 = nil
	if animePtr.duration_value != nil {
		duration = (*float32)(animePtr.duration_value)
	}

	var anime = domain.NewAnime(
		uint32(animePtr.id),
		C.GoString(animePtr.sources),
		C.GoString(animePtr.title),
		uint32(animePtr._type),
		uint32(animePtr.episodes),
		uint32(animePtr.status),
		C.GoString(animePtr.picture),
		C.GoString(animePtr.thumbnail),
		duration,
	)

	C.free_anime(&animePtr)

	return anime, nil
}
