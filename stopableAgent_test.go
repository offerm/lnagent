package lnagent

// TODO: test via grpc

import (
	"github.com/offerm/lnagent/lightning"
	"github.com/offerm/lnagent/rebalancer/mocking"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStoppableAgent(t *testing.T) {
	//go newMockCoordinator()
	agent := NewAgent(&Config{Host: "127.0.0.1", Port: 8888}, mocking.NewMockService(&lightning.Config{}))
	agent.Run()
	time.Sleep(3 * time.Second)
	agent.Stop()
	assert.Equal(t, 1, 1)

	//serviceA := mocking.NewMockService(&lightning.Config{})
	//agentA := NewAgent(nil, serviceA)
	//
	//serviceB := mocking.NewMockService(&lightning.Config{})
	//agentB := NewAgent(nil, serviceB)
	//
	//serviceC := mocking.NewMockService(&lightning.Config{})
	//agentC := NewAgent(nil, serviceC)
	//
	//pubkeyA, _ := hex.DecodeString("02aeb304f6282f6ab93bf1b7cfdf0a7e842ccef33455f201484cc3a3d316edabb7")
	//pubkeyB, _ := hex.DecodeString("02b998d8c3f065f3e0a8b383bd00dff56aeeac05c52ea2b7a5c936ff8ab2fb369a")
	//pubkeyC, _ := hex.DecodeString("02848fffeb2ebaafdcd6b795b3a45d1e2397181e1c0d4424e86661276bfbe815a9")
	//AtoB := uint64(2542708600164515840)
	//CtoA := uint64(2542877924953358337)
	//BtoC := uint64(2542904313232752641)
	//amount := uint64(1000 * 1000)
	//sid := rebalancer.SwapID(uuid.NewString())
	//
	//// swap order A->B->C->A
	//// init order A, C, B
	//
	////	go agent.Run()
	//agentA.rebalancer.TaskInit(sid, &protobuf.Task_Init{
	//	Role: protobuf.Task_INITIATOR,
	//	From: &protobuf.Payment{
	//		PeerPubKey: pubkeyC,
	//		ChanId:     CtoA,
	//		AmtMsat:    amount,
	//		FeeMsat:    0,
	//		TimeLock:   0,
	//	},
	//	To: &protobuf.Payment{ //to B
	//		PeerPubKey: pubkeyB,
	//		ChanId:     AtoB,
	//		AmtMsat:    amount,
	//		FeeMsat:    0,
	//		TimeLock:   0,
	//	},
	//	PaymentRequest: "",
	//})
	//initA := <-agentA.events
	//
	//assert.IsType(t, &protobuf.TaskResponse_InitDoneType{}, initA.Response)
	//assert.Equal(t, initA.Swap_ID, string(sid))
	//
	//agentC.rebalancer.TaskInit(sid, &protobuf.Task_Init{
	//	Role: protobuf.Task_FOLLOWER,
	//	From: &protobuf.Payment{
	//		PeerPubKey: pubkeyB,
	//		ChanId:     BtoC,
	//		AmtMsat:    amount,
	//		FeeMsat:    0,
	//		TimeLock:   0,
	//	},
	//	To: &protobuf.Payment{
	//		PeerPubKey: pubkeyA,
	//		ChanId:     CtoA,
	//		AmtMsat:    amount,
	//		FeeMsat:    0,
	//		TimeLock:   0,
	//	},
	//	PaymentRequest: initA.GetInitDoneType().PaymentRequest,
	//})
	//initC := <-agentC.events
	//
	//assert.IsType(t, &protobuf.TaskResponse_InitDoneType{}, initC.Response)
	//assert.Equal(t, initC.Swap_ID, string(sid))
	//
	//agentB.rebalancer.TaskInit(sid, &protobuf.Task_Init{
	//	Role: protobuf.Task_FOLLOWER,
	//	From: &protobuf.Payment{ // from A
	//		PeerPubKey: pubkeyA,
	//		ChanId:     AtoB,
	//		AmtMsat:    amount,
	//		FeeMsat:    0,
	//		TimeLock:   0,
	//	},
	//	To: &protobuf.Payment{
	//		PeerPubKey: pubkeyC,
	//		ChanId:     BtoC,
	//		AmtMsat:    amount,
	//		FeeMsat:    0,
	//		TimeLock:   0,
	//	},
	//	PaymentRequest: initC.GetInitDoneType().PaymentRequest,
	//})
	//initB := <-agentB.events
	//
	//assert.IsType(t, &protobuf.TaskResponse_InitDoneType{}, initC.Response)
	//assert.Equal(t, initB.Swap_ID, string(sid))
	//
	//// we give to A B's payment request so A can start the payment
	//agentA.rebalancer.TaskSwap(sid, &protobuf.Task_Swap{
	//	PaymentRequest: initB.GetInitDoneType().PaymentRequest,
	//})
	//swapA := <-agentA.events
	//
	//assert.IsType(t, &protobuf.TaskResponse_PaymentInitiatedType{}, swapA.Response)
	//assert.Equal(t, swapA.Swap_ID, string(sid))
	//
	//for {
	//	event := <-agentA.events
	//	if event.GetSwapDoneType() != nil {
	//		break
	//	}
	//}
}

//package lnagent
//
//import (
//	"github.com/offerm/lnagent/lightning"
//	"github.com/stretchr/testify/assert"
//	"testing"
//	"time"
//)
//
//func TestStoppableAgent(t *testing.T) {
//	go newMockCoordinator()
//agent := NewAgent(&Config{Host: "127.0.0.1", Port: 8888}, lightning.NewService(&lightning.Config{Host: "127.0.0.1", Port: 8888}))
//agent.Run()
//time.Sleep(3 * time.Second)
//agent.Stop()
//assert.Equal(t, 1, 1)
//const (
//	validPayReq       = "validPayReq"
//	emptyPayReq       = ""
//	invalidPayReq     = "invalidPayReq"
//	decodeTestPayReq  = "decode"
//	swapComparePayReq = "swapCompareTestPayReq"
//)

//pubkeyB, _ := hex.DecodeString("02b998d8c3f065f3e0a8b383bd00dff56aeeac05c52ea2b7a5c936ff8ab2fb369a")
//pubkeyC, _ := hex.DecodeString("02848fffeb2ebaafdcd6b795b3a45d1e2397181e1c0d4424e86661276bfbe815a9")
//pubkey123, _ := hex.DecodeString("123")
//AtoB := uint64(2542708600164515840)
//CtoA := uint64(2542877924953358337)
//amount := uint64(1000 * 1000)
//
//sid := rebalancer.SwapID(uuid.NewString())

//tests := []struct {
//	name             string
//	initTask         *protobuf.Task_Init
//	swapTask         *protobuf.Task_Swap
//	expectedResponse *protobuf.TaskResponse
//}{
//	{
//		//	name: "validPayReq Swap",
//		//	initTask: &protobuf.Task_Init{
//		//		Role: protobuf.Task_INITIATOR,
//		//		From: &protobuf.Payment{
//		//			PeerPubKey: pubkeyC,
//		//			ChanId:     CtoA,
//		//			AmtMsat:    amount,
//		//			FeeMsat:    0,
//		//			TimeLock:   0,
//		//		},
//		//		To: &protobuf.Payment{
//		//			PeerPubKey: pubkeyB,
//		//			ChanId:     AtoB,
//		//			AmtMsat:    amount,
//		//			FeeMsat:    0,
//		//			TimeLock:   0,
//		//		},
//		//		PaymentRequest: emptyPayReq,
//		//	},
//		//	swapTask: &protobuf.Task_Swap{
//		//		PaymentRequest: validPayReq,
//		//	},
//		//	expectedResponse: &protobuf.TaskResponse{
//		//		Swap_ID:  string(sid),
//		//		Response: &protobuf.TaskResponse_PaymentInitiatedType{},
//		//	},
//	},
//}
//
//for _, tt := range tests {
//	t.Run(tt.name, func(t *testing.T) {
//		//	events := make(chan *protobuf.TaskResponse, 1)
//		//	service := &mocking.SwapMockService{}
//		//	rebalancer := rebalancer.NewRebalancer(events, service)
//		//
//		//	if tt.initTask != nil {
//		//		if tt.initTask.Role == protobuf.Task_INITIATOR {
//		//			service.SaveInfo(pubkeyC, tt.initTask.To.PeerPubKey, CtoA, AtoB, tt.initTask.To.AmtMsat, tt.swapTask.PaymentRequest, t)
//		//		} else {
//		//			service.SaveInfo(pubkeyC, tt.initTask.To.PeerPubKey, CtoA, AtoB, tt.initTask.To.AmtMsat, tt.initTask.PaymentRequest, t)
//		//		}
//		//		rebalancer.TaskInit(sid, tt.initTask)
//		//		<-events
//		//	} else {
//		//		service.SaveInfo(pubkeyC, pubkeyB, CtoA, AtoB, amount, "", t)
//		//	}
//		//
//		//	rebalancer.TaskSwap(sid, tt.swapTask)
//		//	finalResponse, _ := <-events
//		//	assert.IsType(t, tt.expectedResponse.Response, finalResponse.Response)
//		//	if finalResponse.GetErrorType() != nil { //error type
//		//		assert.Contains(t, finalResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error, tt.expectedResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error)
//		//	} else {
//		//		assert.Equal(t, tt.expectedResponse.Response.(*protobuf.TaskResponse_PaymentInitiatedType).PaymentInitiatedType, finalResponse.GetPaymentInitiatedType())
//		//	}
//		//	assert.Equal(t, tt.expectedResponse.Swap_ID, finalResponse.Swap_ID)
//		//
//	})
//}
