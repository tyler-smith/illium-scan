package viewmodels

import (
	"strconv"

	"github.com/tyler-smith/iexplorer/internal/db/models"
)

type Transaction struct {
	ID       string
	BlockID  string
	TXORoot  string
	Type     int
	Locktime int
	Fee      string
	Proof    string

	Outputs    []Output
	Nullifiers []string

	Coinbase         Coinbase
	Stake            Stake
	TreasuryProposal TreasuryProposal
	Mint             Mint
}

type Output struct {
	Index      string
	Commitment string
	Ciphertext string
}

type Coinbase struct {
	ValidatorID string
	Signature   string
	NewCoins    uint64
}

type Stake struct {
	ValidatorID string
	Amount      uint64
}

type TreasuryProposal struct {
	ProposalHash string
	Amount       uint64
}

type Mint struct {
	MintType     int
	AssetID      string
	DocumentHash string
	NewTokens    uint64
	MintKey      string
}

type TransactionsIndex struct {
	Transactions []Transaction
}

func NewTransactionsIndex(txs []models.Transaction) TransactionsIndex {
	vm := TransactionsIndex{}
	for _, tx := range txs {
		vm.Transactions = append(vm.Transactions, TransactionModelToViewModel(tx))
	}
	return vm
}

type TransactionsShow struct {
	Transaction Transaction
}

func NewTransactionsShow(tx models.Transaction) TransactionsShow {
	return TransactionsShow{
		Transaction: TransactionModelToViewModel(tx),
	}
}

func TransactionModelToViewModel(tx models.Transaction) Transaction {
	outputs := make([]Output, len(tx.Outputs))
	for i, output := range tx.Outputs {
		outputs[i] = Output{
			Index:      strconv.Itoa(output.OutputIndex),
			Commitment: output.Commitment,
			Ciphertext: output.Ciphertext,
		}
	}

	vm := Transaction{
		ID:       tx.ID,
		BlockID:  tx.BlockID,
		TXORoot:  tx.TXORoot,
		Type:     tx.Type,
		Locktime: tx.Locktime,
		Fee:      strconv.Itoa(int(tx.Fee)),
		Proof:    tx.Proof,

		Outputs:    outputs,
		Nullifiers: tx.Nullifiers,

		//Coinbase: Coinbase{
		//	ValidatorID: tx.Coinbase.ValidatorID,
		//	Signature:   tx.Coinbase.Signature,
		//	NewCoins:    tx.Coinbase.NewCoins,
		//},
		//
		//Stake: Stake{
		//	ValidatorID: tx.Stake.ValidatorID,
		//	Amount:      tx.Stake.Amount,
		//},
		//
		//TreasuryProposal: TreasuryProposal{
		//	ProposalHash: tx.TreasuryProposal.ProposalHash,
		//	Amount:       tx.TreasuryProposal.Amount,
		//},
		//
		//Mint: Mint{
		//	MintType:     tx.Mint.MintType,
		//	AssetID:      tx.Mint.AssetID,
		//	DocumentHash: tx.Mint.DocumentHash,
		//	NewTokens:    tx.Mint.NewTokens,
		//	MintKey:      tx.Mint.MintKey,
		//},
	}

	if tx.Coinbase != nil {
		vm.Coinbase = Coinbase{
			ValidatorID: tx.Coinbase.ValidatorID,
			Signature:   tx.Coinbase.Signature,
			NewCoins:    tx.Coinbase.NewCoins,
		}
	}

	if tx.Stake != nil {
		vm.Stake = Stake{
			ValidatorID: tx.Stake.ValidatorID,
			Amount:      tx.Stake.Amount,
		}
	}

	if tx.TreasuryProposal != nil {
		vm.TreasuryProposal = TreasuryProposal{
			ProposalHash: tx.TreasuryProposal.ProposalHash,
			Amount:       tx.TreasuryProposal.Amount,
		}
	}

	if tx.Mint != nil {
		vm.Mint = Mint{
			MintType:     tx.Mint.MintType,
			AssetID:      tx.Mint.AssetID,
			DocumentHash: tx.Mint.DocumentHash,
			NewTokens:    tx.Mint.NewTokens,
			MintKey:      tx.Mint.MintKey,
		}
	}

	return vm
}
