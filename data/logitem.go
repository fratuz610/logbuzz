package data

import (
	"time"
)

type LogItem struct {
	Timestamp int64    `json:"timestamp"`
	Level     string   `json:"level"`
	Message   string   `json:"message"`
	TagList   []string `json:"tagList"`
}

func NewLogItem(level, message string, tagList []string) *LogItem {
	ret := LogItem{}
	ret.Timestamp = time.Now().Unix()
	ret.Level = level
	ret.Message = message
	ret.TagList = tagList

	if ret.TagList == nil {
		ret.TagList = make([]string, 0)
	}

	return &ret
}
