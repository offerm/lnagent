package rebalancer

import (
	"bytes"
	"crypto/sha256"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/offerm/lnagent/protobuf"
	"sync"
)

type SwapID string

// TODO: add swap state to prevent double execution

type swap struct {
	initTask      *protobuf.Task_Init
	swapID        SwapID
	secret        [sha256.Size]byte
	hash          [sha256.Size]byte
	rebalanceHash []byte
	holdInvoice   *invoicesrpc.AddHoldInvoiceResp
	decodedPayReq *lnrpc.PayReq
}

type swaps struct {
	lock        sync.Mutex
	activeSwaps map[SwapID]*swap
}

func NewSwap(init *protobuf.Task_Init, id SwapID) *swap {
	return &swap{
		initTask: init,
		swapID:   id,
	}
}

func NewSwaps() *swaps {
	return &swaps{
		lock:        sync.Mutex{},
		activeSwaps: make(map[SwapID]*swap),
	}
}

func (s *swaps) Get(id SwapID) *swap {
	s.lock.Lock()
	defer s.lock.Unlock()

	// make sure we don't have an active swap with this ID
	swap, ok := s.activeSwaps[id]
	if !ok {
		return nil
	}
	return swap
}

func (s *swaps) Add(swap *swap) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.activeSwaps[swap.swapID] = swap
}

func (s *swaps) SwapByHash(hash []byte) *swap {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, swap := range s.activeSwaps {
		if bytes.Compare(swap.hash[:], hash) == 0 {
			return swap
		}
	}
	return nil
}
