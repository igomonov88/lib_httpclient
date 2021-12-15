package lib_httpclient

import (
	"bytes"
	"io"
	"io/ioutil"
	"sync"
)

// RestorableReadCloser allows to read data once and then restore it from backup if you need to re-read it.
type RestorableReadCloser struct {
	rs   io.ReadSeeker
	lock sync.Mutex
}

// NewRestorableReadCloser creates new RestorableReadCloser.
func NewRestorableReadCloser(r io.Reader) (*RestorableReadCloser, error) {
	rc := &RestorableReadCloser{}

	var backup []byte

	if r != nil {
		var err error
		backup, err = ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
	}

	rc.rs = bytes.NewReader(backup)

	return rc, nil
}

// Read reads bytes from internal buffer.
func (rc *RestorableReadCloser) Read(p []byte) (int, error) {
	rc.lock.Lock()
	defer rc.lock.Unlock()

	return rc.rs.Read(p)
}

// Close just supports io.ReadCloser interface, but does nothing.
func (rc *RestorableReadCloser) Close() error {
	rc.lock.Lock()
	defer rc.lock.Unlock()

	return nil
}

// Restore internal buffer and return itself. It is ready to be read again.
func (rc *RestorableReadCloser) Restore() *RestorableReadCloser {
	rc.lock.Lock()
	defer rc.lock.Unlock()

	rc.rs.Seek(0, io.SeekStart)

	return rc
}
