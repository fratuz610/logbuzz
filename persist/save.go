package persist

import (
	"encoding/gob"
	golog "log"
	"logbuzz/data"
	"logbuzz/list"
	"os"
)

const (
	PERSIST_FILE      string = "logbuzz.bin"
	TEMP_PERSIST_FILE string = "logbuzz_old.bin"
)

var saveChannel chan data.LogItem = make(chan data.LogItem)
var quit chan int = make(chan int)

func Persist(item *data.LogItem) {
	saveChannel <- *item
}

func PersistanceLoop() {

	golog.Println("Starting persistance cycle")

	// we open the target file for writing and keep it open
	file, err := os.OpenFile(PERSIST_FILE, os.O_APPEND, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(PERSIST_FILE)
			if err != nil {
				golog.Fatal("Unable to create file:", PERSIST_FILE, " because ", err)
			}
		} else {
			golog.Fatal("Unable to open file:", PERSIST_FILE, " because ", err)
		}
	}

	// we create an encoder
	encoder := gob.NewEncoder(file)

	for {

		select {
		case item := <-saveChannel:
			err = encoder.Encode(item)
			if err != nil {
				golog.Fatal("Unable to encode logitem because: '", err, "'")
			}
		case <-quit:
			break
		}
	}

	golog.Println("Ending persistance cycle")

	file.Close()
}

func StopPersistanceLoop() {
	quit <- 0
}

func RenamePersistanceLog() {
	// we remove the old TEMP persistence file
	os.Remove(TEMP_PERSIST_FILE)

	err := os.Rename(PERSIST_FILE, TEMP_PERSIST_FILE)
	if err != nil {
		golog.Println("Unable to rename: ", PERSIST_FILE, " to ", TEMP_PERSIST_FILE, " because: ", err)
		return
	}
}

func LoadSavedLogs() {

	// we open the target file for reading and keep it open
	file, err := os.Open(TEMP_PERSIST_FILE)
	if err != nil {
		golog.Println("Unable to open log file:", TEMP_PERSIST_FILE, " because ", err, "file: ", file)
		return
	}

	// we close the file
	defer func() {
		if file != nil {
			file.Close()
		}
		err := os.Remove(TEMP_PERSIST_FILE)
		if err != nil {
			golog.Println("Unable to remove log file:", TEMP_PERSIST_FILE, " because ", err)
		}
	}()

	decoder := gob.NewDecoder(file)
	count := 0

	for {
		logItem := data.LogItem{}
		err := decoder.Decode(&logItem)

		if err != nil {
			break
		}

		list.AddLogItem(&logItem)
		Persist(&logItem)

		count++
	}

	golog.Println("Loading complete: ", count, " items loaded")

}
