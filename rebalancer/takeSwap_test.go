package rebalancer

import (
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/offerm/lnagent/lightning"
	"github.com/offerm/lnagent/protobuf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTakeSwap(t *testing.T) {
	pubkeyB, _ := hex.DecodeString("02b998d8c3f065f3e0a8b383bd00dff56aeeac05c52ea2b7a5c936ff8ab2fb369a")
	pubkeyC, _ := hex.DecodeString("02848fffeb2ebaafdcd6b795b3a45d1e2397181e1c0d4424e86661276bfbe815a9")
	AtoB := uint64(2542708600164515840)
	CtoA := uint64(2542877924953358337)
	amount := uint64(1000 * 1000)

	events := make(chan *protobuf.TaskResponse, 10)
	service := lightning.NewMockService(&lightning.Config{})
	rebalancerA := NewRebalancer(events, service)
	sid := SwapID(uuid.NewString())

	events2 := make(chan *protobuf.TaskResponse, 10)
	service2 := lightning.NewMockService(&lightning.Config{})
	rebalancerB := NewRebalancer(events2, service2)

	rebalancerA.TaskInit(sid, &protobuf.Task_Init{
		Role: protobuf.Task_INITIATOR,
		From: &protobuf.Payment{
			PeerPubKey: pubkeyC,
			ChanId:     CtoA,
			AmtMsat:    amount,
			FeeMsat:    0,
			TimeLock:   0,
		},
		To: &protobuf.Payment{
			PeerPubKey: pubkeyB,
			ChanId:     AtoB,
			AmtMsat:    amount,
			FeeMsat:    0,
			TimeLock:   0,
		},
		PaymentRequest: "",
	})
	initA := <-events

	rebalancerB.TaskInit(sid, &protobuf.Task_Init{
		Role: protobuf.Task_FOLLOWER,
		From: &protobuf.Payment{
			PeerPubKey: pubkeyC,
			ChanId:     CtoA,
			AmtMsat:    amount,
			FeeMsat:    0,
			TimeLock:   0,
		},
		To: &protobuf.Payment{
			PeerPubKey: pubkeyB,
			ChanId:     AtoB,
			AmtMsat:    amount,
			FeeMsat:    0,
			TimeLock:   0,
		},
		PaymentRequest: initA.GetInitDoneType().PaymentRequest,
	})

	type args struct {
		task       *protobuf.Task_Swap
		swapID     SwapID
		rebalancer *Rebalancer
	}
	tests := []struct {
		name            string
		args            args
		isExpectedError bool
	}{
		{
			name: "valid Swap",
			args: args{
				task: &protobuf.Task_Swap{
					PaymentRequest: initA.GetInitDoneType().PaymentRequest,
				},
				swapID:     sid,
				rebalancer: rebalancerA,
			},
			isExpectedError: false,
		},
		{
			name: "new swapID ",
			args: args{
				task: &protobuf.Task_Swap{
					PaymentRequest: initA.GetInitDoneType().PaymentRequest,
				},
				swapID:     SwapID(uuid.NewString()),
				rebalancer: rebalancerA,
			},
			isExpectedError: true,
		},
		{
			name: "invalid role ",
			args: args{
				task: &protobuf.Task_Swap{
					PaymentRequest: initA.GetInitDoneType().PaymentRequest,
				},
				swapID:     sid,
				rebalancer: rebalancerB,
			},
			isExpectedError: true,
		},
		{
			name: "invalid payRequest ", //todo merge this with the other branch
			args: args{
				task: &protobuf.Task_Swap{
					PaymentRequest: "1111",
				},
				swapID:     sid,
				rebalancer: rebalancerA,
			},
			isExpectedError: true,
		},
		//{
		//	name: "MakeHashPaymentAndMonitor test",
		//	args: args{
		//		task: &protobuf.Task_Swap{
		//			PaymentRequest: initD.GetInitDoneType().PaymentRequest,
		//		},
		//		swapID:     sid,
		//		rebalancer: rebalancerD,
		//	},
		//	isExpectedError: true,
		//},
		{
			name: "swap hash compare test",
			args: args{
				task: &protobuf.Task_Swap{
					PaymentRequest: "1234",
				},
				swapID:     sid,
				rebalancer: rebalancerA,
			},
			isExpectedError: true,
		},
		{
			name: "hex decode test",
			args: args{
				task: &protobuf.Task_Swap{
					PaymentRequest: "5678",
				},
				swapID:     sid,
				rebalancer: rebalancerA,
			},
			isExpectedError: true,
		},
		//{
		//	name: "amount test",
		//	args: args{
		//		task: &protobuf.Task_Swap{
		//			PaymentRequest: "1010",
		//		},
		//		swapID:     sid,
		//		rebalancer: rebalancerA,
		//	},
		//	isExpectedError: true,
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.rebalancer.TaskSwap(tt.args.swapID, tt.args.task)
			//swapA := <-tt.args.rebalancer.events
			if tt.isExpectedError {
				//for event := range tt.args.rebalancer.events {
				//	finalResponse = event
				//}
				select {
				case finalResponse, ok := <-tt.args.rebalancer.events:
					if !ok {
						assert.IsType(t, &protobuf.TaskResponse_ErrorType{}, finalResponse.Response)
					}
				default:

				}

			}

		})
	}

}
