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

	// save the c information in a go struct
	anime, err := domain.NewAnime(
		uint32(animePtr.id),
		C.GoString(animePtr.url),
		C.GoString(animePtr.title),
		domain.AnimeType(animePtr._type),
		C.GoString(animePtr.source),
		uint32(animePtr.episodes),
		domain.AnimeStatus(animePtr.status),
		bool(animePtr.airing),
		C.GoString(animePtr.duration),
		C.GoString(animePtr.start_date),
		C.GoString(animePtr.end_date),
		domain.SeasonType(animePtr.season.season),
		uint16(animePtr.season.year),
		C.GoString(animePtr.broadcast.day),
		C.GoString(animePtr.broadcast.time),
		C.GoString(animePtr.broadcast.timezone),
		C.GoString(animePtr.image_url),
		C.GoString(animePtr.small_image_url),
		C.GoString(animePtr.large_image_url),
		C.GoString(animePtr.trailer_embed_url),
	)

	if err != nil {
		C.free_anime(&animePtr)
		return nil, err
	}

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

		anime, err := domain.NewAnime(
			uint32(a.id),
			C.GoString(a.url),
			C.GoString(a.title),
			domain.AnimeType(a._type),
			C.GoString(a.source),
			uint32(a.episodes),
			domain.AnimeStatus(a.status),
			bool(a.airing),
			C.GoString(a.duration),
			C.GoString(a.start_date),
			C.GoString(a.end_date),
			domain.SeasonType(a.season.season),
			uint16(a.season.year),
			C.GoString(a.broadcast.day),
			C.GoString(a.broadcast.time),
			C.GoString(a.broadcast.timezone),
			C.GoString(a.image_url),
			C.GoString(a.small_image_url),
			C.GoString(a.large_image_url),
			C.GoString(a.trailer_embed_url),
		)

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

		anime, err := domain.NewAnime(
			uint32(a.id),
			C.GoString(a.url),
			C.GoString(a.title),
			domain.AnimeType(a._type),
			C.GoString(a.source),
			uint32(a.episodes),
			domain.AnimeStatus(a.status),
			bool(a.airing),
			C.GoString(a.duration),
			C.GoString(a.start_date),
			C.GoString(a.end_date),
			domain.SeasonType(a.season.season),
			uint16(a.season.year),
			C.GoString(a.broadcast.day),
			C.GoString(a.broadcast.time),
			C.GoString(a.broadcast.timezone),
			C.GoString(a.image_url),
			C.GoString(a.small_image_url),
			C.GoString(a.large_image_url),
			C.GoString(a.trailer_embed_url),
		)

		if err != nil {
			C.free_partial_anime_array(animeArray, count)
			return nil, err
		}

		results[i] = anime
	}

	C.free_partial_anime_array(animeArray, count)
	return results, nil
}
