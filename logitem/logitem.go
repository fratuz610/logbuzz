package logitem

import (
	"container/list"
	"strings"
	"sync"
	"time"
)

type LogItem struct {
	InputID   string   `json:"inputID"`
	Timestamp int64    `json:"timestamp"`
	Level     string   `json:"level"`
	Message   string   `json:"message"`
	Env       string   `json:"env"`
	TagList   []string `json:"tagList"`
}

func NewLogItem(inputID, level, message, env string) *LogItem {
	ret := LogItem{}
	ret.InputID = inputID
	ret.Level = level
	ret.Message = message
	ret.Env = env
	ret.Timestamp = time.Now().Unix()
	ret.TagList = []string{}
	return &ret
}

var log *list.List = list.New()
var logMutex sync.Mutex

func AddLogItem(item *LogItem) {
	logMutex.Lock()
	defer logMutex.Unlock()

	log.PushFront(item)
}

func Search(inputID, level, env string, fromTS, toTS int64, skip int) []LogItem {

	logMutex.Lock()
	defer logMutex.Unlock()

	// 500 is the maximum number of items returned per query
	retList := make([]LogItem, 0, 500)

	addedSoFar := 0

	for cursor := log.Front(); cursor != nil; cursor = cursor.Next() {

		logItem := cursor.Value.(*LogItem)

		if strings.ToLower(inputID) != logItem.InputID {
			continue
		}

		if level != "" && strings.ToLower(level) != logItem.Level {
			continue
		}

		if env != "" && strings.ToLower(env) != logItem.Env {
			continue
		}

		if toTS != 0 && logItem.Timestamp > toTS {
			continue
		}

		if fromTS != 0 && logItem.Timestamp < fromTS {
			break
		}

		// we skip if needed
		if skip > 0 {
			skip--
			continue
		}

		retList = append(retList, *logItem)
		addedSoFar++

		if addedSoFar >= 500 {
			break
		}
	}

	return retList
}

func TrimEnd() {
	logMutex.Lock()
	defer logMutex.Unlock()

	var removeCount int = 5 * log.Len() / 100

	//golog.Printf("About to remove %d items", removeCount)

	for {
		item := log.Back()
		if item == nil {
			return
		}

		log.Remove(item)

		removeCount--

		if removeCount <= 0 {
			return
		}
	}

}

func GetNumItems() int {
	logMutex.Lock()
	defer logMutex.Unlock()
	return log.Len()
}

func GetLatestTimestamp() int64 {
	logMutex.Lock()
	defer logMutex.Unlock()

	if log.Front() == nil {
		return 1
	}

	logItem := log.Front().Value.(*LogItem)

	return logItem.Timestamp
}
