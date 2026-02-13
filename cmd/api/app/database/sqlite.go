package database

/*
#cgo linux LDFLAGS: -L${SRCDIR}/../../../../drivers -Wl,-rpath,${SRCDIR}/../../../../drivers -lanime_facts
#cgo windows LDFLAGS: -L${SRCDIR}/../../../../drivers -lanime_facts
#cgo CFLAGS: -I${SRCDIR}/../../../../../anime-facts-core/include

#include "anime_facts_api.h"
#include <stdlib.h>
*/
import "C"

import (
	"log"

	"github.com/afuradanime/backend/config"
)

func InitSQLite(Config config.Config) {
	// Set anime database path
	C.set_database_path(C.CString(Config.AnimeDatabasePath))
	log.Println("Anime database path set to ./../anime.db")
}
