package golang

import (
	"fmt"
	"time"

	"github.com/ellcrys/util"
	"github.com/ncodes/cocoon/core/common"
	"github.com/ncodes/cocoon/core/connector/server/proto_connector"
	"github.com/ncodes/cocoon/core/types"
)

// LockCB represents the type of function to pass to the Lock functions.
// The callback will be passed the folling methods:
// isAcquirer: Checks if the current lock session still has the lock.
// release: Used to release the lock
// refresh: Used to refresh the lock
type LockCB func(isAcquirer func() bool, release, refresh func() error)

// Link provides access to all platform services available to
// the cocoon code.
type Link struct {
	cocoonID string
	native   bool
}

// NewLink creates a new link to a cocoon
func NewLink(cocoonID string) *Link {
	return &Link{
		cocoonID: cocoonID,
	}
}

// NewNativeLink create a new native link to a cocoon
func newNativeLink(cocoonID string) *Link {
	return &Link{
		cocoonID: cocoonID,
		native:   true,
	}
}

// IsNative checks whether the link is a native link
func (link *Link) IsNative() bool {
	return link.native
}

// GetCocoonID returns the cocoon id attached to this link
func (link *Link) GetCocoonID() string {
	return link.cocoonID
}

// NewRangeGetter creates an instance of a RangeGetter for a specified ledger.
func (link *Link) NewRangeGetter(ledgerName, start, end string, inclusive bool) *RangeGetter {
	return NewRangeGetter(ledgerName, link.GetCocoonID(), start, end, inclusive)
}

// CreateLedger creates a new ledger by sending an
// invoke transaction (TxCreateLedger) to the connector.
// If chained is set to true, a blockchain is created and subsequent
// PUT operations to the ledger will be included in the types. Otherwise,
// PUT operations will only be included in the types.
func (link *Link) CreateLedger(name string, chained, public bool) (*types.Ledger, error) {

	if !common.IsValidResName(name) {
		return nil, types.ErrInvalidResourceName
	}

	result, err := sendLedgerOp(&proto_connector.LedgerOperation{
		ID:     util.UUID4(),
		Name:   types.TxCreateLedger,
		LinkTo: link.GetCocoonID(),
		Params: []string{name, fmt.Sprintf("%t", chained), fmt.Sprintf("%t", public)},
	})

	if err != nil {
		return nil, err
	}

	var ledger types.Ledger
	if err = util.FromJSON(result, &ledger); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data")
	}

	return &ledger, nil
}

// GetLedger fetches a ledger
func (link *Link) GetLedger(ledgerName string) (*types.Ledger, error) {

	result, err := sendLedgerOp(&proto_connector.LedgerOperation{
		ID:     util.UUID4(),
		Name:   types.TxGetLedger,
		LinkTo: link.GetCocoonID(),
		Params: []string{ledgerName},
	})

	if err != nil {
		return nil, err
	}

	var ledger types.Ledger
	if err = util.FromJSON(result, &ledger); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data")
	}

	return &ledger, nil
}

// Put puts a transaction in a ledger
func (link *Link) Put(ledgerName string, key string, value []byte) (*types.Transaction, error) {

	start := time.Now()

	if !common.IsValidResName(key) {
		return nil, types.ErrInvalidResourceName
	}

	ledger, err := link.GetLedger(ledgerName)
	if err != nil {
		return nil, err
	}

	tx := &types.Transaction{
		ID:             util.UUID4(),
		Ledger:         ledger.Name,
		LedgerInternal: types.MakeLedgerName(link.GetCocoonID(), ledger.Name),
		Key:            key,
		KeyInternal:    types.MakeTxKey(link.GetCocoonID(), key),
		Value:          string(value),
		CreatedAt:      time.Now().Unix(),
	}

	tx.Hash = tx.MakeHash()

	if ledger.Chained {
		respChan := make(chan interface{})
		blockMaker.Add(&Entry{
			Tx:       tx,
			RespChan: respChan,
			LinkTo:   link.GetCocoonID(),
		})
		result := <-respChan

		log.Debug("Put(): Time taken: ", time.Since(start))

		switch v := result.(type) {
		case error:
			return nil, v
		case *types.Block:
			tx.Block = v
			return tx, err
		default:
			return nil, fmt.Errorf("unexpected response %s", err)
		}
	}

	txJSON, _ := util.ToJSON([]*types.Transaction{tx})
	putTxResultBs, err := sendLedgerOp(&proto_connector.LedgerOperation{
		ID:     util.UUID4(),
		Name:   types.TxPut,
		LinkTo: link.GetCocoonID(),
		Params: []string{ledgerName},
		Body:   txJSON,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to put transaction: %s", err)
	}

	var putTxResult types.PutResult
	if err := util.FromJSON(putTxResultBs, &putTxResult); err != nil {
		return nil, common.JSONCoerceErr("putTxResult", err)
	}

	// ensure transaction was successful
	for _, txResult := range putTxResult.TxReceipts {
		if txResult.ID == tx.ID && len(txResult.Err) > 0 {
			return nil, fmt.Errorf("failed to put transaction: %s", txResult.Err)
		}
	}

	log.Debug("Put(): Time taken: ", time.Since(start))

	return tx, nil
}

// Get gets a transaction from a ledger
func (link *Link) Get(ledgerName, key string) (*types.Transaction, error) {

	result, err := sendLedgerOp(&proto_connector.LedgerOperation{
		ID:     util.UUID4(),
		Name:   types.TxGet,
		LinkTo: link.GetCocoonID(),
		Params: []string{ledgerName, key},
	})

	if err != nil {
		return nil, err
	}

	var tx types.Transaction
	if err = util.FromJSON(result, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data")
	}

	if tx.Block.ID == "" {
		tx.Block = nil
	}

	return &tx, nil
}

// GetBlock gets a block from a ledger by a block id
func (link *Link) GetBlock(ledgerName, id string) (*types.Block, error) {

	result, err := sendLedgerOp(&proto_connector.LedgerOperation{
		ID:     util.UUID4(),
		Name:   types.TxGetBlockByID,
		LinkTo: link.GetCocoonID(),
		Params: []string{ledgerName, id},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get block: %s", err)
	}

	var blk types.Block
	if err = util.FromJSON(result, &blk); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data")
	}

	return &blk, nil
}

// Lock acquires a lock on the specified key. The onAcquired method is called
// when the lock has been acquired. Use the ttl to decide how long the lock
// will be held for. Minimum ttl is 10 seconds and max is 30 minutes.
// If the link is pointed to a external cocoon, the lock will also be enforced
// in the associated cocoon.
func (link *Link) Lock(key string, ttl int, onAcquired LockCB) {

}
