package cloudsync

import "sync/atomic"

type Stats struct {
	currentUploadJobs uint64
	totalUploadJobs   uint64
}

var DefaultStats = &Stats{}

func (s *Stats) increaseUploadJobs() {
	atomic.AddUint64(&s.currentUploadJobs, 1)
	atomic.AddUint64(&s.totalUploadJobs, 1)
}

func (s *Stats) decreaseUploadJobs() {
	prev := atomic.LoadUint64(&s.currentUploadJobs)
	atomic.StoreUint64(&s.currentUploadJobs, prev-1)
}

func (s Stats) GetTotalUploadJobs() uint64 {
	return atomic.LoadUint64(&s.totalUploadJobs)
}

func (s Stats) GetCurrentUploadJobs() uint64 {
	return atomic.LoadUint64(&s.currentUploadJobs)
}
