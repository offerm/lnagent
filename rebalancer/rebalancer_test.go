package rebalancer_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/offerm/lnagent/lightning"
	"github.com/offerm/lnagent/protobuf"
	"github.com/offerm/lnagent/rebalancer"

	//"github.com/offerm/lnagent/rebalancer"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
)

type mockService struct {
}

func (ms *mockService) DecodePayReq(*lnrpc.PayReqString) (*lnrpc.PayReq, error) {
	return nil, nil
}
func (ms *mockService) NewHoldInvoice([]byte, uint64, string, lightning.InvoiceCallBack) (*invoicesrpc.AddHoldInvoiceResp, error) {
	return nil, fmt.Errorf("sorry, can't do that")
}

func (ms *mockService) MakeHashPaymentAndMonitor([]byte, uint64, []byte, []byte, uint64, lightning.PaymentCallBack) error {
	return nil
}

func (ms *mockService) SettleInvoice(*invoicesrpc.SettleInvoiceMsg) (*invoicesrpc.SettleInvoiceResp, error) {
	return nil, nil
}

func (ms *mockService) GetInfo(ctx context.Context, request *lnrpc.GetInfoRequest) (*lnrpc.GetInfoResponse, error) {
	return nil, nil
}

func (ms *mockService) ListChannels(ctx context.Context, request *lnrpc.ListChannelsRequest) (*lnrpc.ListChannelsResponse, error) {
	return nil, nil
}

func (ms *mockService) FeeReport(ctx context.Context, request *lnrpc.FeeReportRequest, opts ...grpc.CallOption) (*lnrpc.FeeReportResponse, error) {
	return nil, nil
}

func (ms *mockService) SignMessage(ctx context.Context, request *lnrpc.SignMessageRequest, opts ...grpc.CallOption) (*lnrpc.SignMessageResponse, error) {
	return nil, nil
}

func (ms *mockService) ChanInfo(ctx context.Context, request *lnrpc.ChanInfoRequest) (*lnrpc.ChannelEdge, error) {
	return nil, nil
}

func (ms *mockService) DescribeGraph(ctx context.Context, request *lnrpc.ChannelGraphRequest, opts ...grpc.CallOption) (*lnrpc.ChannelGraph, error) {

	return nil, nil
}
func (ms *mockService) Close() {
	return
}

func TestRebalancer_TaskInit(t *testing.T) {
	pubkeyB, _ := hex.DecodeString("02b998d8c3f065f3e0a8b383bd00dff56aeeac05c52ea2b7a5c936ff8ab2fb369a")
	pubkeyC, _ := hex.DecodeString("02848fffeb2ebaafdcd6b795b3a45d1e2397181e1c0d4424e86661276bfbe815a9")
	AtoB := uint64(2542708600164515840)
	CtoA := uint64(2542877924953358337)
	amount := uint64(1000 * 1000)

	events := make(chan *protobuf.TaskResponse, 1)
	rebalancerA := rebalancer.NewRebalancer(events, &mockService{})

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
	assert.IsType(t, &protobuf.TaskResponse_ErrorType{}, initA.Response)
	assert.Equal(t, initA.Swap_ID, string(sid))
}
