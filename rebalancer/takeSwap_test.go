package rebalancer

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/offerm/lnagent/lightning"
	"github.com/offerm/lnagent/protobuf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTakeSwap(t *testing.T) {
	pubkeyB, _ := hex.DecodeString("02b998d8c3f065f3e0a8b383bd00dff56aeeac05c52ea2b7a5c936ff8ab2fb369a")
	pubkeyC, _ := hex.DecodeString("02848fffeb2ebaafdcd6b795b3a45d1e2397181e1c0d4424e86661276bfbe815a9")
	pubkey123, _ := hex.DecodeString("123")
	AtoB := uint64(2542708600164515840)
	CtoA := uint64(2542877924953358337)
	amount := uint64(1000 * 1000)

	eventsA := make(chan *protobuf.TaskResponse, 1)
	serviceA := &lightning.TestMockService{}
	rebalancerA := NewRebalancer(eventsA, serviceA)
	sid := SwapID(uuid.NewString())
	newSid := SwapID(uuid.NewString())

	eventsB := make(chan *protobuf.TaskResponse, 1) // for the wrong role test
	serviceB := &lightning.TestMockService{}
	rebalancerB := NewRebalancer(eventsB, serviceB)

	eventsC := make(chan *protobuf.TaskResponse, 1) // for the MakeHashPaymentAndMonitor test
	serviceC := &lightning.TestMockService{}
	rebalancerC := NewRebalancer(eventsC, serviceC)

	eventsD := make(chan *protobuf.TaskResponse, 1)
	serviceD := &lightning.TestMockService{}
	rebalancerD := NewRebalancer(eventsD, serviceD)

	serviceA.SaveInfo(pubkeyC, pubkeyB, CtoA, AtoB, amount, "", t) //todo do for all

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

	initA := <-eventsA
	serviceA.UpdatePayReq(initA.GetInitDoneType().PaymentRequest)
	serviceB.SaveInfo(pubkeyC, pubkeyB, CtoA, AtoB, amount, initA.GetInitDoneType().PaymentRequest, t)

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

	initB := <-eventsB
	serviceB.UpdatePayReq(initB.GetInitDoneType().PaymentRequest)

	serviceC.SaveInfo(pubkeyC, pubkey123, CtoA, AtoB, amount, "", t) //todo do for all

	rebalancerC.TaskInit(sid, &protobuf.Task_Init{ //for hash and monitor test
		Role: protobuf.Task_INITIATOR,
		From: &protobuf.Payment{
			PeerPubKey: pubkeyC,
			ChanId:     CtoA,
			AmtMsat:    amount,
			FeeMsat:    0,
			TimeLock:   0,
		},
		To: &protobuf.Payment{
			PeerPubKey: pubkey123,
			ChanId:     AtoB,
			AmtMsat:    amount,
			FeeMsat:    0,
			TimeLock:   0,
		},
		PaymentRequest: "",
	})
	initC := <-eventsC
	serviceC.UpdatePayReq(initC.GetInitDoneType().PaymentRequest)
	serviceD.SaveInfo(pubkeyC, pubkeyB, CtoA, AtoB, 10, "", t) //todo do for all

	rebalancerD.TaskInit(sid, &protobuf.Task_Init{
		Role: protobuf.Task_INITIATOR,
		From: &protobuf.Payment{
			PeerPubKey: pubkeyC,
			ChanId:     CtoA,
			AmtMsat:    10,
			FeeMsat:    0,
			TimeLock:   0,
		},
		To: &protobuf.Payment{
			PeerPubKey: pubkeyB,
			ChanId:     AtoB,
			AmtMsat:    10,
			FeeMsat:    0,
			TimeLock:   0,
		},
		PaymentRequest: "",
	})

	initD := <-eventsD
	serviceD.UpdatePayReq(initD.GetInitDoneType().PaymentRequest)

	type args struct {
		task       *protobuf.Task_Swap
		swapID     SwapID
		rebalancer *Rebalancer
	}

	tests := []struct {
		name             string
		args             args
		expectedResponse *protobuf.TaskResponse
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
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID:  string(sid),
				Response: &protobuf.TaskResponse_PaymentInitiatedType{},
			},
		},
		{
			name: "new swapID ",
			args: args{
				task: &protobuf.Task_Swap{
					PaymentRequest: initA.GetInitDoneType().PaymentRequest,
				},
				swapID:     newSid,
				rebalancer: rebalancerA,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(newSid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("swap %v does not exist", newSid),
					},
				},
			},
		},
		{
			name: "invalid role ",
			args: args{
				task: &protobuf.Task_Swap{
					PaymentRequest: "",
				},
				swapID:     sid,
				rebalancer: rebalancerB,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(newSid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("swap can't start by non initiator"),
					},
				},
			},
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
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(newSid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("failed to decode payreq 1111 - %v", lnwire.NewError()),
					},
				},
			},
		},
		{
			name: "MakeHashPaymentAndMonitor test",
			args: args{
				task: &protobuf.Task_Swap{
					PaymentRequest: initC.GetInitDoneType().PaymentRequest,
				},
				swapID:     sid,
				rebalancer: rebalancerC,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(newSid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("initiator failed to initiate payment - %v ", lnwire.NewError()),
					},
				},
			},
		},
		{
			name: "swap hash compare test",
			args: args{
				task: &protobuf.Task_Swap{
					PaymentRequest: "1234",
				},
				swapID:     sid,
				rebalancer: rebalancerA,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(newSid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("wrong hash %v received for payment %v", "1234", sid),
					},
				},
			},
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
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(newSid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("invalid payment request hash g"),
					},
				},
			},
		},
		{
			name: "amount test",
			args: args{
				task: &protobuf.Task_Swap{
					PaymentRequest: initD.GetInitDoneType().PaymentRequest,
				},
				swapID:     sid,
				rebalancer: rebalancerD,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(newSid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("different amounts. expected 10 invoiced 1000000"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.rebalancer.TaskSwap(tt.args.swapID, tt.args.task)
			select {
			case finalResponse, _ := <-tt.args.rebalancer.events:
				assert.IsType(t, tt.expectedResponse.Response, finalResponse.Response)
				if finalResponse.GetErrorType() != nil { //error type
					assert.Equal(t, tt.expectedResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error, finalResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error)
				} else {
					assert.Equal(t, tt.expectedResponse.Response.(*protobuf.TaskResponse_PaymentInitiatedType).PaymentInitiatedType, finalResponse.GetPaymentInitiatedType())
				}
				//default:

			}

		})
	}

}
