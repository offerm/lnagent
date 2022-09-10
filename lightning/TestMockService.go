package lightning

import (
	"bytes"
	"encoding/hex"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestMockService struct {
	mockService
	info *info
}

type info struct {
	fromPeerPubKey []byte
	toPeerPubKey   []byte
	fromChanID     uint64
	toChanID       uint64
	amount         uint64
	swapHash       [32]uint8
	t              *testing.T
	payReq         string
}

const (
	invalidPayReq      = "invalidPayReq"
	invalidRole        = "invalid role"
	swapCompareTest    = "swapCompareTest"
	decode             = "decode"
	payHash            = "1234"
	makeHashPaymenTest = "123"
)

func (ms *TestMockService) MakeHashPaymentAndMonitor(peerPubKey []byte, chanID uint64, hash []byte, payAddress []byte, amount uint64, cb PaymentCallBack) error {
	assert.Equal(ms.info.t, ms.info.toPeerPubKey, peerPubKey)
	assert.Equal(ms.info.t, ms.info.toChanID, chanID)
	assert.Equal(ms.info.t, ms.info.amount, amount)

	x, _ := hex.DecodeString(makeHashPaymenTest)
	if bytes.Equal(peerPubKey, x) {
		return lnwire.NewError()
	}
	return nil
}

func (ms *TestMockService) DecodePayReq(payReqString *lnrpc.PayReqString) (*lnrpc.PayReq, error) {
	assert.Equal(ms.info.t, ms.info.payReq, payReqString.PayReq)
	switch payReqString.PayReq {
	case invalidPayReq: //invalid payReq test
		return nil, lnwire.NewError()

	case invalidRole:
		return &lnrpc.PayReq{
			PaymentHash: payHash,
			PaymentAddr: []byte{},
			NumMsat:     100000,
		}, nil

	case swapCompareTest: //pay hash test
		return &lnrpc.PayReq{
			PaymentHash: payHash,
			PaymentAddr: []byte{},
			NumMsat:     1000000,
		}, nil

	case decode:
		return &lnrpc.PayReq{
			PaymentHash: "g",
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
