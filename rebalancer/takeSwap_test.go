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

	sid := SwapID(uuid.NewString())
	newSid := SwapID(uuid.NewString())

	type args struct {
		initTask *protobuf.Task_Init
		swapTask *protobuf.Task_Swap
		swapID   SwapID
	}

	tests := []struct {
		name             string
		args             args
		expectedResponse *protobuf.TaskResponse
	}{
		{
			name: "valid Swap",
			args: args{
				swapID: sid,
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
					PaymentRequest: "",
				},
				swapTask: &protobuf.Task_Swap{
					PaymentRequest: "normal",
				},
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID:  string(sid),
				Response: &protobuf.TaskResponse_PaymentInitiatedType{},
			},
		},
		{
			name: "new swapID ",
			args: args{
				swapID: newSid,
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
					PaymentRequest: "",
				},
				swapTask: &protobuf.Task_Swap{
					PaymentRequest: "normal",
				},
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(newSid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("does not exist"),
					},
				},
			},
		},
		{
			name: "invalid role ",
			args: args{
				swapID: sid,
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
					PaymentRequest: "invalid role",
				},
				swapTask: &protobuf.Task_Swap{
					PaymentRequest: "",
				},
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
				swapID: sid,
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
					PaymentRequest: "",
				},
				swapTask: &protobuf.Task_Swap{
					PaymentRequest: "invalidPayReq",
				},
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(newSid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("failed to decode payreq invalidPayReq - %v", lnwire.NewError()),
					},
				},
			},
		},
		{
			name: "MakeHashPaymentAndMonitor test",
			args: args{
				swapID: sid,
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
					PaymentRequest: "",
				},
				swapTask: &protobuf.Task_Swap{
					PaymentRequest: "normal",
				},
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
				swapID: sid,
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
					PaymentRequest: "",
				},
				swapTask: &protobuf.Task_Swap{
					PaymentRequest: "swapCompareTest",
				},
			},
			expectedResponse: &protobuf.TaskResponse{
				Swap_ID: string(newSid),
				Response: &protobuf.TaskResponse_ErrorType{
					ErrorType: &protobuf.TaskResponse_Error{
						Error: fmt.Sprintf("wrong hash 1234 received for payment"),
					},
				},
			},
		},
		{
			name: "hex decode test",
			args: args{
				swapID: sid,
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
					PaymentRequest: "",
				},
				swapTask: &protobuf.Task_Swap{
					PaymentRequest: "decode",
				},
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
				swapID: sid,
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
					PaymentRequest: "",
				},
				swapTask: &protobuf.Task_Swap{
					PaymentRequest: "normal",
				},
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
			events := make(chan *protobuf.TaskResponse, 1)
			service := &lightning.TestMockService{}
			rebalancer := NewRebalancer(events, service)

			service.SaveInfo(pubkeyC, tt.args.initTask.To.PeerPubKey, CtoA, AtoB, tt.args.initTask.To.AmtMsat, tt.args.initTask.PaymentRequest, t)

			rebalancer.TaskInit(sid, tt.args.initTask)

			initA := <-events

			if tt.args.swapTask.PaymentRequest == "normal" {
				tt.args.swapTask.PaymentRequest = initA.GetInitDoneType().PaymentRequest
			}

			service.UpdatePayReq(tt.args.swapTask.PaymentRequest)
			rebalancer.TaskSwap(tt.args.swapID, tt.args.swapTask)
			select {
			case finalResponse, _ := <-rebalancer.events:
				assert.IsType(t, tt.expectedResponse.Response, finalResponse.Response)
				if finalResponse.GetErrorType() != nil { //error type
					assert.Contains(t, finalResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error, tt.expectedResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error)
				} else {
					assert.Equal(t, tt.expectedResponse.Response.(*protobuf.TaskResponse_PaymentInitiatedType).PaymentInitiatedType, finalResponse.GetPaymentInitiatedType())
				}
				//default:

			}

		})
	}

}
