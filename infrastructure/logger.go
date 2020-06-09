package infrastructure

import (
	"log"
)

//Logger ...
type Logger struct{}

//Log log messages to console
func (logger Logger) Log(args ...interface{}) {
	log.Println(args...)
}
