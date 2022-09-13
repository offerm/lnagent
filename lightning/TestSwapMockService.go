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
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
)

type TestMockService struct {
	info *info
}

type info struct {
	fromPeerPubKey []byte
	toPeerPubKey   []byte
	fromChanID     uint64
	toChanID       uint64
	amount         uint64
	swapHash       []byte
	t              *testing.T
	payReq         string
}

const (
	invalidPayReq             = "invalidPayReq"
	invalidRolePayReq         = "invalid role"
	swapCompareTestPayReq     = "swapCompareTestPayReq"
	decodeTestPayReq          = "decode"
	payHash                   = "1234"
	makeHashPaymentTestPubKey = "123"
	valid                     = "validPayReq"
)

func (ms *TestMockService) MakeHashPaymentAndMonitor(peerPubKey []byte, chanID uint64, hash []byte, payAddress []byte, amount uint64, cb PaymentCallBack) error {
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

func (ms *TestMockService) DecodePayReq(payReqString *lnrpc.PayReqString) (*lnrpc.PayReq, error) {
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
		pd := payReqToPayData[payReqString.PayReq]
		return &lnrpc.PayReq{
			PaymentHash: pd.hash,
			PaymentAddr: pd.payAddress,
			NumMsat:     1000000,
		}, nil

	}
}

func (ms *TestMockService) NewHoldInvoice(hash []byte, amount uint64, swapID string, cb InvoiceCallBack) (*invoicesrpc.AddHoldInvoiceResp, error) {
	// saving the hash to be used in the decodePayReq func
	ms.info.swapHash = hash
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

func (ms *TestMockService) SaveInfo(fromPeerPubKey []byte, toPeerPubKey []byte, fromChanID uint64, toChanID uint64, amount uint64, payReq string, t *testing.T) {
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

func (ms *TestMockService) UpdatePayReq(payReq string) {
	ms.info.payReq = payReq
}

func (ms *TestMockService) SettleInvoice(msg *invoicesrpc.SettleInvoiceMsg) (*invoicesrpc.SettleInvoiceResp, error) {
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

func (ms *TestMockService) GetInfo(ctx context.Context, request *lnrpc.GetInfoRequest) (*lnrpc.GetInfoResponse, error) {
	panic("not supposed to use this func")
	return nil, nil
}

func (ms *TestMockService) ListChannels(ctx context.Context, request *lnrpc.ListChannelsRequest) (*lnrpc.ListChannelsResponse, error) {
	panic("not supposed to use this func")
	return nil, nil
}

func (ms *TestMockService) FeeReport(ctx context.Context, request *lnrpc.FeeReportRequest, opts ...grpc.CallOption) (*lnrpc.FeeReportResponse, error) {
	panic("not supposed to use this func")
	return nil, nil
}

func (ms *TestMockService) SignMessage(ctx context.Context, request *lnrpc.SignMessageRequest, opts ...grpc.CallOption) (*lnrpc.SignMessageResponse, error) {
	panic("not supposed to use this func")
	return nil, nil
}

func (ms *TestMockService) ChanInfo(ctx context.Context, request *lnrpc.ChanInfoRequest) (*lnrpc.ChannelEdge, error) {
	panic("not supposed to use this func")
	return nil, nil
}

func (ms *TestMockService) DescribeGraph(ctx context.Context, request *lnrpc.ChannelGraphRequest, opts ...grpc.CallOption) (*lnrpc.ChannelGraph, error) {
	panic("not supposed to use this func")
	return nil, nil
}
func (ms *TestMockService) Close() {
	panic("not supposed to use this func")
	return
}
