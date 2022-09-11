package rebalancer_test

import (
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/offerm/lnagent/lightning"
	"github.com/offerm/lnagent/protobuf"
	"github.com/offerm/lnagent/rebalancer"

	//"github.com/offerm/lnagent/rebalancer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRebalancer_TaskInit(t *testing.T) {
	pubkeyB, _ := hex.DecodeString("02b998d8c3f065f3e0a8b383bd00dff56aeeac05c52ea2b7a5c936ff8ab2fb369a")
	pubkeyC, _ := hex.DecodeString("02848fffeb2ebaafdcd6b795b3a45d1e2397181e1c0d4424e86661276bfbe815a9")
	AtoB := uint64(2542708600164515840)
	CtoA := uint64(2542877924953358337)
	amount := uint64(1000 * 1000)

	events := make(chan *protobuf.TaskResponse, 1)
	service := lightning.NewMockService(&lightning.Config{})
	rebalancerA := rebalancer.NewRebalancer(events, service)

	sid := rebalancer.SwapID(uuid.NewString())
	rebalancerA.TaskInit(sid, &protobuf.Task_Init{
		Role: protobuf.Task_INITIATOR,
		From: &protobuf.Payment{
			PeerPubKey: pubkeyC,
			ChanId:     CtoA,
			AmtMsat:    amount,
			FeeMsat:    0,
			TimeLock:   0,
		},
		To: &protobuf.Payment{ //to B
			PeerPubKey: pubkeyB,
			ChanId:     AtoB,
			AmtMsat:    amount,
			FeeMsat:    0,
			TimeLock:   0,
		},
		PaymentRequest: "",
	})
	initA := <-events
	assert.IsType(t, &protobuf.TaskResponse_InitDoneType{}, initA.Response)
	assert.Equal(t, initA.Swap_ID, string(sid))
}
