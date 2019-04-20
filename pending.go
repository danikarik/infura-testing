package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// PendingTransactions polling new transactions from node, when first time called creates filter and gets transactions from pending block.
func PendingTransactions(ctx context.Context, rpc *rpc.Client, client *ethclient.Client, opts *PendingTransactionsQuery) (*PendingTransactionsResponse, error) {
	defer newTimer().Start().Finish("PendingTransactions")
	if opts == nil || opts.FilterID == "" {
		return initPendingTransactions(ctx, rpc, opts)
	}

	var hashes []string
	err := rpc.CallContext(ctx, &hashes, "eth_getFilterChanges", opts.FilterID)
	if err != nil {
		if err.Error() == "Filter not found" {
			return initPendingTransactions(ctx, rpc, opts)
		}
		return nil, err
	}

	header, err := pendingBlockHeader(ctx, rpc)
	if err != nil {
		return nil, err
	}

	transactions := []*Transaction{}
	for _, h := range hashes {
		ethTx, isPending, err := client.TransactionByHash(ctx, common.HexToHash(h))
		if err != nil {
			return nil, fmt.Errorf("TransactionByHash: %v", err)
		}

		if !isPending {
			continue
		}

		tx, err := transactionFromETH(ethTx, (*uint64)(header.Time))
		if err != nil {
			return nil, fmt.Errorf("transactionFromETH: %v", err)
		}

		transactions = append(transactions, tx)
	}

	if opts != nil && opts.From != nil {
		transactions, err = transactionsFilterByFrom(transactions, *opts.From)
		if err != nil {
			return nil, fmt.Errorf("transactionsFilterByFrom: %v", err)
		}
	}

	return &PendingTransactionsResponse{
		FilterID:     opts.FilterID,
		Transactions: transactions,
	}, nil
}

func initPendingTransactions(ctx context.Context, rpc *rpc.Client, opts *PendingTransactionsQuery) (*PendingTransactionsResponse, error) {
	defer newTimer().Start().Finish("initPendingTransactions")
	tx, err := pendingTransactions(ctx, rpc)
	if err != nil {
		return nil, fmt.Errorf("pendingTransactions: %v", err)
	}

	filterID, err := newPendingTransactionFilter(ctx, rpc)
	if err != nil {
		return nil, fmt.Errorf("newPendingTransactionFilter: %v", err)
	}

	if opts != nil && opts.From != nil {
		tx, err = transactionsFilterByFrom(tx, *opts.From)
		if err != nil {
			return nil, fmt.Errorf("transactionsFilterByFrom: %v", err)
		}
	}

	return &PendingTransactionsResponse{
		Transactions: tx,
		FilterID:     filterID,
	}, nil
}

func pendingTransactions(ctx context.Context, rpc *rpc.Client) ([]*Transaction, error) {
	defer newTimer().Start().Finish("pendingTransactions")
	resp, err := pendingBlockTransactions(ctx, rpc)
	if err != nil {
		return nil, err
	}

	txs := make([]*Transaction, len(resp.Transactions))
	for i, t := range resp.Transactions {
		tx, err := transactionFromETH(t, (*uint64)(resp.Time))
		if err != nil {
			return nil, err
		}
		txs[i] = tx
	}

	return txs, nil
}

func newPendingTransactionFilter(ctx context.Context, rpc *rpc.Client) (string, error) {
	defer newTimer().Start().Finish("newPendingTransactionFilter")
	var result string
	err := rpc.CallContext(ctx, &result, "eth_newPendingTransactionFilter")
	if err != nil {
		return "", err
	}
	return result, nil
}

func pendingBlockHeader(ctx context.Context, rpc *rpc.Client) (*PendingBlockHeader, error) {
	defer newTimer().Start().Finish("pendingBlockHeader")
	var result PendingBlockHeader
	if err := rpc.CallContext(ctx, &result, "eth_getBlockByNumber", "pending", false); err != nil {
		return nil, err
	}
	return &result, nil
}

func pendingBlockTransactions(ctx context.Context, rpc *rpc.Client) (*PendingBlockTransactions, error) {
	defer newTimer().Start().Finish("pendingBlockTransactions")
	var result PendingBlockTransactions
	if err := rpc.CallContext(ctx, &result, "eth_getBlockByNumber", "pending", true); err != nil {
		return nil, err
	}
	return &result, nil
}

type timer struct {
	start time.Time
}

func newTimer() *timer {
	return &timer{}
}

func (t *timer) Start() *timer {
	t.start = time.Now()
	return t
}

func (t *timer) Finish(tag string) {
	log.Printf("%s: [%v]\n", tag, time.Since(t.start))
}
