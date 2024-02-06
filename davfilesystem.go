package main

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/studio-b12/gowebdav"
)

type driver struct {
	c *gowebdav.Client
}

type DavFileSystem struct {
	driverImpl *driver
}

func NewDavFileSystem(url, username, password string) *DavFileSystem {
	return &DavFileSystem{
		driverImpl: &driver{
			c: gowebdav.NewClient(url, username, password),
		},
	}
}

func (d DavFileSystem) Open(name string) (http.File, error) {
	name = path.Clean(name)
	info, err := d.Stat(name)
	if err != nil {
		if strings.HasSuffix(name, "index.html") {
			return nil, err
		}
		log.Debugf("Name: %s", name)
		info, err = d.Stat(name + "/")
		if err != nil {
			log.Debugf("Error: %v", err)
			return nil, fs.ErrNotExist
		}
		name = name + "/"

	}
	log.Debugf("Info: %v", info)

	file := NewDavFile(d, name, info)
	httpFile, err := file.OpenAs(name)
	if err != nil {
		log.Debugf("Error: %v", err)
	}
	return httpFile, err
}

func (d DavFileSystem) Stat(name string) (os.FileInfo, error) {
	fileInfo, err := d.driverImpl.c.Stat(name)
	return fileInfo, err
}
