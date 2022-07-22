package utils

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

func GetCurrentTimestampSec() int {
	ts, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))
	return ts
}

func NewUUID() (string, error) {
	uuid := uuid.New()
	key := uuid.String()
	return key, nil
}
