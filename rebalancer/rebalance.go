package rebalancer

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/offerm/lnagent/lightning"
	"github.com/offerm/lnagent/protobuf"
	log "github.com/sirupsen/logrus"
)

// TODO: use one content per swap (with cancel)

type Rebalancer struct {
	pubkey    string
	swaps     *swaps
	events    chan *protobuf.TaskResponse
	lnservice lightning.Service
}

func NewRebalancer(events chan *protobuf.TaskResponse, lnservice lightning.Service) *Rebalancer {
	r := &Rebalancer{
		swaps:     NewSwaps(),
		events:    events,
		lnservice: lnservice,
	}
	return r

}

func (rebalancer *Rebalancer) TaskInit(swapID SwapID, init *protobuf.Task_Init) {
	var err error

	// make sure we don't have an active swap with this ID
	swap := rebalancer.swaps.Get(swapID)
	if swap != nil {
		rebalancer.errorTaskResponse(swap, fmt.Sprintf("swap %v is already active", swapID))
		return
	}
	/*
		checks to add:
		if initiator - no payment request, if not - valid payment request
		validate channels (in and out)
		validate that these channels should be rebalanced

	*/

	// create a new swap for this task
	swap = NewSwap(init, swapID)

	switch init.Role {
	case protobuf.Task_INITIATOR:
		log.Info("time to initiate a swap")

		if init.PaymentRequest != "" {
			rebalancer.errorTaskResponse(swap, "initiator should not accept a payment request ")
			return
		}
		// make a secret for this rebalance. Only the initiator knows it.
		// Keep it safe until payment from "from" node is locked

		if _, err := rand.Read(swap.secret[:]); err != nil {
			log.Fatal(err)
		}
		swap.hash = sha256.Sum256(swap.secret[:])

	// Follower accepts a payment request and extract the payment hash from it
	case protobuf.Task_FOLLOWER:
		swap.decodedPayReq, err = rebalancer.lnservice.DecodePayReq(&lnrpc.PayReqString{
			PayReq: init.PaymentRequest,
		})
		if err != nil {
			rebalancer.errorTaskResponse(swap, fmt.Sprintf("failed to decode payreq %v - %v", init.PaymentRequest, err))
			return
		}
		h, err := hex.DecodeString(swap.decodedPayReq.PaymentHash)
		if err != nil {
			rebalancer.errorTaskResponse(swap, fmt.Sprintf("invalid payment request hash %v", swap.decodedPayReq.PaymentHash))
			return
		}
		copy(swap.hash[:], h)

	default:
		rebalancer.errorTaskResponse(swap, fmt.Sprintf("unknown rolee %v for swap %v ", init.Role.String(), swapID))
		return
	}

	// create a hold invoice for the payment
	swap.holdInvoice, err = rebalancer.lnservice.NewHoldInvoice(swap.hash[:], init.From.AmtMsat, string(swap.swapID), rebalancer.onInvoice)
	if err != nil {
		rebalancer.errorTaskResponse(swap, fmt.Sprintf("can't create invoice for swap %v - %v", swapID, err))
		return
	}

	rebalancer.swaps.Add(swap)
	rebalancer.initDoneResponse(swap)
	return
}

