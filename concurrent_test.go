package cloudsync

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

// required component to schedule goroutines eagerly, avoiding routine bad counters and
// collisions in workers listening to shared queue channels
type concurrentTestSuite struct {
	suite.Suite
	mu sync.Mutex
}

func TestConcurrentComponents(t *testing.T) {
	suite.Run(t, &concurrentTestSuite{})
}

func (s *concurrentTestSuite) SetupSuite() {
}

func (s *concurrentTestSuite) SetupTest() {
	s.mu.Lock()
	if objectUploadJobQueue == nil {
		objectUploadJobQueue = make(chan Object, 0)
	}
	if objectUploadJobQueueErr == nil {
		objectUploadJobQueueErr = make(chan ErrFileUpload, 0)
	}
}

func (s *concurrentTestSuite) TearDownSuite() {
}

func (s *concurrentTestSuite) TearDownTest() {
	defer s.mu.Unlock()
	if objectUploadJobQueue != nil {
		close(objectUploadJobQueue)
		objectUploadJobQueue = nil
	}
	if objectUploadJobQueueErr != nil {
		close(objectUploadJobQueueErr)
		objectUploadJobQueueErr = nil
	}
}

func (s *concurrentTestSuite) Test_ScheduleFileUploads() {
	testScheduleFileUploads(s.T())
}

func (s *concurrentTestSuite) Test_ListenForSysInterruption() {
	testListenForSysInterruption(s.T())
}

func (s *concurrentTestSuite) Test_ListenUpload() {
	testListenUpload(s.T())
}

func (s *concurrentTestSuite) Test_ListenUploadErrors() {
	testListenUploadErrors(s.T())
}

func (s *concurrentTestSuite) Test_ShutdownUploadWorkers() {
	testShutdownUploadWorkers(s.T())
}

func (s *concurrentTestSuite) Test_NewScanner() {
	testNewScanner(s.T())
}
