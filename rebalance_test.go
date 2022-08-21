package lnagent

// TODO: test via grpc

import (
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/offerm/lnagent/lightning"
	"github.com/offerm/lnagent/protobuf"
	"github.com/offerm/lnagent/rebalancer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCycle(t *testing.T) {
	agentA := NewAgent(nil, &lightning.Config{
		Host:           "localhost",
		Port:           10010,
		Network:        "testnet",
		Implementation: "",
		DataDir:        "~/lnd/lnda",
	})
	//agentA.Start(agentA, agentA.lnsdkConfig)

	agentB := NewAgent(nil, &lightning.Config{
		Host:           "localhost",
		Port:           10011,
		Network:        "testnet",
		Implementation: "",
		DataDir:        "~/lnd/lndb",
	})
	//agentB.Start(agentB, agentB.lnsdkConfig)

	agentC := NewAgent(nil, &lightning.Config{
		Host:           "localhost",
		Port:           10012,
		Network:        "testnet",
		Implementation: "",
		DataDir:        "~/lnd/lndc",
	})
	//agentC.Start(agentC, agentC.lnsdkConfig)

	pubkeyA, _ := hex.DecodeString("02aeb304f6282f6ab93bf1b7cfdf0a7e842ccef33455f201484cc3a3d316edabb7")
	pubkeyB, _ := hex.DecodeString("02b998d8c3f065f3e0a8b383bd00dff56aeeac05c52ea2b7a5c936ff8ab2fb369a")
	pubkeyC, _ := hex.DecodeString("02848fffeb2ebaafdcd6b795b3a45d1e2397181e1c0d4424e86661276bfbe815a9")
	AtoB := uint64(2542708600164515840)
	CtoA := uint64(2542877924953358337)
	BtoC := uint64(2542904313232752641)
	amount := uint64(1000 * 1000)
	sid := rebalancer.SwapID(uuid.NewString())

	// swap order A->B->C->A
	// init order A, C, B

	//	go agent.Run()
	agentA.rebalancer.TaskInit(sid, &protobuf.Task_Init{
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
	initA := <-agentA.events

	assert.IsType(t, &protobuf.TaskResponse_InitDoneType{}, initA.Response)
	assert.Equal(t, initA.Swap_ID, string(sid))

	agentC.rebalancer.TaskInit(sid, &protobuf.Task_Init{
		Role: protobuf.Task_FOLLOWER,
		From: &protobuf.Payment{
			PeerPubKey: pubkeyB,
			ChanId:     BtoC,
			AmtMsat:    amount,
			FeeMsat:    0,
			TimeLock:   0,
		},
		To: &protobuf.Payment{
			PeerPubKey: pubkeyA,
			ChanId:     CtoA,
			AmtMsat:    amount,
			FeeMsat:    0,
			TimeLock:   0,
		},
		PaymentRequest: initA.GetInitDoneType().PaymentRequest,
	})
	initC := <-agentC.events

	assert.IsType(t, &protobuf.TaskResponse_InitDoneType{}, initC.Response)
	assert.Equal(t, initC.Swap_ID, string(sid))

	agentB.rebalancer.TaskInit(sid, &protobuf.Task_Init{
		Role: protobuf.Task_FOLLOWER,
		From: &protobuf.Payment{ // from A
			PeerPubKey: pubkeyA,
			ChanId:     AtoB,
			AmtMsat:    amount,
			FeeMsat:    0,
			TimeLock:   0,
		},
		To: &protobuf.Payment{
			PeerPubKey: pubkeyC,
			ChanId:     BtoC,
			AmtMsat:    amount,
			FeeMsat:    0,
			TimeLock:   0,
		},
		PaymentRequest: initC.GetInitDoneType().PaymentRequest,
	})
	initB := <-agentB.events

	assert.IsType(t, &protobuf.TaskResponse_InitDoneType{}, initC.Response)
	assert.Equal(t, initB.Swap_ID, string(sid))

	// we give to A B's payment request so A can start the payment
	agentA.rebalancer.TaskSwap(sid, &protobuf.Task_Swap{
		PaymentRequest: initB.GetInitDoneType().PaymentRequest,
	})
	swapA := <-agentA.events

	assert.IsType(t, &protobuf.TaskResponse_PaymentInitiatedType{}, swapA.Response)
	assert.Equal(t, swapA.Swap_ID, string(sid))

	for {
		event := <-agentA.events
		if event.GetSwapDoneType() != nil {
			break
		}
	}
}
