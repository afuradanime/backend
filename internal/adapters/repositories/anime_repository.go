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
	"unsafe"

	"github.com/afuradanime/backend/internal/core/domain"
)

type AnimeRepository struct{}

func NewAnimeRepository() *AnimeRepository {
	return &AnimeRepository{}
}

func (r *AnimeRepository) FetchAnimeByID(animeID uint32) (*domain.Anime, error) {

	// Query the dll for an anime id, get the pointed data
	var animePtr C.anime_t
	rc := C.fetch_anime_by_id(C.uint(animeID), &animePtr)

	if rc != 0 {
		return nil, errors.New("There's no Anime with id " + strconv.Itoa(int(animeID)))
	}

	// decode the duration pointer - copy to Go memory before freeing C memory
	var duration *float32 = nil
	if animePtr.duration_value != nil {
		durationVal := float32(*animePtr.duration_value)
		duration = &durationVal
	}

	// save the c information in a go struct
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

	// now it's handled by the go GC
	// so we may free it
	C.free_anime(&animePtr)

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
	var animeArray *C.anime_t

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
		animePtr := &animeSlice[i]

		var duration *float32 = nil
		if animePtr.duration_value != nil {
			// Copy the value to Go memory before freeing C memory
			durationVal := float32(*animePtr.duration_value)
			duration = &durationVal
		}

		results[i] = domain.NewAnime(
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
	}

	C.free_anime_array(animeArray, C.uint(count))

	return results, nil
}

func (r *AnimeRepository) FetchAnimeThisSeason() ([]*domain.Anime, error) {
	var count C.uint
	var animeArray *C.anime_t

	var rc = C.fetch_anime_this_season(&count, &animeArray)

	if rc != 0 {
		return nil, errors.New("Failed to fetch anime for this season")
	}

	if count == 0 {
		return []*domain.Anime{}, nil
	}

	// Convert C array in Go slice
	var results []*domain.Anime
	animeSlice := unsafe.Slice(animeArray, count)
	results = make([]*domain.Anime, count)

	for i := 0; i < int(count); i++ {
		animePtr := &animeSlice[i]
		var duration *float32 = nil
		if animePtr.duration_value != nil {
			// Copy the value to Go memory before freeing C memory
			durationVal := float32(*animePtr.duration_value)
			duration = &durationVal
		}
		results[i] = domain.NewAnime(
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
	}

	C.free_anime_array(animeArray, C.uint(count))

	return results, nil
}
