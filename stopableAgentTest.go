package lnagent

import (
	"github.com/offerm/lnagent/protobuf"
	"testing"
)

func TestStopableAgent(t *testing.T) {

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

	tests := []struct {
		name             string
		initTask         *protobuf.Task_Init
		swapTask         *protobuf.Task_Swap
		expectedResponse *protobuf.TaskResponse
	}{
		{
			//	name: "validPayReq Swap",
			//	initTask: &protobuf.Task_Init{
			//		Role: protobuf.Task_INITIATOR,
			//		From: &protobuf.Payment{
			//			PeerPubKey: pubkeyC,
			//			ChanId:     CtoA,
			//			AmtMsat:    amount,
			//			FeeMsat:    0,
			//			TimeLock:   0,
			//		},
			//		To: &protobuf.Payment{
			//			PeerPubKey: pubkeyB,
			//			ChanId:     AtoB,
			//			AmtMsat:    amount,
			//			FeeMsat:    0,
			//			TimeLock:   0,
			//		},
			//		PaymentRequest: emptyPayReq,
			//	},
			//	swapTask: &protobuf.Task_Swap{
			//		PaymentRequest: validPayReq,
			//	},
			//	expectedResponse: &protobuf.TaskResponse{
			//		Swap_ID:  string(sid),
			//		Response: &protobuf.TaskResponse_PaymentInitiatedType{},
			//	},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//	events := make(chan *protobuf.TaskResponse, 1)
			//	service := &mocking.SwapMockService{}
			//	rebalancer := rebalancer.NewRebalancer(events, service)
			//
			//	if tt.initTask != nil {
			//		if tt.initTask.Role == protobuf.Task_INITIATOR {
			//			service.SaveInfo(pubkeyC, tt.initTask.To.PeerPubKey, CtoA, AtoB, tt.initTask.To.AmtMsat, tt.swapTask.PaymentRequest, t)
			//		} else {
			//			service.SaveInfo(pubkeyC, tt.initTask.To.PeerPubKey, CtoA, AtoB, tt.initTask.To.AmtMsat, tt.initTask.PaymentRequest, t)
			//		}
			//		rebalancer.TaskInit(sid, tt.initTask)
			//		<-events
			//	} else {
			//		service.SaveInfo(pubkeyC, pubkeyB, CtoA, AtoB, amount, "", t)
			//	}
			//
			//	rebalancer.TaskSwap(sid, tt.swapTask)
			//	finalResponse, _ := <-events
			//	assert.IsType(t, tt.expectedResponse.Response, finalResponse.Response)
			//	if finalResponse.GetErrorType() != nil { //error type
			//		assert.Contains(t, finalResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error, tt.expectedResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error)
			//	} else {
			//		assert.Equal(t, tt.expectedResponse.Response.(*protobuf.TaskResponse_PaymentInitiatedType).PaymentInitiatedType, finalResponse.GetPaymentInitiatedType())
			//	}
			//	assert.Equal(t, tt.expectedResponse.Swap_ID, finalResponse.Swap_ID)
			//
		})
	}

}
