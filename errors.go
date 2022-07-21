package cloudsync

import (
	"errors"
	"fmt"
)

// ErrFatalStorage non-recovery error issued by the blob storage. Programs should panic once they receive this error.
var ErrFatalStorage = errors.New("cloudsync: Got fatal error from blob storage")

// ErrFileUpload generic error generated from a blob upload job.
type ErrFileUpload struct {
	Key    string
	Parent error
}

var _ error = ErrFileUpload{}

func (e ErrFileUpload) Error() string {
	if e.Key == "" {
		return "cloudsync: File upload failed"
	}
	return fmt.Sprintf("cloudsync: File upload failed with key %s", e.Key)
}
