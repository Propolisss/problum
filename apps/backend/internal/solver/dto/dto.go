package dto

import "time"

type Result struct {
	Duration     time.Duration `db:"duration"`
	MemoryUsage  int64         `db:"memory_usage"`
	Status       string        `db:"status"`
	ErrorMessage *string       `db:"error_message"`
}

type Limits struct {
	TimeLimit   time.Duration
	MemoryLimit int64
}
