package lnagent

// TODO: test via grpc

import (
	"github.com/offerm/lnagent/rebalancer/mocking"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

func TestStoppableAgent(t *testing.T) {
	startGoRoutins := runtime.NumGoroutine()
	agent := NewAgent(&Config{Host: "127.0.0.1", Port: 8888}, mocking.NewAgentMockService())
	go func() {
		time.Sleep(11 * time.Second)
		// 2 new goroutines added by the loop func (not 3 - no connection to coordinator)
		// 1 new goroutine added by this test
		assert.Equal(t, startGoRoutins+3, runtime.NumGoroutine())
		agent.Stop()
		assert.Equal(t, startGoRoutins+1, runtime.NumGoroutine()) // no new goroutines left except this new one

	}()
	agent.Run()
}
