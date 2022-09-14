package mocking

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/offerm/lnagent/lightning"
)

type mockService struct {
	baseMockService
}

type payData struct {
	hash       string
	cb         lightning.InvoiceCallBack
	payAddress []byte
	memo       string
}

var (
	payReqToPayData     map[string]*payData
	payAddressToPayData map[string]*payData
	callBackStack       []lightning.PaymentCallBack
)

func init() {
	payReqToPayData = make(map[string]*payData)
	payAddressToPayData = make(map[string]*payData)
}

func (ms *mockService) DecodePayReq(payReqString *lnrpc.PayReqString) (*lnrpc.PayReq, error) {
	pd := payReqToPayData[payReqString.PayReq]
	return &lnrpc.PayReq{
		PaymentHash: pd.hash,
		PaymentAddr: pd.payAddress,
		NumMsat:     1000000,
	}, nil
}

func (ms *mockService) NewHoldInvoice(hash []byte, amount uint64, swapID string, cb lightning.InvoiceCallBack) (*invoicesrpc.AddHoldInvoiceResp, error) {
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

func (ms *mockService) MakeHashPaymentAndMonitor(peerPubKey []byte, chanID uint64, hash []byte, payAddress []byte, amount uint64, cb lightning.PaymentCallBack) error {
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

func NewMockService(config *lightning.Config) lightning.Service {
	return &mockService{}
}
