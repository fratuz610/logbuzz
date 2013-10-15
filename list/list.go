package list

import (
	"container/list"
	"logbuzz/data"
	"runtime"
	"strings"
	"sync"
	"time"
)

const MAX_MEMORY uint64 = 20 * 1024 * 1024

var log *list.List = list.New()
var logMutex sync.Mutex

func AddLogItem(item *data.LogItem) {
	logMutex.Lock()
	defer logMutex.Unlock()

	// we update the TS (internally generated)
	item.Timestamp = time.Now().Unix()

	// we save in front
	log.PushFront(item)

	// we check the memory usage
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// if we are over the limit, we remove an old object
	if m.HeapAlloc > MAX_MEMORY {
		item := log.Back()
		log.Remove(item)
	}

}

func Search(level string, tagList []string, fromTS, toTS int64, skip int) []*data.LogItem {

	logMutex.Lock()
	defer logMutex.Unlock()

	// 500 is the maximum number of items returned per query
	retList := make([]*data.LogItem, 0, 500)

	addedSoFar := 0

	for cursor := log.Front(); cursor != nil; cursor = cursor.Next() {

		logItem := cursor.Value.(*data.LogItem)

		// we skip if needed
		if skip > 0 {
			skip--
			continue
		}

		if level != "" && strings.ToLower(level) != logItem.Level {
			continue
		}

		if tagList != nil {
			found := false

			for _, tag := range tagList {
				for _, savedTag := range logItem.TagList {
					if tag == savedTag {
						found = true
					}
				}
			}

			if !found {
				continue
			}
		}

		if toTS != 0 && logItem.Timestamp > toTS {
			continue
		}

		if fromTS != 0 && logItem.Timestamp < fromTS {
			continue
		}

		retList = append(retList, logItem)
		addedSoFar++

		if addedSoFar >= 500 {
			break
		}
	}

	return retList
}

func GetNumItems() int {
	logMutex.Lock()
	defer logMutex.Unlock()
	return log.Len()
}

func GetTimestampRange() (int64, int64) {
	logMutex.Lock()
	defer logMutex.Unlock()

	var startTS int64 = 1
	var endTS int64 = 1
	if log.Front() != nil {
		endTS = log.Front().Value.(*data.LogItem).Timestamp
	}
	if log.Back() != nil {
		startTS = log.Back().Value.(*data.LogItem).Timestamp
	}

	return startTS, endTS
}
