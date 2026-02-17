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
	"github.com/afuradanime/backend/internal/core/domain/value"
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

		anime, err := r.animeMapper.CToGoPartial(unsafe.Pointer(a))

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

		anime, err := r.animeMapper.CToGoPartial(unsafe.Pointer(a))

		if err != nil {
			C.free_partial_anime_array(animeArray, count)
			return nil, err
		}

		results[i] = anime
	}

	C.free_partial_anime_array(animeArray, count)
	return results, nil
}

func (r *AnimeRepository) FetchStudioByID(studioID uint32, pageNumber, pageSize int) (*value.Studio, []*domain.Anime, error) {
	var studioPtr C.studio_t
	var count C.uint
	var animeArray *C.partial_anime_t
	var page = C.pageable_t{
		page_number: C.ushort(pageNumber),
		page_size:   C.ushort(pageSize),
	}

	rc := C.fetch_studio_by_id(C.uint(studioID), &studioPtr, page, &count, &animeArray)
	if rc != 0 {
		return nil, nil, errors.New("No studio found with id " + strconv.Itoa(int(studioID)))
	}
	defer C.free_studio(&studioPtr)

	studio, err := r.animeMapper.CToGoStudio(unsafe.Pointer(&studioPtr))
	if err != nil {
		return nil, nil, err
	}

	if count == 0 {
		return studio, []*domain.Anime{}, nil
	}

	animeSlice := unsafe.Slice(animeArray, count)
	results := make([]*domain.Anime, count)
	for i := 0; i < int(count); i++ {
		a := &animeSlice[i]
		anime, err := r.animeMapper.CToGoPartial(unsafe.Pointer(a))
		if err != nil {
			C.free_partial_anime_array(animeArray, count)
			return nil, nil, err
		}
		results[i] = anime
	}
	C.free_partial_anime_array(animeArray, count)

	return studio, results, nil
}

func (r *AnimeRepository) FetchProducerByID(producerID uint32, pageNumber, pageSize int) (*value.Producer, []*domain.Anime, error) {
	var producerPtr C.producer_t
	var count C.uint
	var animeArray *C.partial_anime_t
	var page = C.pageable_t{
		page_number: C.ushort(pageNumber),
		page_size:   C.ushort(pageSize),
	}

	rc := C.fetch_producer_by_id(C.uint(producerID), &producerPtr, page, &count, &animeArray)
	if rc != 0 {
		return nil, nil, errors.New("No producer found with id " + strconv.Itoa(int(producerID)))
	}
	defer C.free_producer(&producerPtr)

	producer, err := r.animeMapper.CToGoProducer(unsafe.Pointer(&producerPtr))
	if err != nil {
		return nil, nil, err
	}

	if count == 0 {
		return producer, []*domain.Anime{}, nil
	}

	animeSlice := unsafe.Slice(animeArray, count)
	results := make([]*domain.Anime, count)
	for i := 0; i < int(count); i++ {
		a := &animeSlice[i]
		anime, err := r.animeMapper.CToGoPartial(unsafe.Pointer(a))
		if err != nil {
			C.free_partial_anime_array(animeArray, count)
			return nil, nil, err
		}
		results[i] = anime
	}
	C.free_partial_anime_array(animeArray, count)

	return producer, results, nil
}

func (r *AnimeRepository) FetchLicensorByID(licensorID uint32, pageNumber, pageSize int) (*value.Licensor, []*domain.Anime, error) {
	var licensorPtr C.licensor_t
	var count C.uint
	var animeArray *C.partial_anime_t
	var page = C.pageable_t{
		page_number: C.ushort(pageNumber),
		page_size:   C.ushort(pageSize),
	}

	rc := C.fetch_licensor_by_id(C.uint(licensorID), &licensorPtr, page, &count, &animeArray)
	if rc != 0 {
		return nil, nil, errors.New("No licensor found with id " + strconv.Itoa(int(licensorID)))
	}
	defer C.free_licensor(&licensorPtr)

	licensor, err := r.animeMapper.CToGoLicensor(unsafe.Pointer(&licensorPtr))
	if err != nil {
		return nil, nil, err
	}

	if count == 0 {
		return licensor, []*domain.Anime{}, nil
	}

	animeSlice := unsafe.Slice(animeArray, count)
	results := make([]*domain.Anime, count)
	for i := 0; i < int(count); i++ {
		a := &animeSlice[i]
		anime, err := r.animeMapper.CToGoPartial(unsafe.Pointer(a))
		if err != nil {
			C.free_partial_anime_array(animeArray, count)
			return nil, nil, err
		}
		results[i] = anime
	}
	C.free_partial_anime_array(animeArray, count)

	return licensor, results, nil
}
