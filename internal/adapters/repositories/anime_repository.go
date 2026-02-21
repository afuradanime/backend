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
	"strconv"
	"unsafe"

	"github.com/afuradanime/backend/internal/adapters/mappers"
	"github.com/afuradanime/backend/internal/core/domain"
	"github.com/afuradanime/backend/internal/core/domain/filters"
	"github.com/afuradanime/backend/internal/core/domain/value"
	domain_errors "github.com/afuradanime/backend/internal/core/errors"
	"github.com/afuradanime/backend/internal/core/utils"
)

type AnimeRepository struct {
	animeMapper *mappers.AnimeMapper
}

func buildCFilter(f filters.AnimeFilter) (C.anime_filter_t, []*C.char) {
	var filter C.anime_filter_t
	var toFree []*C.char

	if f.Name != nil {
		cs := C.CString(*f.Name)
		toFree = append(toFree, cs)
		filter.name = cs
	}
	if f.Type != nil {
		ct := C.enum_anime_type(*f.Type)
		filter._type = &ct
	}
	if f.Status != nil {
		cs := C.enum_anime_status(*f.Status)
		filter.status = &cs
	}
	if f.StartDate != nil {
		ct := C.time_t(*f.StartDate)
		filter.start_date = &ct
	}
	if f.EndDate != nil {
		ct := C.time_t(*f.EndDate)
		filter.end_date = &ct
	}
	if f.MinEpisodes != nil {
		cm := C.uint(*f.MinEpisodes)
		filter.min_episodes = &cm
	}
	if f.MaxEpisodes != nil {
		cm := C.uint(*f.MaxEpisodes)
		filter.max_episodes = &cm
	}

	return filter, toFree
}

func freeCFilter(ptrs []*C.char) {
	for _, p := range ptrs {
		C.free(unsafe.Pointer(p))
	}
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
		return nil, domain_errors.AnimeNotFoundError{
			AnimeID: strconv.Itoa(int(animeID)),
		}
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

func (r *AnimeRepository) FetchRandomAnime() (*domain.Anime, error) {

	// Query the dll for an anime id, get the pointed data
	var animePtr C.anime_t
	rc := C.fetch_random_anime(&animePtr)

	if rc != 0 {
		return nil, domain_errors.AnimeNotFoundError{
			AnimeID: "",
		}
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

func (r *AnimeRepository) FetchAnimeFromQuery(filters filters.AnimeFilter, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error) {
	cFilter, toFree := buildCFilter(filters)
	defer freeCFilter(toFree)

	var totalPages C.uint
	var page = C.pageable_t{
		page_number: C.ushort(pageNumber),
		page_size:   C.ushort(pageSize),
	}
	var count C.uint
	var animeArray *C.partial_anime_t

	rc := C.fetch_anime_from_query(cFilter, page, &count, &totalPages, &animeArray)
	if rc != 0 {
		return nil, utils.Pagination{}, domain_errors.AnimeFetchFailedError{}
	}

	if count == 0 {
		return []*domain.Anime{}, utils.Pagination{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			TotalPages: int(totalPages),
		}, nil
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
			return nil, utils.Pagination{}, err
		}

		results[i] = anime
	}
	C.free_partial_anime_array(animeArray, C.uint(count))

	return results, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: int(totalPages),
	}, nil
}

func (r *AnimeRepository) FetchAnimeThisSeason(filters filters.AnimeFilter, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error) {
	cFilter, toFree := buildCFilter(filters)
	defer freeCFilter(toFree)

	var count C.uint
	var animeArray *C.partial_anime_t
	var totalPages C.uint
	var page = C.pageable_t{
		page_number: C.ushort(pageNumber),
		page_size:   C.ushort(pageSize),
	}

	rc := C.fetch_anime_this_season(cFilter, page, &count, &totalPages, &animeArray)
	if rc != 0 {
		return nil, utils.Pagination{}, domain_errors.AnimeFetchFailedError{}
	}

	if count == 0 {
		return []*domain.Anime{}, utils.Pagination{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			TotalPages: int(totalPages),
		}, nil
	}

	animeSlice := unsafe.Slice(animeArray, count)
	results := make([]*domain.Anime, count)

	for i := 0; i < int(count); i++ {
		a := &animeSlice[i]

		anime, err := r.animeMapper.CToGoPartial(unsafe.Pointer(a))

		if err != nil {
			C.free_partial_anime_array(animeArray, count)
			return nil, utils.Pagination{}, err
		}

		results[i] = anime
	}

	C.free_partial_anime_array(animeArray, count)
	return results, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: int(totalPages),
	}, nil
}

func (r *AnimeRepository) FetchStudioByID(studioID uint32, filters filters.AnimeFilter, pageNumber, pageSize int) (*value.Studio, []*domain.Anime, utils.Pagination, error) {
	cFilter, toFree := buildCFilter(filters)
	defer freeCFilter(toFree)

	var studioPtr C.studio_t
	var count C.uint
	var animeArray *C.partial_anime_t
	var totalPages C.uint
	var page = C.pageable_t{
		page_number: C.ushort(pageNumber),
		page_size:   C.ushort(pageSize),
	}

	rc := C.fetch_studio_by_id(C.uint(studioID), cFilter, &studioPtr, page, &count, &totalPages, &animeArray)
	if rc != 0 {
		return nil, nil, utils.Pagination{}, domain_errors.StudioNotFoundError{}
	}
	defer C.free_studio(&studioPtr)

	studio, err := r.animeMapper.CToGoStudio(unsafe.Pointer(&studioPtr))
	if err != nil {
		return nil, nil, utils.Pagination{}, err
	}

	if count == 0 {
		return studio, []*domain.Anime{}, utils.Pagination{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			TotalPages: int(totalPages),
		}, nil
	}

	animeSlice := unsafe.Slice(animeArray, count)
	results := make([]*domain.Anime, count)
	for i := 0; i < int(count); i++ {
		a := &animeSlice[i]
		anime, err := r.animeMapper.CToGoPartial(unsafe.Pointer(a))
		if err != nil {
			C.free_partial_anime_array(animeArray, count)
			return nil, nil, utils.Pagination{}, err
		}
		results[i] = anime
	}
	C.free_partial_anime_array(animeArray, count)

	return studio, results, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: int(totalPages),
	}, nil
}

