package intlog

import (
	"logbuzz/data"
	"logbuzz/list"
	"strings"
)

type InternalLogWriter struct {
	TagList []string
}

func (this *InternalLogWriter) Write(p []byte) (n int, err error) {
	list.AddLogItem(data.NewLogItem("INFO", strings.Trim(string(p), "\n\r\t"), this.TagList))
	return len(p), nil
}
