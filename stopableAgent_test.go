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
	startGoRoutines := runtime.NumGoroutine()
	agent := NewAgent(&Config{Host: "127.0.0.1", Port: 8888}, mocking.NewAgentMockService())
	go func() {
		agent.Run()
	}()

	agent.Stop()
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, startGoRoutines+1, runtime.NumGoroutine()) // no new goroutines

}

func TestSignalStoppableAgent(t *testing.T) {
	startGoRoutines := runtime.NumGoroutine()
	agent := NewAgent(&Config{Host: "127.0.0.1", Port: 8888}, mocking.NewAgentMockService())
	go func() {
		agent.Run()
	}()

	time.Sleep(10 * time.Millisecond)

	//agent.signalChan <- os.Interrupt

	assert.Equal(t, startGoRoutines+1, runtime.NumGoroutine()) // no new goroutines

}
