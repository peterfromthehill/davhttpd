package main

import (
	"errors"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type Reader struct {
	driver    *driver
	path      string
	fileInfo  os.FileInfo
	davReader io.ReadCloser
	i         int64 // current reading index
}

func NewWebDavReader(path string, d *driver, fileInfo os.FileInfo, reader io.ReadCloser) *Reader {
	w := &Reader{
		driver:    d,
		path:      path,
		fileInfo:  fileInfo,
		davReader: reader,
	}
	return w
}

// Len returns the number of bytes of the unread portion of the
// slice.
func (r *Reader) Len() int {
	log.Debugf("Len: %d", int(r.fileInfo.Size()-r.i))
	if r.i >= r.fileInfo.Size() {
		return 0
	}
	return int(r.fileInfo.Size() - r.i)
}

// Size returns the original length of the underlying byte slice.
// Size is the number of bytes available for reading via ReadAt.
// The result is unaffected by any method calls except Reset.
func (r *Reader) Size() int64 { return r.fileInfo.Size() }

// Read implements the io.Reader interface.
func (r *Reader) Read(b []byte) (n int, err error) {
	if r.i >= r.fileInfo.Size() {
		return 0, io.EOF
	}
	if r.davReader != nil {
		n, err := r.davReader.Read(b)
		r.i += int64(n)
		return n, err
	}
	reader, err := r.driver.c.ReadStreamRange(r.path, r.i, r.fileInfo.Size()-r.i)
	if err != nil {
		return 0, err
	}
	defer reader.Close()
	n, err = reader.Read(b)
	if err != nil {
		return n, err
	}
	r.i += int64(n)
	log.Debugf("Read %d bytes", n)
	return
}

// Seek implements the io.Seeker interface.
func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = r.fileInfo.Size() + offset
	default:
		return 0, errors.New("bytes.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("bytes.Reader.Seek: negative position")
	}
	r.i = abs
	var err error
	r.davReader, err = r.driver.c.ReadStreamRange(r.path, abs, r.Size())
	if err != nil {
		return 0, errors.New("cannot Seek because of webdav api")
	}
	return abs, nil
}
