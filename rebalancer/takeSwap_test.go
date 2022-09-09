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

var initA *protobuf.TaskResponse

func TestTakeSwap(t *testing.T) {
	pubkeyB, _ := hex.DecodeString("02b998d8c3f065f3e0a8b383bd00dff56aeeac05c52ea2b7a5c936ff8ab2fb369a")
	pubkeyC, _ := hex.DecodeString("02848fffeb2ebaafdcd6b795b3a45d1e2397181e1c0d4424e86661276bfbe815a9")
	pubkey123, _ := hex.DecodeString("123")
	AtoB := uint64(2542708600164515840)
	CtoA := uint64(2542877924953358337)
	amount := uint64(1000 * 1000)

	events := make(chan *protobuf.TaskResponse, 1)
	service := &lightning.TestMockService{}
	rebalancer := NewRebalancer(events, service)

	sid := SwapID(uuid.NewString())
	newSid := SwapID(uuid.NewString())

	type args struct {
		task         *protobuf.Task_Swap
		newSid       bool
		role         protobuf.Task_Role
		amount       uint64
		payReq       string
		toPeerPubKey []byte
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
					PaymentRequest: "normal",
				},
				newSid:       false,
				role:         protobuf.Task_INITIATOR,
				amount:       amount,
				payReq:       "",
				toPeerPubKey: pubkeyB,
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
					PaymentRequest: "normal",
				},
				newSid:       true,
				role:         protobuf.Task_INITIATOR,
				amount:       amount,
				payReq:       "",
				toPeerPubKey: pubkeyB,
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
				task: &protobuf.Task_Swap{
					PaymentRequest: "",
				},
				newSid:       false,
				role:         protobuf.Task_FOLLOWER,
				amount:       amount,
				payReq:       "initA.GetInitDoneType().PaymentRequest",
				toPeerPubKey: pubkeyB,
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
					PaymentRequest: "invalidPayReq",
				},
				newSid:       false,
				role:         protobuf.Task_INITIATOR,
				amount:       amount,
				payReq:       "",
				toPeerPubKey: pubkeyB,
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
				task: &protobuf.Task_Swap{
					PaymentRequest: "normal",
				},
				newSid:       false,
				role:         protobuf.Task_INITIATOR,
				amount:       amount,
				payReq:       "",
				toPeerPubKey: pubkey123,
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
					PaymentRequest: "swap compare test",
				},
				newSid:       false,
				role:         protobuf.Task_INITIATOR,
				amount:       amount,
				payReq:       "",
				toPeerPubKey: pubkeyB,
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
				task: &protobuf.Task_Swap{
					PaymentRequest: "decode",
				},
				newSid:       false,
				role:         protobuf.Task_INITIATOR,
				amount:       amount,
				payReq:       "",
				toPeerPubKey: pubkeyB,
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
					PaymentRequest: "normal",
				},
				newSid:       false,
				role:         protobuf.Task_INITIATOR,
				amount:       10,
				payReq:       "",
				toPeerPubKey: pubkeyB,
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
			sid = SwapID(uuid.NewString())
			if tt.args.payReq != "" {

				rebalancer.TaskInit(SwapID(uuid.NewString()), &protobuf.Task_Init{
					Role: protobuf.Task_INITIATOR,
					From: &protobuf.Payment{
						PeerPubKey: pubkeyC,
						ChanId:     CtoA,
						AmtMsat:    tt.args.amount,
						FeeMsat:    0,
						TimeLock:   0,
					},
					To: &protobuf.Payment{
						PeerPubKey: pubkeyB,
						ChanId:     AtoB,
						AmtMsat:    tt.args.amount,
						FeeMsat:    0,
						TimeLock:   0,
					},
					PaymentRequest: "",
				})
				initA = <-events

				tt.args.payReq = initA.GetInitDoneType().PaymentRequest

				service.SaveInfo(pubkeyC, tt.args.toPeerPubKey, CtoA, AtoB, tt.args.amount, tt.args.payReq, t)

			}

			rebalancer.TaskInit(sid, &protobuf.Task_Init{
				Role: tt.args.role,
				From: &protobuf.Payment{
					PeerPubKey: pubkeyC,
					ChanId:     CtoA,
					AmtMsat:    tt.args.amount,
					FeeMsat:    0,
					TimeLock:   0,
				},
				To: &protobuf.Payment{
					PeerPubKey: tt.args.toPeerPubKey,
					ChanId:     AtoB,
					AmtMsat:    tt.args.amount,
					FeeMsat:    0,
					TimeLock:   0,
				},
				PaymentRequest: tt.args.payReq,
			})
			service.SaveInfo(pubkeyC, tt.args.toPeerPubKey, CtoA, AtoB, tt.args.amount, tt.args.payReq, t)

			initA = <-events

			if tt.args.task.PaymentRequest == "normal" {
				tt.args.task.PaymentRequest = initA.GetInitDoneType().PaymentRequest
			}
			if tt.args.newSid {
				sid = SwapID(uuid.NewString())
			}
			service.UpdatePayReq(tt.args.task.PaymentRequest)
			rebalancer.TaskSwap(sid, tt.args.task)
			select {
			case finalResponse, _ := <-rebalancer.events:
				assert.IsType(t, tt.expectedResponse.Response, finalResponse.Response)
				if finalResponse.GetErrorType() != nil { //error type
					assert.Contains(t, finalResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error, tt.expectedResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error)
					//assert.Equal(t, tt.expectedResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error, finalResponse.Response.(*protobuf.TaskResponse_ErrorType).ErrorType.Error)
				} else {
					assert.Equal(t, tt.expectedResponse.Response.(*protobuf.TaskResponse_PaymentInitiatedType).PaymentInitiatedType, finalResponse.GetPaymentInitiatedType())
				}
				//default:

			}

		})
	}

}