func (r *AnimeRepository) FetchProducerByID(producerID uint32, filters filters.AnimeFilter, pageNumber, pageSize int) (*value.Producer, []*domain.Anime, utils.Pagination, error) {
	cFilter, toFree := buildCFilter(filters)
	defer freeCFilter(toFree)

	var producerPtr C.producer_t
	var count C.uint
	var animeArray *C.partial_anime_t
	var totalPages C.uint
	var page = C.pageable_t{
		page_number: C.ushort(pageNumber),
		page_size:   C.ushort(pageSize),
	}

	rc := C.fetch_producer_by_id(C.uint(producerID), cFilter, &producerPtr, page, &count, &totalPages, &animeArray)
	if rc != 0 {
		return nil, nil, utils.Pagination{}, domain_errors.ProducerNotFoundError{}
	}
	defer C.free_producer(&producerPtr)

	producer, err := r.animeMapper.CToGoProducer(unsafe.Pointer(&producerPtr))
	if err != nil {
		return nil, nil, utils.Pagination{}, err
	}

	if count == 0 {
		return producer, []*domain.Anime{}, utils.Pagination{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			TotalPages: int(totalPages),
		}, nil
	}

	animeSlice := unsafe.Slice(animeArray, count)
	results := make([]*domain.Anime, count)
	for i := 0; i < int(count); i++ {
		a := &animeSlice[i]
		anime, err := r.animeMapper.CToGoPartial(unsafe.Pointer(a))
		if err != nil {
			C.free_partial_anime_array(animeArray, count)
			return nil, nil, utils.Pagination{}, err
		}
		results[i] = anime
	}
	C.free_partial_anime_array(animeArray, count)

	return producer, results, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: int(totalPages),
	}, nil
}

func (r *AnimeRepository) FetchLicensorByID(licensorID uint32, filters filters.AnimeFilter, pageNumber, pageSize int) (*value.Licensor, []*domain.Anime, utils.Pagination, error) {
	cFilter, toFree := buildCFilter(filters)
	defer freeCFilter(toFree)

	var licensorPtr C.licensor_t
	var count C.uint
	var animeArray *C.partial_anime_t
	var totalPages C.uint
	var page = C.pageable_t{
		page_number: C.ushort(pageNumber),
		page_size:   C.ushort(pageSize),
	}

	rc := C.fetch_licensor_by_id(C.uint(licensorID), cFilter, &licensorPtr, page, &count, &totalPages, &animeArray)
	if rc != 0 {
		return nil, nil, utils.Pagination{}, domain_errors.LicensorNotFoundError{}
	}
	defer C.free_licensor(&licensorPtr)

	licensor, err := r.animeMapper.CToGoLicensor(unsafe.Pointer(&licensorPtr))
	if err != nil {
		return nil, nil, utils.Pagination{}, err
	}

	if count == 0 {
		return licensor, []*domain.Anime{}, utils.Pagination{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			TotalPages: int(totalPages),
		}, nil
	}

	animeSlice := unsafe.Slice(animeArray, count)
	results := make([]*domain.Anime, count)
	for i := 0; i < int(count); i++ {
		a := &animeSlice[i]
		anime, err := r.animeMapper.CToGoPartial(unsafe.Pointer(a))
		if err != nil {
			C.free_partial_anime_array(animeArray, count)
			return nil, nil, utils.Pagination{}, err
		}
		results[i] = anime
	}
	C.free_partial_anime_array(animeArray, count)

	return licensor, results, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: int(totalPages),
	}, nil
}

func (r *AnimeRepository) FetchAnimeFromTag(tagID uint32, f filters.AnimeFilter, pageNumber, pageSize int) ([]*domain.Anime, utils.Pagination, error) {
	cFilter, toFree := buildCFilter(f)
	defer freeCFilter(toFree)

	var count C.uint
	var animeArray *C.partial_anime_t
	var totalPages C.uint
	var page = C.pageable_t{
		page_number: C.ushort(pageNumber),
		page_size:   C.ushort(pageSize),
	}

	rc := C.fetch_anime_from_tag(C.uint(tagID), cFilter, page, &count, &totalPages, &animeArray)
	if rc != 0 {
		return nil, utils.Pagination{}, domain_errors.AnimeFetchFailedError{}
	}

	if count == 0 {
		return []*domain.Anime{}, utils.Pagination{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			TotalPages: int(totalPages),
		}, nil
	}

	animeSlice := unsafe.Slice(animeArray, count)
	results := make([]*domain.Anime, count)
	for i := 0; i < int(count); i++ {
		a := &animeSlice[i]
		anime, err := r.animeMapper.CToGoPartial(unsafe.Pointer(a))
		if err != nil {
			C.free_partial_anime_array(animeArray, count)
			return nil, utils.Pagination{}, err
		}
		results[i] = anime
	}
	C.free_partial_anime_array(animeArray, count)

	return results, utils.Pagination{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: int(totalPages),
	}, nil
}
