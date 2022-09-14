package lightning

import (
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestInitMockService struct {
	mockService
	info *info
}

const (
	invalidDecodePayHash = "g"
	amount               = 1000000
	normalPayReq         = "1234"
	testSID              = "5656"
)

func (ms *TestInitMockService) DecodePayReq(payReqString *lnrpc.PayReqString) (*lnrpc.PayReq, error) {
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
		pd := payReqToPayData[payReqString.PayReq]
		return &lnrpc.PayReq{
			PaymentHash: pd.hash,
			PaymentAddr: pd.payAddress,
			NumMsat:     amount,
		}, nil
	}
}

func (ms *TestInitMockService) NewHoldInvoice(hash []byte, amount uint64, swapID string, cb InvoiceCallBack) (*invoicesrpc.AddHoldInvoiceResp, error) {
	// saving the hash to be used in the decodePayReq func
	ms.info.swapHash = hash
	payReq := normalPayReq

	payAddress := uuid.New()
	bytesAddress := []byte(payAddress[:])
	if swapID == testSID {
		return nil, lnwire.NewError()
	}

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

func (ms *TestInitMockService) SaveInfo(fromPeerPubKey []byte, toPeerPubKey []byte, fromChanID uint64, toChanID uint64, amount uint64, payReq string, t *testing.T) {
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