func (rebalancer *Rebalancer) TaskSwap(swapID SwapID, payment *protobuf.Task_Swap) {
	// make sure we don't have an active swap with this ID
	swap := rebalancer.swaps.Get(swapID)
	if swap == nil {
		rebalancer.errorTaskResponse(swap, fmt.Sprintf("swap %v does not exist", swapID))
		return
	}

	if swap.initTask == nil || swap.initTask.Role != protobuf.Task_INITIATOR {
		rebalancer.errorTaskResponse(swap, "swap can't start by non initiator")
		return
	}

	var err error
	swap.decodedPayReq, err = rebalancer.lnservice.DecodePayReq(&lnrpc.PayReqString{
		PayReq: payment.PaymentRequest,
	})
	if err != nil {
		rebalancer.errorTaskResponse(swap, fmt.Sprintf("failed to decode payreq %v - %v", payment.PaymentRequest, err))
		return
	}

	h, err := hex.DecodeString(swap.decodedPayReq.PaymentHash)
	if err != nil {
		rebalancer.errorTaskResponse(swap, fmt.Sprintf("invalid payment request hash %v", swap.decodedPayReq.PaymentHash))
		return
	}

	if bytes.Compare(h, swap.hash[:]) != 0 {
		rebalancer.errorTaskResponse(swap, fmt.Sprintf("wrong hash %v received for payment %v", swap.decodedPayReq.PaymentHash, swapID))
		return
	}

	to := swap.initTask.To
	if swap.decodedPayReq.NumMsat != int64(to.AmtMsat) {
		rebalancer.errorTaskResponse(swap, fmt.Sprintf("different amounts. expected %v invoiced %v", to.AmtMsat, swap.decodedPayReq.NumMsat))
		return
	}

	//log.Info("init loop closed, time to start the payment")
	err = rebalancer.lnservice.MakeHashPaymentAndMonitor(to.PeerPubKey, to.ChanId, swap.hash[:], swap.decodedPayReq.PaymentAddr,
		to.AmtMsat, rebalancer.onPayment)
	if err != nil {
		rebalancer.errorTaskResponse(swap, fmt.Sprintf("initiator failed to initiate payment - %v ", err))
		return
	}

	rebalancer.paymentInitiatedResponse(swap)
	return
}

func (rebalancer *Rebalancer) TaskCancel(_ SwapID, _ *protobuf.Task_Cancel) *protobuf.TaskResponse {
	return nil

}

func (rebalancer *Rebalancer) TaskUnKnow(_ SwapID) *protobuf.TaskResponse {
	return nil
}

func (rebalancer *Rebalancer) RebalanceOnPayment(payment *lnrpc.Payment) error {

	if payment.Status != lnrpc.Payment_SUCCEEDED && payment.Status != lnrpc.Payment_FAILED {
		return nil
	}

	//&& payment.PaymentHash == hex.EncodeToString(rebalancer.lnservice.rebalanceHash)
	hash, _ := hex.DecodeString(payment.PaymentHash)
	swap := rebalancer.swaps.SwapByHash(hash)

	// is swap can't be found, this payment is not related to an active swap
	if swap == nil {
		return nil
	}

	if payment.Status == lnrpc.Payment_FAILED {
		rebalancer.errorTaskResponse(swap, fmt.Sprintf("payment failed with reason %v", payment.FailureReason.String()))
		return nil
	}

	switch swap.initTask.Role {
	case protobuf.Task_FOLLOWER:
		preimage, err := hex.DecodeString(payment.PaymentPreimage)
		if err != nil {
			log.Error(err)
			return err
		}
		_, err = rebalancer.lnservice.SettleInvoice(&invoicesrpc.SettleInvoiceMsg{
			Preimage: preimage,
		})
		if err != nil {
			log.Error(err)
			return err
		}
		rebalancer.paymentSettledResponse(swap)

	case protobuf.Task_INITIATOR:
		rebalancer.swapDoneResponse(swap)
		log.Info("looks like we are done!!!!")
	}

	return nil
}

