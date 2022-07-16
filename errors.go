package cloudsync

import "fmt"

type ErrFileUpload struct {
	Key    string
	parent error
}

var _ error = ErrFileUpload{}

func (e ErrFileUpload) Error() string {
	return fmt.Sprintf("cloudsync: Failed to parse file with key %s", e.Key)
}
