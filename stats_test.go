package cloudsync

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStats_GetTotalUploadJobs(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	stats := &Stats{}
	go func() {
		stats.increaseUploadJobs()
		wg.Done()
	}()
	go func() {
		stats.increaseUploadJobs()
		wg.Done()
	}()
	go func() {
		stats.decreaseUploadJobs()
		wg.Done()
	}()
	wg.Wait()
	assert.Equal(t, uint64(2), stats.GetTotalUploadJobs())
}

func TestStats_GetCurrentUploadJobs(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	stats := &Stats{}
	go func() {
		stats.increaseUploadJobs()
		wg.Done()
	}()
	go func() {
		stats.increaseUploadJobs()
		wg.Done()
	}()
	go func() {
		stats.decreaseUploadJobs()
		wg.Done()
	}()
	wg.Wait()
	assert.Equal(t, uint64(1), stats.GetCurrentUploadJobs())
}

func TestStats_GetTotalFailedJobs(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	stats := &Stats{}
	go func() {
		stats.increaseFailedJobs()
		wg.Done()
	}()
	go func() {
		stats.increaseFailedJobs()
		wg.Done()
	}()
	go func() {
		stats.increaseFailedJobs()
		wg.Done()
	}()
	wg.Wait()
	assert.Equal(t, uint64(3), stats.GetTotalFailedJobs())
}
