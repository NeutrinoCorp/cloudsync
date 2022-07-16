package cloudsync

import (
	"errors"
	"fmt"
)

type ErrFileUpload struct {
	Key    string
	Parent error
}

var _ error = ErrFileUpload{}

func (e ErrFileUpload) Error() string {
	return fmt.Sprintf("cloudsync: Failed to parse file with key %s", e.Key)
}

var ErrFatalStorage = errors.New("cloudsync: Got fatal error from blob storage")
