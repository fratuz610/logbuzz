package logitem

import (
	"fmt"
)

struct LogItem {
	InputID string `json:"inputID"`
	Timestamp uint64
	Level string
	Message string
	Env string
}