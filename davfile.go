package main

import (
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type DavFile struct {
	http.File
	path     string
	root     DavFileSystem
	fileInfo os.FileInfo
}

type DavFileInstance struct {
	root      DavFileSystem
	fileInfo  os.FileInfo
	name      string
	davReader io.ReadCloser
	*Reader
}

func NewDavFile(root DavFileSystem, name string, fileInfo os.FileInfo) *DavFile {
	return &DavFile{
		path:     name,
		root:     root,
		fileInfo: fileInfo,
	}
}

func (f DavFile) OpenAs(name string) (http.File, error) {
	var davReader io.ReadCloser
	var err error
	if !f.fileInfo.IsDir() {
		davReader, err = f.root.driverImpl.c.ReadStream(name)
		if err != nil {
			return nil, err
		}
	}
	reader := NewWebDavReader(name, f.root.driverImpl, f.fileInfo, davReader)
	davFileInstance := DavFileInstance{
		name:      name,
		root:      f.root,
		fileInfo:  f.fileInfo,
		Reader:    reader,
		davReader: davReader,
	}

	return davFileInstance, nil
}

func (d DavFile) Stat() (os.FileInfo, error) {
	return d.fileInfo, nil
}

func (d DavFileInstance) Readdir(count int) ([]os.FileInfo, error) {
	log.Debugf("Readdir()")
	dirInfo, err := d.root.driverImpl.c.ReadDir(d.name)
	return dirInfo, err
}

func (d DavFileInstance) Close() error {
	log.Debugf("Close()")
	if d.davReader != nil {
		d.davReader.Close()
	}
	return nil
}
func (d DavFileInstance) Stat() (os.FileInfo, error) {
	log.Debugf("Stat")
	return d.fileInfo, nil
}
func (d DavFileInstance) Name() string {
	log.Debugf("Name")
	return d.name
}
func (d DavFileInstance) Size() int64 {
	log.Debugf("Size")
	return d.fileInfo.Size()
}
func (d DavFileInstance) Mode() os.FileMode {
	log.Debugf("Mode")
	return d.fileInfo.Mode()
}
func (d DavFileInstance) ModTime() time.Time {
	log.Debugf("Modtime")
	return d.fileInfo.ModTime()
}
func (d DavFileInstance) IsDir() bool {
	log.Debugf("IsDir")
	return d.fileInfo.IsDir()
}
func (d DavFileInstance) Sys() interface{} {
	log.Debugf("Sys?")
	return d.fileInfo.Sys()
}
