// logbuzz project main.go
package main

import (
	hosieweb "github.com/hoisie/web"
	"logbuzz/persist"
	myweb "logbuzz/web"
)

func main() {

	// first we rename the persistance file
	persist.RenamePersistanceLog()

	// then we start the load routine
	go persist.LoadSavedLogs()

	// then we start the persistance loop (shared among web and internal calls)
	go persist.PersistanceLoop()

	// we setup the web interface
	hosieweb.Post("/api/input", myweb.PostInput)
	hosieweb.Get("/api/query", myweb.GetQuery)
	hosieweb.Get("/api/stats", myweb.GetStats)
	hosieweb.Run("0.0.0.0:1212")

}

/*
func rotateLogs() {

	ticker := time.NewTicker(time.Second * 10)

	for _ = range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		//log.Printf("Heap Sys %d, Heap Alloc %d, Heap Idle %d, Heap Relased %d\n", m.HeapSys, m.HeapAlloc, m.HeapIdle, m.HeapReleased)
		if m.HeapAlloc > maxMemory { // 10 mb
			log.Printf("The heap size is %d. trimming and calling GC\n", m.HeapAlloc)
			list.TrimEnd()
			debug.FreeOSMemory()
		}

	}

}
*/
