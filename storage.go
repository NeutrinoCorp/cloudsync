package cloudsync

import (
	"context"
	"io"
	"time"
)

type ReadSeekerAt interface {
	io.ReadSeeker
	io.ReaderAt
	io.Closer
}

type File struct {
	Key  string
	Data ReadSeekerAt
}

type BlobStorage interface {
	Upload(ctx context.Context, f File) error
	CheckMod(ctx context.Context, key string, modTime time.Time, size int64) (bool, error)
}
