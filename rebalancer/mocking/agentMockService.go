package mocking

import (
	"bytes"
	"context"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/offerm/lnagent/lightning"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
)

type agentMockService struct {
	baseMockService
	info *info
}

func NewAgentMockService() *agentMockService {
	return &agentMockService{info: &info{}}

}

func (ms *agentMockService) MakeHashPaymentAndMonitor(peerPubKey []byte, chanID uint64, hash []byte, payAddress []byte, amount uint64, cb lightning.PaymentCallBack) error {
	assert.Equal(ms.info.t, ms.info.toPeerPubKey, peerPubKey)
	assert.Equal(ms.info.t, ms.info.toChanID, chanID)
	assert.Equal(ms.info.t, ms.info.amount, amount)

	//for testing the MakeHashPaymentAndMonitor error response
	x, _ := hex.DecodeString(makeHashPaymentTestPubKey)
	if bytes.Equal(peerPubKey, x) {
		return lnwire.NewError()
	}
	return nil
}

func (ms *agentMockService) DecodePayReq(payReqString *lnrpc.PayReqString) (*lnrpc.PayReq, error) {
	assert.Equal(ms.info.t, ms.info.payReq, payReqString.PayReq)
	switch payReqString.PayReq {
	//for testing the decodePayReq error response
	case invalidPayReq: //invalid payReq test
		return nil, lnwire.NewError()

	//for testing the invalidRole response, mock returns a "valid" pay req so the swap task will work
	case invalidRolePayReq:
		return &lnrpc.PayReq{
			PaymentHash: payHash,
			PaymentAddr: []byte{},
			NumMsat:     100000,
		}, nil

	//for testing the pay hash response
	case swapCompareTestPayReq:
		return &lnrpc.PayReq{
			PaymentHash: payHash,
			PaymentAddr: []byte{},
			NumMsat:     1000000,
		}, nil

	//for testing the hex decode error response
	case decodeTestPayReq:
		return &lnrpc.PayReq{
			PaymentHash: "g",
			PaymentAddr: []byte{},
			NumMsat:     1000000,
		}, nil

	// a "valid" PayReq for the next tests, not containing true information
	case valid:
		return &lnrpc.PayReq{
			PaymentHash: hex.EncodeToString(ms.info.swapHash), // the hash was saved at the newHoldInvoice func
			PaymentAddr: []byte{},
			NumMsat:     1000000,
		}, nil

	default:
		return nil, nil

	}
}

func (ms *agentMockService) NewHoldInvoice(hash []byte, amount uint64, swapID string, cb lightning.InvoiceCallBack) (*invoicesrpc.AddHoldInvoiceResp, error) {
	// saving the hash to be used in the decodePayReq func
	ms.info.swapHash = hash
	payReq := uuid.New().String()
	payAddress := uuid.New()
	bytesAddress := []byte(payAddress[:])

	return &invoicesrpc.AddHoldInvoiceResp{
		PaymentRequest: payReq,
		PaymentAddr:    bytesAddress,
	}, nil
}

func (ms *agentMockService) SaveInfo(fromPeerPubKey []byte, toPeerPubKey []byte, fromChanID uint64, toChanID uint64, amount uint64, payReq string, t *testing.T) {
	ms.info = &info{
		fromPeerPubKey: fromPeerPubKey,
		toPeerPubKey:   toPeerPubKey,
		fromChanID:     fromChanID,
		toChanID:       toChanID,
		amount:         amount,
		payReq:         payReq,
		t:              t,
	}
}

func (ms *agentMockService) UpdatePayReq(payReq string) {
	ms.info.payReq = payReq
}

func (ms *agentMockService) GetInfo(ctx context.Context, request *lnrpc.GetInfoRequest) (*lnrpc.GetInfoResponse, error) {
	return &lnrpc.GetInfoResponse{IdentityPubkey: ""}, nil
}

func (ms *agentMockService) ListChannels(ctx context.Context, request *lnrpc.ListChannelsRequest) (*lnrpc.ListChannelsResponse, error) {
	return &lnrpc.ListChannelsResponse{}, nil
}

func (ms *agentMockService) FeeReport(ctx context.Context, request *lnrpc.FeeReportRequest, opts ...grpc.CallOption) (*lnrpc.FeeReportResponse, error) {
	return &lnrpc.FeeReportResponse{}, nil
}
func (ms *agentMockService) SignMessage(ctx context.Context, request *lnrpc.SignMessageRequest, opts ...grpc.CallOption) (*lnrpc.SignMessageResponse, error) {
	return &lnrpc.SignMessageResponse{Signature: "bla"}, nil
}

func (ms *agentMockService) Close() {
}
