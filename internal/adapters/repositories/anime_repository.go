package repositories

/*
#cgo linux LDFLAGS: -L${SRCDIR}/../../../drivers -Wl,-rpath,${SRCDIR}/../../../drivers -lanime_facts
#cgo windows LDFLAGS: -L${SRCDIR}/../../../drivers -lanime_facts
#cgo CFLAGS: -I${SRCDIR}/../../../../anime-facts-core/include

#include "anime_facts_api.h"
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"strconv"
	"unsafe"

	"github.com/afuradanime/backend/internal/adapters/mappers"
	"github.com/afuradanime/backend/internal/core/domain"
)

type AnimeRepository struct {
	animeMapper *mappers.AnimeMapper
}

func NewAnimeRepository() *AnimeRepository {
	return &AnimeRepository{
		animeMapper: mappers.NewAnimeMapper(),
	}
}

func (r *AnimeRepository) FetchAnimeByID(animeID uint32) (*domain.Anime, error) {

	// Query the dll for an anime id, get the pointed data
	var animePtr C.anime_t
	rc := C.fetch_anime_by_id(C.uint(animeID), &animePtr)

	if rc != 0 {
		return nil, errors.New("There's no Anime with id " + strconv.Itoa(int(animeID)))
	}

	// Convert the C struct to a Go struct
	anime, err := r.animeMapper.CtoGo(unsafe.Pointer(&animePtr))

	// now it's handled by the go GC
	// so we may free it
	C.free_anime(&animePtr)

	if err != nil {
		return nil, err
	}

	return anime, nil
}

func (r *AnimeRepository) FetchAnimeFromQuery(name string, pageNumber, pageSize int) ([]*domain.Anime, error) {

	// convert Go string to C string
	var cName = C.CString(name)
	defer C.free(unsafe.Pointer(cName)) // set to clean on scope end

	// Create pageable struct
	var page = C.pageable_t{
		page_number: C.ushort(pageNumber),
		page_size:   C.ushort(pageSize),
	}

	var count C.uint
	var animeArray *C.partial_anime_t

	var rc = C.fetch_anime_from_query(cName, page, &count, &animeArray)
	if rc != 0 {
		return nil, errors.New("Failed to fetch anime from query: " + name)
	}

	if count == 0 {
		return []*domain.Anime{}, nil
	}

	// Convert C array in Go slice
	var results []*domain.Anime

	animeSlice := unsafe.Slice(animeArray, count)
	results = make([]*domain.Anime, count)

	for i := 0; i < int(count); i++ {
		a := &animeSlice[i]

		anime, err := r.animeMapper.CtoGo(unsafe.Pointer(a))

		if err != nil {
			C.free_partial_anime_array(animeArray, count)
			return nil, err
		}

		results[i] = anime
	}
	C.free_partial_anime_array(animeArray, C.uint(count))

	return results, nil
}

func (r *AnimeRepository) FetchAnimeThisSeason() ([]*domain.Anime, error) {
	var count C.uint
	var animeArray *C.partial_anime_t

	rc := C.fetch_anime_this_season(&count, &animeArray)
	if rc != 0 {
		return nil, errors.New("Failed to fetch anime for this season")
	}

	if count == 0 {
		return []*domain.Anime{}, nil
	}

	animeSlice := unsafe.Slice(animeArray, count)
	results := make([]*domain.Anime, count)

	for i := 0; i < int(count); i++ {
		a := &animeSlice[i]

		anime, err := r.animeMapper.CtoGo(unsafe.Pointer(a))

		if err != nil {
			C.free_partial_anime_array(animeArray, count)
			return nil, err
		}

		results[i] = anime
	}

	C.free_partial_anime_array(animeArray, count)
	return results, nil
}
