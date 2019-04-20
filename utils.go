package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func transactionFromETH(t *types.Transaction, timestamp *uint64) (*Transaction, error) {
	V, _, _ := t.RawSignatureValues()
	var signer types.Signer
	if V.Sign() != 0 && t.Protected() {
		signer = types.NewEIP155Signer(t.ChainId())
	} else {
		signer = types.HomesteadSigner{}
	}

	from, err := types.Sender(signer, t)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Transaction: t,
		From:        &from,
		Timestamp:   timestamp,
	}, nil
}

func transactionsFilterByFrom(t []*Transaction, address common.Address) ([]*Transaction, error) {
	newTx := []*Transaction{}
	for _, tx := range t {
		if tx.From != nil && *tx.From == address {
			newTx = append(newTx, tx)
		}
	}
	return newTx, nil
}
