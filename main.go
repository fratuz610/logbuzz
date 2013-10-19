// logbuzz project main.go
package main

import (
	hosieweb "github.com/hoisie/web"
	//"log"
	//"logbuzz/intlog"
	"logbuzz/config"
	"logbuzz/persist"
	myweb "logbuzz/web"
	"runtime/debug"
	"strconv"
	"time"
)

func main() {

	// we load the redo log (synchronous)
	persist.LoadSavedLogs()

	// then we start the persistance loop
	go persist.PersistanceLoop()

	//log.SetOutput(&intlog.InternalLogWriter{[]string{"internal-log", "persistance"}})
	//hosieweb.SetLogger(log.New(&intlog.InternalLogWriter{[]string{"internal-log", "web"}}, "", 0))

	// we start the memory release routine
	go releaseMemory()

	// we setup the web interface
	hosieweb.Post("/api/input", myweb.PostInput)
	hosieweb.Get("/api/query", myweb.GetQuery)
	hosieweb.Get("/api/stats", myweb.GetStats)
	hosieweb.Get("/api/config", myweb.GetConfig)

	hosieweb.Post("/api/config/httpPort", myweb.PostConfigHttpPort)
	hosieweb.Post("/api/config/memoryLimit", myweb.PostConfigMemoryLimit)

	for {
		hosieweb.Run("0.0.0.0:" + strconv.Itoa(config.GetConfig().HttpPort))
	}

}

func releaseMemory() {

	ticker := time.NewTicker(time.Minute * 1)

	for _ = range ticker.C {
		debug.FreeOSMemory()
	}

}
