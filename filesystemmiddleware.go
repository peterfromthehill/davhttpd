package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Middleware interface {
	BeforeOpen(http.FileSystem, string)
	AfterOpen(http.FileSystem, string, http.File, error)
}

type FileSystemMiddleware struct {
	FileSystem  http.FileSystem
	Middlewares []Middleware
}

func (f FileSystemMiddleware) Open(name string) (http.File, error) {
	for _, middleware := range f.Middlewares {
		middleware.BeforeOpen(f, name)
	}
	file, err := f.FileSystem.Open(name)
	if err != nil {
		log.Debugf("Error: %s", err)
		return file, err
	}
	for _, middleware := range f.Middlewares {
		middleware.AfterOpen(f, name, file, err)
	}
	return file, nil
}