// rebalanceOnInvoice is called when on invoice event
func (rebalancer *Rebalancer) RebalanceOnInvoice(invoice *lnrpc.Invoice) error {

	// we are only looking for Invoice_ACCEPTED status which means that we locked a payment
	// from the peer and that we need to make our payment
	if invoice.State != lnrpc.Invoice_ACCEPTED {
		return nil
	}

	// if we can't find a swap based on Memo (that we created with the SwapID) -
	// it is not our invoice that was paid
	swap := rebalancer.swaps.Get(SwapID(invoice.Memo))
	if swap == nil {
		return nil
	}

	// TODO:
	if bytes.Compare(swap.hash[:], invoice.RHash) != 0 {
		return nil
		// todo: return with error
	}

	rebalancer.paymentLockedResponse(swap)
	switch swap.initTask.Role {
	case protobuf.Task_FOLLOWER:
		to := swap.initTask.To
		// TODO: check from
		if swap.decodedPayReq.NumMsat != int64(to.AmtMsat) {
			return fmt.Errorf("different amounts. expected %v invoiced %v", to.AmtMsat, swap.decodedPayReq.NumMsat)
		}

		err := rebalancer.lnservice.MakeHashPaymentAndMonitor(to.PeerPubKey, to.ChanId, swap.hash[:], swap.decodedPayReq.PaymentAddr, to.AmtMsat, rebalancer.onPayment)
		if err != nil {
			rebalancer.errorTaskResponse(swap, fmt.Sprintf("for error from MakeHashPaymentAndMonitor - %v", err))
			return err
		}
		rebalancer.paymentInitiatedResponse(swap)
	case protobuf.Task_INITIATOR:
		// if the initiator got paid it can settle the payment by sending the secret.
		// This exposes the secret to the payer which will do the same for his payer
		_, err := rebalancer.lnservice.SettleInvoice(&invoicesrpc.SettleInvoiceMsg{
			Preimage: swap.secret[:],
		})
		if err != nil {
			return err
		}
		rebalancer.paymentSettledResponse(swap)

	}
	return nil
}

func (rebalancer *Rebalancer) errorTaskResponse(swap *swap, error string) *protobuf.TaskResponse {
	event := &protobuf.TaskResponse{
		Swap_ID: string(swap.swapID),
		Response: &protobuf.TaskResponse_ErrorType{
			ErrorType: &protobuf.TaskResponse_Error{
				Error: error,
			},
		},
	}
	rebalancer.sendEvent(event)
	return event
}

func (rebalancer *Rebalancer) initDoneResponse(swap *swap) *protobuf.TaskResponse {
	event := &protobuf.TaskResponse{
		Swap_ID: string(swap.swapID),
		Response: &protobuf.TaskResponse_InitDoneType{
			InitDoneType: &protobuf.TaskResponse_Init_Done{
				PaymentRequest: swap.holdInvoice.PaymentRequest,
			},
		},
	}
	rebalancer.sendEvent(event)
	return event
}

func (rebalancer *Rebalancer) paymentInitiatedResponse(swap *swap) *protobuf.TaskResponse {
	event := &protobuf.TaskResponse{
		Swap_ID:  string(swap.swapID),
		Response: &protobuf.TaskResponse_PaymentInitiatedType{},
	}
	rebalancer.sendEvent(event)
	return event
}

func (rebalancer *Rebalancer) paymentLockedResponse(swap *swap) *protobuf.TaskResponse {
	event := &protobuf.TaskResponse{
		Swap_ID:  string(swap.swapID),
		Response: &protobuf.TaskResponse_PaymentLockedType{},
	}
	rebalancer.sendEvent(event)
	return event
}

func (rebalancer *Rebalancer) paymentSettledResponse(swap *swap) *protobuf.TaskResponse {
	event := &protobuf.TaskResponse{
		Swap_ID:  string(swap.swapID),
		Response: &protobuf.TaskResponse_PaymentSettledType{},
	}
	rebalancer.sendEvent(event)
	return event
}

func (rebalancer *Rebalancer) swapDoneResponse(swap *swap) *protobuf.TaskResponse {
	event := &protobuf.TaskResponse{
		Swap_ID: string(swap.swapID),
		Response: &protobuf.TaskResponse_SwapDoneType{
			SwapDoneType: &protobuf.TaskResponse_Swap_Done{},
		},
	}
	rebalancer.sendEvent(event)
	return event
}

func (rebalancer *Rebalancer) sendEvent(event *protobuf.TaskResponse) {
	log.Info(event)
	rebalancer.events <- event
}

func (rebalancer *Rebalancer) onInvoice(invoice *lnrpc.Invoice) error {
	if err := rebalancer.RebalanceOnInvoice(invoice); err != nil {
		return err
	}
	return nil
}

func (rebalancer *Rebalancer) onPayment(payment *lnrpc.Payment) error {
	if err := rebalancer.RebalanceOnPayment(payment); err != nil {
		return err
	}

	return nil
}
