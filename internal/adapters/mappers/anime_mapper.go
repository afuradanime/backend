package mappers

/*
#cgo linux LDFLAGS: -L${SRCDIR}/../../../drivers -Wl,-rpath,${SRCDIR}/../../../drivers -lanime_facts
#cgo windows LDFLAGS: -L${SRCDIR}/../../../drivers -lanime_facts
#cgo CFLAGS: -I${SRCDIR}/../../../../anime-facts-core/include

#include "anime_facts_api.h"
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/afuradanime/backend/internal/core/domain"
)

type AnimeMapper struct{}

func NewAnimeMapper() *AnimeMapper {
	return &AnimeMapper{}
}

func (m *AnimeMapper) CtoGo(animePtr unsafe.Pointer) (*domain.Anime, error) {
	// Cast unsafe.Pointer back to C.anime_t pointer
	cAnime := (*C.anime_t)(animePtr)

	// save the c information in a go struct
	anime, err := domain.NewAnime(
		uint32(cAnime.id),
		C.GoString(cAnime.url),
		C.GoString(cAnime.title),
		domain.AnimeType(cAnime._type),
		C.GoString(cAnime.source),
		uint32(cAnime.episodes),
		domain.AnimeStatus(cAnime.status),
		bool(cAnime.airing),
		C.GoString(cAnime.duration),
		C.GoString(cAnime.start_date),
		C.GoString(cAnime.end_date),
		domain.SeasonType(cAnime.season.season),
		uint16(cAnime.season.year),
		C.GoString(cAnime.broadcast.day),
		C.GoString(cAnime.broadcast.time),
		C.GoString(cAnime.broadcast.timezone),
		C.GoString(cAnime.image_url),
		C.GoString(cAnime.small_image_url),
		C.GoString(cAnime.large_image_url),
		C.GoString(cAnime.trailer_embed_url),
	)

	if err != nil {
		return nil, err
	}

	// Fill synonyms
	if cAnime.synonyms.count > 0 {
		synonymsSlice := unsafe.Slice(cAnime.synonyms.items, cAnime.synonyms.count)
		for _, synonymPtr := range synonymsSlice {
			if synonymPtr != nil {
				anime.AddSynonym(C.GoString(synonymPtr))
			}
		}
	}

	// Fill descriptions
	if cAnime.descriptions.count > 0 {
		descriptionsSlice := unsafe.Slice(cAnime.descriptions.items, cAnime.descriptions.count)
		for _, descPtr := range descriptionsSlice {
			if descPtr.description != nil {
				desc := domain.Description{
					Language:    domain.Language(descPtr.language),
					Description: C.GoString(descPtr.description),
				}
				anime.AddDescription(desc)
			}
		}
	}

	// Fill tags
	if cAnime.tags.count > 0 {
		tagsSlice := unsafe.Slice(cAnime.tags.items, cAnime.tags.count)
		for _, tagPtr := range tagsSlice {
			tag := domain.Tag{
				ID:   uint32(tagPtr.id),
				Name: C.GoString(tagPtr.name),
				Type: domain.TagType(tagPtr._type),
				URL:  C.GoString(tagPtr.url),
			}
			anime.AddTag(tag)
		}
	}

	// Fill producers
	if cAnime.producers.count > 0 {
		producersSlice := unsafe.Slice(cAnime.producers.items, cAnime.producers.count)
		for _, producerPtr := range producersSlice {
			producer := domain.Producer{
				ID:   uint32(producerPtr.id),
				Name: C.GoString(producerPtr.name),
				Type: C.GoString(producerPtr._type),
				URL:  C.GoString(producerPtr.url),
			}
			anime.AddProducer(producer)
		}
	}

	// Fill licensors
	if cAnime.licensors.count > 0 {
		licensorsSlice := unsafe.Slice(cAnime.licensors.items, cAnime.licensors.count)
		for _, licensorPtr := range licensorsSlice {
			licensor := domain.Licensor{
				ID:   uint32(licensorPtr.id),
				Name: C.GoString(licensorPtr.name),
				Type: C.GoString(licensorPtr._type),
				URL:  C.GoString(licensorPtr.url),
			}
			anime.AddLicensor(licensor)
		}
	}

	// Fill studios
	if cAnime.studios.count > 0 {
		studiosSlice := unsafe.Slice(cAnime.studios.items, cAnime.studios.count)
		for _, studioPtr := range studiosSlice {
			studio := domain.Studio{
				ID:   uint32(studioPtr.id),
				Name: C.GoString(studioPtr.name),
				URL:  C.GoString(studioPtr.url),
			}
			anime.AddStudio(studio)
		}
	}

	return anime, nil
}
