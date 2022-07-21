package cloudsync

import (
	"context"
	"io"
	"time"
)

// ReadSeekerAt a custom read buffer type used to perform memory-efficient allocations when
// reading large objects.
//
// In deepness, io.Reader requires a complete buffer allocation when reading a slice of bytes while the
// combination of io.ReadSeeker and io.ReaderAt allows to read specific parts of the given slice of bytes
// (avoiding unnecessary memory allocations, i.e. full buffer allocation) while still satisfying the io.Reader interface.
//
// Finally, this increases application performance drastically when reading a big slice of bytes (i.e. a large
// PDF or docx file) as underlying upload APIs from third party vendors might partition these files using
// a multipart strategy.
//
// For more information, please read: https://aws.github.io/aws-sdk-go-v2/docs/sdk-utilities/s3/.
type ReadSeekerAt interface {
	io.ReadSeeker
	io.ReaderAt
}

// Object also known as file, information unit stored within a directory composed of an io.Reader holding
// binary data (Data) and a Key.
type Object struct {
	// Key file's path + name or name.
	Key string
	// Data Binary Large Object reader instance.
	Data ReadSeekerAt
	// CleanupFunc frees resources like underlying buffers.
	CleanupFunc func() error
}

// BlobStorage unit of non-volatile binary large objects (BLOB) persistence.
type BlobStorage interface {
	// Upload stores an Object in a remote blob storage.
	Upload(ctx context.Context, obj Object) error

	// CheckMod verifies if an Object (using its key) was modified prior a specified time or differs from
	// size compared to the Object stored in the remote storage.
	//
	// Returns ErrFatalStorage if non-recovery operation was returned from remote storage server
	// (e.g. insufficient permissions, bucket does not exists).
	CheckMod(ctx context.Context, key string, modTime time.Time, size int64) (bool, error)
}

type NoopBlobStorage struct {
	UploadErr    error
	CheckModBool bool
	CheckModErr  error
}

var _ BlobStorage = NoopBlobStorage{}

func (n NoopBlobStorage) Upload(_ context.Context, _ Object) error {
	return n.UploadErr
}

func (n NoopBlobStorage) CheckMod(_ context.Context, _ string, _ time.Time, _ int64) (bool, error) {
	return n.CheckModBool, n.CheckModErr
}
