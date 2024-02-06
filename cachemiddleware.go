package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type RedisCache string

func (rc RedisCache) BeforeOpen(fileSystem http.FileSystem, name string) {
	log.Debugf("Open fs: %s file: %s", fileSystem, name)
}

func (rc RedisCache) AfterOpen(fileSystem http.FileSystem, name string, file http.File, err error) {
	//fileInfo, err := file.Stat()
	//log.Printf("Add Cache: FS %s / File: %v / Fileinfo: %v / Err: %s", fileSystem, file, fileInfo, err)
}
