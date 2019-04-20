package main

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// PendingBlockHeader holds pending block header.
type PendingBlockHeader struct {
	Time *hexutil.Uint64 `json:"timestamp"`
}

// PendingBlockTransactions holds pending block transactions.
type PendingBlockTransactions struct {
	PendingBlockHeader
	Transactions []*types.Transaction
}

// PendingTransactionsQuery contains query parameters.
type PendingTransactionsQuery struct {
	FilterID string
	From     *common.Address
}

// PendingTransactionsResponse holds response data.
type PendingTransactionsResponse struct {
	FilterID     string         `json:"filterId"`
	Transactions []*Transaction `json:"transactions"`
}

// Transaction wrapper around tx.
type Transaction struct {
	*types.Transaction
	From      *common.Address
	Timestamp *uint64
}

// MarshalJSON encodes transaction wrappper.
func (t *Transaction) MarshalJSON() ([]byte, error) {
	type txdata struct {
		AccountNonce hexutil.Uint64  `json:"nonce"`
		Price        *hexutil.Big    `json:"gasPrice"`
		GasLimit     hexutil.Uint64  `json:"gas"`
		Recipient    *common.Address `json:"to"`
		From         *common.Address `json:"from"`
		Amount       *hexutil.Big    `json:"value"`
		Payload      hexutil.Bytes   `json:"input"`
		V            *hexutil.Big    `json:"v"`
		R            *hexutil.Big    `json:"r"`
		S            *hexutil.Big    `json:"s"`
		Hash         *common.Hash    `json:"hash"`
		ChainID      *hexutil.Big    `json:"chainId"`
		Timestamp    *hexutil.Uint64 `json:"timestamp"`
	}
	var enc txdata
	enc.AccountNonce = hexutil.Uint64(t.Nonce())
	enc.Price = (*hexutil.Big)(t.GasPrice())
	enc.GasLimit = hexutil.Uint64(t.Gas())
	enc.Recipient = t.To()
	enc.From = t.From
	enc.Amount = (*hexutil.Big)(t.Value())
	enc.Payload = t.Data()

	V, R, S := t.RawSignatureValues()
	enc.V = (*hexutil.Big)(V)
	enc.R = (*hexutil.Big)(R)
	enc.S = (*hexutil.Big)(S)

	hash := t.Hash()
	enc.Hash = &hash
	enc.ChainID = (*hexutil.Big)(t.ChainId())
	enc.Timestamp = (*hexutil.Uint64)(t.Timestamp)

	return json.Marshal(&enc)
}
