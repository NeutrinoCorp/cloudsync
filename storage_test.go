package cloudsync_test

import (
	"errors"
	"testing"
	"time"

	"github.com/neutrinocorp/cloudsync"
	"github.com/stretchr/testify/assert"
)

func TestNoopBlobStorage(t *testing.T) {
	storage := cloudsync.NoopBlobStorage{
		UploadErr:    errors.New("foo err"),
		CheckModBool: true,
		CheckModErr:  errors.New("bar err"),
	}
	assert.Equal(t, "foo err", storage.Upload(nil, cloudsync.Object{}).Error())
	b, err := storage.CheckMod(nil, "", time.Time{}, 0)
	assert.Equal(t, "bar err", err.Error())
	assert.True(t, b)
}
