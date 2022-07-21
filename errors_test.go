package cloudsync_test

import (
	"testing"

	"github.com/neutrinocorp/cloudsync"
	"github.com/stretchr/testify/assert"
)

func TestErrFileUpload_Error(t *testing.T) {
	tests := []struct {
		name string
		key  string
		exp  string
	}{
		{
			name: "Empty",
			key:  "",
			exp:  "cloudsync: File upload failed",
		},
		{
			name: "With key",
			key:  "foo",
			exp:  "cloudsync: File upload failed with key foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := cloudsync.ErrFileUpload{
				Key: tt.key,
			}
			assert.Equal(t, tt.exp, out.Error())
		})
	}
}
