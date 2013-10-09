// logbuzz project main.go
package main

import (
	hosieweb "github.com/hoisie/web"
	"log"
	"logbuzz/logitem"
	myweb "logbuzz/web"
	"runtime"
	"runtime/debug"
	"time"
)

const maxMemory uint64 = 20 * 1024 * 1024

func main() {
	go rotateLogs()
	hosieweb.Post("/api/input", myweb.PostInput)
	hosieweb.Get("/api/query", myweb.GetQuery)
	hosieweb.Get("/api/stats", myweb.GetStats)
	hosieweb.Run("0.0.0.0:1212")
}

func rotateLogs() {

	ticker := time.NewTicker(time.Second * 10)

	for _ = range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		//log.Printf("Heap Sys %d, Heap Alloc %d, Heap Idle %d, Heap Relased %d\n", m.HeapSys, m.HeapAlloc, m.HeapIdle, m.HeapReleased)
		if m.HeapAlloc > maxMemory { // 10 mb
			log.Printf("The heap size is %d. trimming and calling GC\n", m.HeapAlloc)
			logitem.TrimEnd()
			debug.FreeOSMemory()
		}

	}

}
