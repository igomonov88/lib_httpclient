package lib_httpclient

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestorableReadCloser(t *testing.T) {
	body := bytes.NewBufferString("content")

	rc, err := NewRestorableReadCloser(body)
	assert.Nil(t, err)

	readed1, err1 := ioutil.ReadAll(rc)
	assert.Nil(t, err1)
	assert.Equal(t, "content", string(readed1))

	rc.Close()

	readed2, err2 := ioutil.ReadAll(rc.Restore())
	assert.Nil(t, err2)
	assert.Equal(t, "content", string(readed2))
}

func TestRestorableReadCloser_NilReader(t *testing.T) {
	bk, err := NewRestorableReadCloser(nil)
	assert.Nil(t, err)

	readed, err := ioutil.ReadAll(bk)
	assert.Nil(t, err)
	assert.Equal(t, "", string(readed))
}
