package middleware

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

var cacheError map[string]error

type FileCache struct {
	next   http.FileSystem
	Prefix string
}

func (fl *FileCache) SetNext(middleware http.FileSystem) {
	fl.next = middleware
}

func (fl *FileCache) Open(name string) (http.File, error) {
	if cacheError == nil {
		cacheError = make(map[string]error)
	}
	err, ok := cacheError[name]
	if ok {
		log.Println(fl.Prefix, "FileCache.Return from cache:", name, err)
		return nil, err
	}
	if fl.next == nil {
		panic("is nil!")
	}
	a, b := fl.next.Open(name)
	if b != nil {
		log.Println(fl.Prefix, "FileCache.Error:", b)
		cacheError[name] = b
		return a, b
	}
	return a, b
}
