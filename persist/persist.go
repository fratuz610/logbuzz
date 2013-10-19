package persist

import (
	"compress/flate"
	"encoding/json"
	golog "log"
	"logbuzz/data"
	"logbuzz/list"
	"logbuzz/vars"
	"os"
)

var logChannel chan data.LogItem = make(chan data.LogItem)
var quit chan int = make(chan int)

func PersistLogItem(item *data.LogItem) {
	logChannel <- *item
}

func PersistanceLoop() {

	golog.Println("Starting persistance cycle")

	// we open the target file for writing and keep it open
	file, err := os.OpenFile(vars.GetPersistFile(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(vars.GetPersistFile())
			if err != nil {
				golog.Fatal("Unable to create file:", vars.GetPersistFile(), " because ", err)
			}
		} else {
			golog.Fatal("Unable to open file:", vars.GetPersistFile(), " because ", err)
		}
	}

	defer file.Close()

	// we open the flateWriter stream
	flateWriter, err := flate.NewWriter(file, 4)

	if err != nil {
		golog.Fatal("Unable to create flate writer:", err)
	}

	defer flateWriter.Close()

	// we create an encoder
	encoder := json.NewEncoder(flateWriter)

	for {

		select {
		case op := <-logChannel:
			err = encoder.Encode(op)
			flateWriter.Flush()
			if err != nil {
				golog.Fatal("Unable to encode log item because: '", err, "'")
			}
		case <-quit:
			break
		}
	}

	golog.Println("Ending persistance cycle")

}

func StopPersistanceLoop() {
	quit <- 0
}

func LoadSavedLogs() {

	// we open the target file for reading and keep it open
	file, err := os.Open(vars.GetPersistFile())
	if err != nil {
		golog.Println("Unable to open log file:", vars.GetPersistFile(), "because", err)
		return
	}

	// we close the file on defer
	defer file.Close()

	flateReader := flate.NewReader(file)

	defer flateReader.Close()

	decoder := json.NewDecoder(flateReader)
	count := 0

	for {
		logItem := data.LogItem{}
		err := decoder.Decode(&logItem)

		if err != nil {
			break
		}

		list.AddLogItem(&logItem)

		count++
	}

	golog.Println("Loading complete: ", count, " items loaded")

}
