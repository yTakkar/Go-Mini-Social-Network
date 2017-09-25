package config

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// MakeTimestamp function
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// Err Log
func Err(err interface{}) {
	if err != nil {
		log.Fatal(err)
	}
}

// MeOrNot function to checked whether it's me or not
func MeOrNot(c *gin.Context, user string) bool {
	id, _ := AllSessions(c)
	if id != user {
		return false
	}
	return true
}

