package signutil

import (
	"github.com/fibercrypto/fibercryptowallet/src/core"
)

// NewTxnSignerIteratorFromMap allocate transaction signer iterator out of standard signers map
func NewTxnSignerIteratorFromMap(signers map[core.UID]core.TxnSigner) core.TxnSignerIterator {
	var signersSlice []core.TxnSigner
	for _, signer := range signers {
		signersSlice = append(signersSlice, signer)
	}
	return &DefaultTxnSignerIterator{
		signers: signersSlice,
		nextIdx: 0,
	}
}

// FilterSignersFromMap allocate iterator for transaction signers in a map matching custom condition
func FilterSignersFromMap(signers map[core.UID]core.TxnSigner, cond func(core.TxnSigner) bool) core.TxnSignerIterator {
	var signersSlice []core.TxnSigner
	for _, signer := range signers {
		if cond(signer) {
			signersSlice = append(signersSlice, signer)
		}
	}
	return &DefaultTxnSignerIterator{
		signers: signersSlice,
		nextIdx: 0,
	}
}

// DefaultTxnSignerIterator iterate over items in transaction signer slicer
type DefaultTxnSignerIterator struct {
	signers []core.TxnSigner
	nextIdx int
}

// Value of signer at iterator pointer position
func (sm *DefaultTxnSignerIterator) Value() core.TxnSigner {
	if sm.HasNext() {
		return sm.signers[sm.nextIdx]
	}
	return nil
}

// Next discards current value and moves iteration pointer up to next item
func (sm *DefaultTxnSignerIterator) Next() bool {
	sm.nextIdx++
	return sm.HasNext()
}

// HasNext may be used to query whether more items are to be expected in the sequence
func (sm *DefaultTxnSignerIterator) HasNext() bool {
	return sm.nextIdx < len(sm.signers)
}

// Count total number of items in sequence
func (sm *DefaultTxnSignerIterator) Count() int {
	return len(sm.signers)
}
