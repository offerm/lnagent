package mocking

import (
	"bytes"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/offerm/lnagent/lightning"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestMockService struct {
	baseMockService
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

func (ms *TestMockService) MakeHashPaymentAndMonitor(peerPubKey []byte, chanID uint64, hash []byte, payAddress []byte, amount uint64, cb lightning.PaymentCallBack) error {
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
		return nil, nil

	}
}

func (ms *TestMockService) NewHoldInvoice(hash []byte, amount uint64, swapID string, cb lightning.InvoiceCallBack) (*invoicesrpc.AddHoldInvoiceResp, error) {
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
