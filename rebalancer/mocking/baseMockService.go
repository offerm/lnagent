package mocking

import (
	"context"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/offerm/lnagent/lightning"
	"google.golang.org/grpc"
)

type baseMockService struct {
}

func (ms *baseMockService) DecodePayReq(payReqString *lnrpc.PayReqString) (*lnrpc.PayReq, error) {
	panic("not supposed to use this func")
}

func (ms *baseMockService) NewHoldInvoice(hash []byte, amount uint64, swapID string, cb lightning.InvoiceCallBack) (*invoicesrpc.AddHoldInvoiceResp, error) {
	panic("not supposed to use this func")
}

func (ms *baseMockService) MakeHashPaymentAndMonitor(peerPubKey []byte, chanID uint64, hash []byte, payAddress []byte, amount uint64, cb lightning.PaymentCallBack) error {
	panic("not supposed to use this func")
}

func (ms *baseMockService) SettleInvoice(msg *invoicesrpc.SettleInvoiceMsg) (*invoicesrpc.SettleInvoiceResp, error) {
	panic("not supposed to use this func")
}

func (ms *baseMockService) GetInfo(ctx context.Context, request *lnrpc.GetInfoRequest) (*lnrpc.GetInfoResponse, error) {
	panic("not supposed to use this func")
}

func (ms *baseMockService) ListChannels(ctx context.Context, request *lnrpc.ListChannelsRequest) (*lnrpc.ListChannelsResponse, error) {
	panic("not supposed to use this func")
}

func (ms *baseMockService) FeeReport(ctx context.Context, request *lnrpc.FeeReportRequest, opts ...grpc.CallOption) (*lnrpc.FeeReportResponse, error) {
	panic("not supposed to use this func")
}

func (ms *baseMockService) SignMessage(ctx context.Context, request *lnrpc.SignMessageRequest, opts ...grpc.CallOption) (*lnrpc.SignMessageResponse, error) {
	panic("not supposed to use this func")
}

func (ms *baseMockService) ChanInfo(ctx context.Context, request *lnrpc.ChanInfoRequest) (*lnrpc.ChannelEdge, error) {
	panic("not supposed to use this func")
}

func (ms *baseMockService) DescribeGraph(ctx context.Context, request *lnrpc.ChannelGraphRequest, opts ...grpc.CallOption) (*lnrpc.ChannelGraph, error) {
	panic("not supposed to use this func")
}
func (ms *baseMockService) Close() {
	panic("not supposed to use this func")
}
