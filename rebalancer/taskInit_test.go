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

const (
	emptyPayReq         = ""
	validResponsePayReq = "1234"
	notEmptyPayReq      = "not empty payment request"
	invalidRole         = 3
	invalidPayReq       = "invalidPayReq"
	decodeTestPayReq    = "decode"
)

func TestTaskInit(t *testing.T) {
	pubkeyB, _ := hex.DecodeString("02b998d8c3f065f3e0a8b383bd00dff56aeeac05c52ea2b7a5c936ff8ab2fb369a")
	pubkeyC, _ := hex.DecodeString("02848fffeb2ebaafdcd6b795b3a45d1e2397181e1c0d4424e86661276bfbe815a9")
	AtoB := uint64(2542708600164515840)
	CtoA := uint64(2542877924953358337)
	amount := uint64(1000 * 1000)

	sid := SwapID(uuid.NewString())
	validSidArr := []SwapID{sid}
	invalidSidArr := []SwapID{sid, sid}

	tests := []struct {
		name             string
		swapIdArr        []SwapID
		task             *protobuf.Task_Init
		expectedResponse *protobuf.TaskResponse
	}{
		{
			// testing the response of a init task with a valid swap ID
			name:      "valid Swap ID test",
			swapIdArr: validSidArr,
			task: &protobuf.Task_Init{
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
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(sid),
				Response: &protobuf.TaskResponse_InitDoneType{
					InitDoneType: &protobuf.TaskResponse_Init_Done{
						PaymentRequest: validResponsePayReq,
					},
				},
			},
		},
		{
			//testing the response of an init task with a known swap ID
			//expecting the second init to fail
			name:      "known Swap ID test",
			swapIdArr: invalidSidArr,
			task: &protobuf.Task_Init{
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
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(sid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("swap %v is already active", sid),
					},
				},
			}},
		{
			//testing the response of an init task with a non-empty payReq and initiator role
			name:      "non-empty payment request with initiator role test",
			swapIdArr: validSidArr,
			task: &protobuf.Task_Init{
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
				PaymentRequest: notEmptyPayReq,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(sid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: "initiator should not accept a payment request ",
					},
				},
			}},
		{
			//testing the response of an init task with an invalid role (not initiator \ follower)
			name:      "invalid role test",
			swapIdArr: validSidArr,
			task: &protobuf.Task_Init{
				Role: invalidRole,
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
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(sid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("unknown rolee %v for swap %v ", invalidRole, sid),
					},
				},
			}},
		{
			//testing the response of an init task with an error returned from the decodePayReq func
			name:      "decodePayReq error test",
			swapIdArr: validSidArr,
			task: &protobuf.Task_Init{
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
				PaymentRequest: invalidPayReq,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(sid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("failed to decode payreq %v - %v", "invalidPayReq", lnwire.NewError()),
					},
				},
			}},
		{
			//testing the response of an init task getting an error from the decodeString func
			name:      "decodeString error test",
			swapIdArr: validSidArr,
			task: &protobuf.Task_Init{
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
				PaymentRequest: decodeTestPayReq,
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(sid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("invalid payment request hash %v", "g"),
					},
				},
			}},
		{
			//testing the response of an init task getting an error from the newHoldInvoice func
			name:      "newHoldInvoice error test",
			swapIdArr: []SwapID{SwapID("5656")},
			task: &protobuf.Task_Init{
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
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: "5656",
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("can't create invoice for swap %v - %v", "5656", lnwire.NewError()),
					},
				},
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events := make(chan *protobuf.TaskResponse, 1)
			service := &lightning.TestInitMockService{}
			rebalancer := NewRebalancer(events, service)

			for i := range tt.swapIdArr {

				service.SaveInfo(pubkeyC, tt.task.To.PeerPubKey, CtoA, AtoB, tt.task.To.AmtMsat, tt.task.PaymentRequest, t)
				rebalancer.TaskInit(tt.swapIdArr[i], tt.task)
				initA := <-events
				//make sure that if the response is not this test's final response the response is not an error response
				if i < len(tt.swapIdArr)-1 {
					assert.Nil(t, initA.GetErrorType())
				} else {
					assert.IsType(t, tt.expectedResponse.Response, initA.Response)
					if initA.GetErrorType() != nil { //error type
						assert.Contains(t, initA.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error, tt.expectedResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error)
					} else {
						assert.Equal(t, tt.expectedResponse.Response.(*protobuf.TaskResponse_InitDoneType).InitDoneType.PaymentRequest, initA.GetInitDoneType().PaymentRequest)
					}
					assert.Equal(t, tt.expectedResponse.Swap_ID, initA.Swap_ID)

				}
			}
		})
	}
}
