package lightning

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lnwire"
	"google.golang.org/grpc"
)

type mockService struct {
}

type payData struct {
	hash       string
	cb         InvoiceCallBack
	payAddress []byte
	memo       string
}

var (
	payReqToPayData     map[string]*payData
	payAddressToPayData map[string]*payData
	callBackStack       []PaymentCallBack
)

func init() {
	payReqToPayData = make(map[string]*payData)
	payAddressToPayData = make(map[string]*payData)
}

func (ms *mockService) DecodePayReq(payReqString *lnrpc.PayReqString) (*lnrpc.PayReq, error) {
	if payReqString.PayReq == "1111" {
		return nil, lnwire.NewError()
	}
	pd := payReqToPayData[payReqString.PayReq]
	if payReqString.PayReq == "1234" { //pay hash test
		return &lnrpc.PayReq{
			PaymentHash: "1234",
			PaymentAddr: []byte{},
			NumMsat:     1000000,
		}, nil
	} else if payReqString.PayReq == "5678" { //hex decode test
		return &lnrpc.PayReq{
			PaymentHash: "g",
			PaymentAddr: []byte{},
			NumMsat:     1000000,
		}, nil
	} else if payReqString.PayReq == "2222" {
		x, _ := hex.DecodeString("1234")
		return &lnrpc.PayReq{
			PaymentHash: "1234",
			PaymentAddr: x,
			NumMsat:     1000000,
		}, nil

	}
	return &lnrpc.PayReq{
		PaymentHash: pd.hash,
		PaymentAddr: pd.payAddress,
		NumMsat:     1000000,
	}, nil
}

func (ms *mockService) NewHoldInvoice(hash []byte, amount uint64, swapID string, cb InvoiceCallBack) (*invoicesrpc.AddHoldInvoiceResp, error) {
	payReq := uuid.New().String()
	payAddress := uuid.New()
	bytesAddress := []byte(payAddress[:])

	currData := payData{
		hash:       hex.EncodeToString(hash[:]),
		cb:         cb,
		payAddress: bytesAddress,
		memo:       swapID,
	}
	payReqToPayData[payReq] = &currData
	payAddressToPayData[hex.EncodeToString(payReqToPayData[payReq].payAddress)] = &currData
	payAddressToPayData[payAddress.String()] = payReqToPayData[payReq]

	return &invoicesrpc.AddHoldInvoiceResp{
		PaymentRequest: payReq,
		PaymentAddr:    bytesAddress,
	}, nil
}

func (ms *mockService) MakeHashPaymentAndMonitor(peerPubKey []byte, chanID uint64, hash []byte, payAddress []byte, amount uint64, cb PaymentCallBack) error {
	x, _ := hex.DecodeString("1234")
	if bytes.Equal(peerPubKey, x) {
		return lnwire.NewError()
	}
	callBackStack = append(callBackStack, cb)

	payAddr := hex.EncodeToString(payAddress)
	pd := payAddressToPayData[payAddr]
	go pd.cb(&lnrpc.Invoice{
		State: lnrpc.Invoice_ACCEPTED,
		Memo:  pd.memo,
		RHash: hash,
	})
	return nil
}

func (ms *mockService) SettleInvoice(msg *invoicesrpc.SettleInvoiceMsg) (*invoicesrpc.SettleInvoiceResp, error) {
	cb := callBackStack[len(callBackStack)-1]
	callBackStack = callBackStack[:len(callBackStack)-1]
	payHash := sha256.Sum256(msg.Preimage[:])

	go cb(&lnrpc.Payment{
		Status:          lnrpc.Payment_SUCCEEDED,
		PaymentPreimage: hex.EncodeToString(msg.Preimage),
		PaymentHash:     hex.EncodeToString(payHash[:]),
	})
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
