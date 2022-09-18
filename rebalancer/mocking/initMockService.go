package mocking

import (
	"github.com/google/uuid"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/offerm/lnagent/lightning"
	"github.com/stretchr/testify/assert"
	"testing"
)

type InitMockService struct {
	baseMockService
	info *info
}

const (
	invalidDecodePayHash = "g"
	amount               = 1000000
	normalPayReq         = "1234"
	testSID              = "5656"
)

func (ms *InitMockService) DecodePayReq(payReqString *lnrpc.PayReqString) (*lnrpc.PayReq, error) {
	assert.Equal(ms.info.t, ms.info.payReq, payReqString.PayReq)

	switch payReqString.PayReq {
	//for testing the decodePayReq error response
	case invalidPayReq:
		return nil, lnwire.NewError()

	//for testing the hex decode error response
	case decodeTestPayReq:
		return &lnrpc.PayReq{
			PaymentHash: invalidDecodePayHash,
			PaymentAddr: []byte{},
			NumMsat:     amount,
		}, nil
	default:
		return nil, nil
	}
}

func (ms *InitMockService) NewHoldInvoice(hash []byte, amount uint64, swapID string, cb lightning.InvoiceCallBack) (*invoicesrpc.AddHoldInvoiceResp, error) {
	// saving the hash to be used in the decodePayReq func
	ms.info.swapHash = hash
	payReq := normalPayReq

	payAddress := uuid.New()
	bytesAddress := []byte(payAddress[:])
	if swapID == testSID {
		return nil, lnwire.NewError()
	}

	return &invoicesrpc.AddHoldInvoiceResp{
		PaymentRequest: payReq,
		PaymentAddr:    bytesAddress,
	}, nil
}

func (ms *InitMockService) SaveInfo(fromPeerPubKey []byte, toPeerPubKey []byte, fromChanID uint64, toChanID uint64, amount uint64, payReq string, t *testing.T) {
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
