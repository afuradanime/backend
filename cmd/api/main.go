package main

/*
#cgo linux LDFLAGS: -L${SRCDIR}/../../drivers -Wl,-rpath,${SRCDIR}/../../../drivers -lanime_facts
#cgo windows LDFLAGS: -L${SRCDIR}/../../drivers -lanime_facts
#cgo CFLAGS: -I${SRCDIR}/../../../anime-facts-core/include

#include "anime_facts_api.h"
#include <stdlib.h>
*/
import "C"

import "github.com/afuradanime/backend/cmd/api/app"

func main() {

	// Set database path
	C.set_database_path(C.CString("./../anime.db"))

	app.New()
}
