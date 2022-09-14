package rebalancer_test

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/offerm/lnagent/protobuf"
	"github.com/offerm/lnagent/rebalancer"
	"github.com/offerm/lnagent/rebalancer/mocking"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTakeSwap(t *testing.T) {
	const (
		validPayReq       = "validPayReq"
		emptyPayReq       = ""
		invalidPayReq     = "invalidPayReq"
		decodeTestPayReq  = "decode"
		swapComparePayReq = "swapCompareTestPayReq"
	)

	pubkeyB, _ := hex.DecodeString("02b998d8c3f065f3e0a8b383bd00dff56aeeac05c52ea2b7a5c936ff8ab2fb369a")
	pubkeyC, _ := hex.DecodeString("02848fffeb2ebaafdcd6b795b3a45d1e2397181e1c0d4424e86661276bfbe815a9")
	pubkey123, _ := hex.DecodeString("123")
	AtoB := uint64(2542708600164515840)
	CtoA := uint64(2542877924953358337)
	amount := uint64(1000 * 1000)

	sid := rebalancer.SwapID(uuid.NewString())

	tests := []struct {
		name             string
		initTask         *protobuf.Task_Init
		swapTask         *protobuf.Task_Swap
		expectedResponse *protobuf.TaskResponse
	}{
		{
			name: "validPayReq Swap",
			initTask: &protobuf.Task_Init{
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
				PaymentRequest: emptyPayReq,
			},
			swapTask: &protobuf.Task_Swap{
				PaymentRequest: validPayReq,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID:  string(sid),
				Response: &protobuf.TaskResponse_PaymentInitiatedType{},
			},
		},
		{
			//testing that an error will be returned if the swapID is new during a swap task
			name:     "unexpected swapID ",
			initTask: nil,
			swapTask: &protobuf.Task_Swap{
				PaymentRequest: validPayReq,
			},

			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(sid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: "does not exist",
					},
				},
			},
		},
		{
			//testing that an error will be returned if the role is not initiator
			name: "invalid role ",
			initTask: &protobuf.Task_Init{
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
				PaymentRequest: validPayReq,
			},
			swapTask: &protobuf.Task_Swap{
				PaymentRequest: emptyPayReq,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(sid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: "swap can't start by non initiator",
					},
				},
			},
		},
		{
			//testing the response of an error returned from the decodePayReq func
			name: "invalid payRequest ",
			initTask: &protobuf.Task_Init{
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
				PaymentRequest: emptyPayReq,
			},
			swapTask: &protobuf.Task_Swap{
				PaymentRequest: invalidPayReq,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(sid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("failed to decode payreq invalidPayReq - %v", lnwire.NewError()),
					},
				},
			},
		},
		{
			//testing the error of MakeHashPaymentAndMonitor func by setting the To.PeerPubKey to be pubkey123
			name: "MakeHashPaymentAndMonitor test",
			initTask: &protobuf.Task_Init{
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
				PaymentRequest: emptyPayReq,
			},
			swapTask: &protobuf.Task_Swap{
				PaymentRequest: validPayReq,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(sid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("initiator failed to initiate payment - %v ", lnwire.NewError()),
					},
				},
			},
		},
		{
			//testing the response when the hash returned from the decodePayReq func is different from expected
			name: "swap hash compare test",
			initTask: &protobuf.Task_Init{
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
				PaymentRequest: emptyPayReq,
			},
			swapTask: &protobuf.Task_Swap{
				PaymentRequest: swapComparePayReq,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(sid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: "wrong hash 1234 received for payment",
					},
				},
			},
		},
		{
			//testing the response after an error is returned from the hex decode func
			name: "hex decode test",
			initTask: &protobuf.Task_Init{
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
				PaymentRequest: emptyPayReq,
			},
			swapTask: &protobuf.Task_Swap{
				PaymentRequest: decodeTestPayReq,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(sid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: "invalid payment request hash g",
					},
				},
			},
		},
		{
			//testing the response when the amount returned from decodePayReq func does not match the expected amount
			name: "amount test",
			initTask: &protobuf.Task_Init{
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
				PaymentRequest: emptyPayReq,
			},
			swapTask: &protobuf.Task_Swap{
				PaymentRequest: validPayReq,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(sid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: "different amounts. expected 10 invoiced 1000000",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events := make(chan *protobuf.TaskResponse, 1)
			service := &mocking.TestMockService{}
			rebalancer := rebalancer.NewRebalancer(events, service)

			if tt.initTask != nil {
				if tt.initTask.Role == protobuf.Task_INITIATOR {
					service.SaveInfo(pubkeyC, tt.initTask.To.PeerPubKey, CtoA, AtoB, tt.initTask.To.AmtMsat, tt.swapTask.PaymentRequest, t)
				} else {
					service.SaveInfo(pubkeyC, tt.initTask.To.PeerPubKey, CtoA, AtoB, tt.initTask.To.AmtMsat, tt.initTask.PaymentRequest, t)
				}
				rebalancer.TaskInit(sid, tt.initTask)
				<-events
			} else {
				service.SaveInfo(pubkeyC, pubkeyB, CtoA, AtoB, amount, "", t)
			}

			rebalancer.TaskSwap(sid, tt.swapTask)
			finalResponse, _ := <-events
			assert.IsType(t, tt.expectedResponse.Response, finalResponse.Response)
			if finalResponse.GetErrorType() != nil { //error type
				assert.Contains(t, finalResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error, tt.expectedResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error)
			} else {
				assert.Equal(t, tt.expectedResponse.Response.(*protobuf.TaskResponse_PaymentInitiatedType).PaymentInitiatedType, finalResponse.GetPaymentInitiatedType())
			}
			assert.Equal(t, tt.expectedResponse.Swap_ID, finalResponse.Swap_ID)

		})
	}

}
