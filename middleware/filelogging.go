package middleware

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type FileLogging struct {
	next   http.FileSystem
	Prefix string
}

func (fl *FileLogging) SetNext(middleware http.FileSystem) {
	fl.next = middleware
}

func (fl *FileLogging) Open(name string) (http.File, error) {
	log.Println(fl.Prefix, "FileLogging.Name:", name)
	if fl.next == nil {
		panic("is nil!")
	}
	a, b := fl.next.Open(name)
	if b != nil {
		log.Println(fl.Prefix, "FileLogging.Error:", b)
		return a, b
	}
	c, _ := a.Stat()
	log.Println(fl.Prefix, "FileLogging.Stats:", c)
	return a, b
}
