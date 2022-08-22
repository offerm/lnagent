package lightning

import (
	"context"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"google.golang.org/grpc"
)

type mockService struct {
}

func (ms *mockService) DecodePayReq(*lnrpc.PayReqString) (*lnrpc.PayReq, error) {
	return nil, nil
}
func (ms *mockService) NewHoldInvoice([]byte, uint64, string, InvoiceCallBack) (*invoicesrpc.AddHoldInvoiceResp, error) {
	return nil, fmt.Errorf("sorry, can't do that")
}

func (ms *mockService) MakeHashPaymentAndMonitor([]byte, uint64, []byte, []byte, uint64, PaymentCallBack) error {
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

func NewMockService(config *Config) Service {
	return &mockService{}
}
