package persist

import (
	"logbuzz/data"
)

type WriteOp int

const (
	ADD WriteOp = iota
	DELETE
)

type LogOperation struct {
	LogItem *data.LogItem `json:"logItem"`
	WriteOp WriteOp       `json:"writeOp"`
}

func NewAddOperation(logItem *data.LogItem) LogOperation {
	return LogOperation{LogItem: logItem, WriteOp: ADD}
}

func NewDeleteOperation() LogOperation {
	return LogOperation{LogItem: nil, WriteOp: DELETE}
}
