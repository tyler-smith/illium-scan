package db

import (
	"context"
	"database/sql"
	"encoding/hex"

	"github.com/project-illium/ilxd/rpc/pb"
	"github.com/project-illium/ilxd/types/transactions"
)

const (
	sqlInsertBlock = `
		INSERT INTO blocks (id, parent_id, producer_id, txo_root, version, height, timestamp, size, tx_count)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	sqlInsertTransaction = `
		INSERT INTO transactions (id, block_id, txo_root, type, locktime, fee, proof)
		VALUES (?, ?, ?, ?, ?, ? , ?);`

	sqlInsertOutput = `
		INSERT INTO outputs (transaction_id, output_index, commitment, ciphertext)
		VALUES (?, ?, ?, ?);`

	sqlInsertNullifier = `
		INSERT INTO nullifiers (transaction_id, nullifier_id)
		VALUES (?, ?);`

	sqlInsertCoinbaseTransaction = `
		INSERT INTO coinbases (transaction_id, validator_id, signature, new_coins)
		VALUES (?, ?, ?, ?);`

	sqlInsertStakeTransaction = `
		INSERT INTO stakes (transaction_id, validator_id, amount)
		VALUES (?, ?, ?);`

	sqlInsertTreasuryTransaction = `
		INSERT INTO treasury_proposals (transaction_id, amount, proposal_hash)
		VALUES (?, ?, ?);`

	sqlInsertMintTransaction = `
		INSERT INTO mints (transaction_id, mint_type, asset_id, document_hash, new_tokens, mint_key)
		VALUES (?, ?, ?, ?, ?, ?);`
)

type txType int

const (
	txTypeStandard txType = 0
	txTypeCoinbase txType = 1
	txTypeStake    txType = 2
	txTypeTreasury txType = 3
	txTypeMint     txType = 4
)

type CachedWriterStmts struct {
	insertBlock       *sql.Stmt
	insertTransaction *sql.Stmt
	insertOutput      *sql.Stmt
	insertNullifier   *sql.Stmt
	insertCoinbaseTx  *sql.Stmt
	insertStakeTx     *sql.Stmt
	insertTreasuryTx  *sql.Stmt
	insertMintTx      *sql.Stmt
}

func NewCachedWriterStmts(db *sql.DB) (CachedWriterStmts, error) {
	ctx := context.Background()
	stmts := CachedWriterStmts{}

	blockStmt, err := db.PrepareContext(ctx, sqlInsertBlock)
	if err != nil {
		return stmts, err
	}
	transactionStmt, err := db.PrepareContext(ctx, sqlInsertTransaction)
	if err != nil {
		return stmts, err
	}
	outputStmt, err := db.PrepareContext(ctx, sqlInsertOutput)
	if err != nil {
		return stmts, err
	}
	nullifierStmt, err := db.PrepareContext(ctx, sqlInsertNullifier)
	if err != nil {
		return stmts, err
	}
	coinbaseStmt, err := db.PrepareContext(ctx, sqlInsertCoinbaseTransaction)
	if err != nil {
		return stmts, err
	}
	stakeStmt, err := db.PrepareContext(ctx, sqlInsertStakeTransaction)
	if err != nil {
		return stmts, err
	}
	treasuryStmt, err := db.PrepareContext(ctx, sqlInsertTreasuryTransaction)
	if err != nil {
		return stmts, err
	}
	mintStmt, err := db.PrepareContext(ctx, sqlInsertMintTransaction)
	if err != nil {
		return stmts, err
	}

	return CachedWriterStmts{
		insertBlock:       blockStmt,
		insertTransaction: transactionStmt,
		insertOutput:      outputStmt,
		insertNullifier:   nullifierStmt,
		insertCoinbaseTx:  coinbaseStmt,
		insertStakeTx:     stakeStmt,
		insertTreasuryTx:  treasuryStmt,
		insertMintTx:      mintStmt,
	}, nil
}

func (c *CachedWriterStmts) Close() error {
	var err error
	if c.insertBlock != nil {
		err = c.insertBlock.Close()
	}
	if c.insertTransaction != nil {
		err = c.insertTransaction.Close()
	}
	if c.insertOutput != nil {
		err = c.insertOutput.Close()
	}
	if c.insertNullifier != nil {
		err = c.insertNullifier.Close()
	}
	if c.insertCoinbaseTx != nil {
		err = c.insertCoinbaseTx.Close()
	}
	if c.insertStakeTx != nil {
		err = c.insertStakeTx.Close()
	}
	if c.insertTreasuryTx != nil {
		err = c.insertTreasuryTx.Close()
	}
	if c.insertMintTx != nil {
		err = c.insertMintTx.Close()
	}
	return err
}

func (c *CachedWriterStmts) ForTx(tx *sql.Tx) CachedWriterStmts {
	return CachedWriterStmts{
		insertBlock:       tx.Stmt(c.insertBlock),
		insertTransaction: tx.Stmt(c.insertTransaction),
		insertOutput:      tx.Stmt(c.insertOutput),
		insertNullifier:   tx.Stmt(c.insertNullifier),
		insertCoinbaseTx:  tx.Stmt(c.insertCoinbaseTx),
		insertStakeTx:     tx.Stmt(c.insertStakeTx),
		insertTreasuryTx:  tx.Stmt(c.insertTreasuryTx),
		insertMintTx:      tx.Stmt(c.insertMintTx),
	}
}

func bytesToHex(b []byte) string {
	return hex.EncodeToString(b)
}

func InsertBlock(stmts CachedWriterStmts, block *pb.BlockInfo) error {
	_, err := stmts.insertBlock.Exec(bytesToHex(block.GetBlock_ID()),
		bytesToHex(block.GetParent()),
		bytesToHex(block.GetProducer_ID()),
		bytesToHex(block.GetTxRoot()),
		block.GetVersion(),
		block.GetHeight(),
		block.GetTimestamp(),
		block.GetSize(),
		block.GetNumTxs(),
	)
	return err
}

func InsertTransaction(stmts CachedWriterStmts, block_id []byte, tx *transactions.Transaction) error {
	id := tx.ID()

	var (
		txtype txType = txTypeStandard

		// Common transaction properties
		outputs    []*transactions.Output
		nullifiers [][]byte
		txoRoot    []byte
		locktime   int64
		fee        uint64
		proof      []byte

		// Coinbase only
		validatorID []byte
		signature   []byte
		newCoins    uint64

		// Stake only
		amount uint64

		// Treasury only
		proposalHash []byte

		// Mint only
		mintType     int32
		assetID      []byte
		documentHash []byte
		newTokens    uint64
		mintKey      []byte
	)

	switch concreteTx := tx.Tx.(type) {
	case *transactions.Transaction_StandardTransaction:
		outputs = concreteTx.StandardTransaction.GetOutputs()
		nullifiers = concreteTx.StandardTransaction.GetNullifiers()
		txoRoot = concreteTx.StandardTransaction.GetTxoRoot()
		locktime = concreteTx.StandardTransaction.GetLocktime()
		fee = concreteTx.StandardTransaction.GetFee()
		proof = concreteTx.StandardTransaction.GetProof()
	case *transactions.Transaction_CoinbaseTransaction:
		txtype = txTypeCoinbase

		outputs = concreteTx.CoinbaseTransaction.GetOutputs()
		proof = concreteTx.CoinbaseTransaction.GetProof()
		validatorID = concreteTx.CoinbaseTransaction.GetValidator_ID()
		signature = concreteTx.CoinbaseTransaction.GetSignature()
		newCoins = concreteTx.CoinbaseTransaction.GetNewCoins()
	case *transactions.Transaction_StakeTransaction:
		txtype = txTypeStake

		validatorID = concreteTx.StakeTransaction.GetValidator_ID()
		amount = concreteTx.StakeTransaction.GetAmount()
		nullifiers = [][]byte{concreteTx.StakeTransaction.GetNullifier()}
		txoRoot = concreteTx.StakeTransaction.GetTxoRoot()
		locktime = concreteTx.StakeTransaction.GetLocktime()
		signature = concreteTx.StakeTransaction.GetSignature()
		proof = concreteTx.StakeTransaction.GetProof()
	case *transactions.Transaction_TreasuryTransaction:
		txtype = txTypeTreasury

		amount = concreteTx.TreasuryTransaction.GetAmount()
		outputs = concreteTx.TreasuryTransaction.GetOutputs()
		proposalHash = concreteTx.TreasuryTransaction.GetProposalHash()
		proof = concreteTx.TreasuryTransaction.GetProof()
	case *transactions.Transaction_MintTransaction:
		txtype = txTypeMint

		mintType = int32(concreteTx.MintTransaction.GetType())
		assetID = concreteTx.MintTransaction.GetAsset_ID()
		documentHash = concreteTx.MintTransaction.GetDocumentHash()
		newTokens = concreteTx.MintTransaction.GetNewTokens()
		mintKey = concreteTx.MintTransaction.GetMintKey()

		outputs = concreteTx.MintTransaction.GetOutputs()
		fee = concreteTx.MintTransaction.GetFee()
		nullifiers = concreteTx.MintTransaction.GetNullifiers()
		txoRoot = concreteTx.MintTransaction.GetTxoRoot()
		locktime = concreteTx.MintTransaction.GetLocktime()
		signature = concreteTx.MintTransaction.GetSignature()
		proof = concreteTx.MintTransaction.GetProof()
	}

	// Insert base transaction.
	_, err := stmts.insertTransaction.Exec(
		bytesToHex(id[:]),
		bytesToHex(block_id),
		bytesToHex(txoRoot),
		txtype,
		locktime,
		fee,
		proof,
	)
	if err != nil {
		return err
	}

	// Insert each output.
	for i, output := range outputs {
		_, err := stmts.insertOutput.Exec(
			bytesToHex(id[:]),
			i,
			bytesToHex(output.GetCommitment()),
			bytesToHex(output.GetCiphertext()),
		)
		if err != nil {
			return err
		}
	}

	// Insert each nullifier.
	for _, nullifier := range nullifiers {
		_, err := stmts.insertNullifier.Exec(
			bytesToHex(id[:]),
			bytesToHex(nullifier),
		)
		if err != nil {
			return err
		}
	}

	// Insert sub-type specific data.
	switch txtype {
	case txTypeCoinbase:
		_, err = stmts.insertCoinbaseTx.Exec(
			bytesToHex(id[:]),
			bytesToHex(validatorID),
			bytesToHex(signature),
			newCoins,
		)
	case txTypeStake:
		_, err = stmts.insertStakeTx.Exec(
			bytesToHex(id[:]),
			bytesToHex(validatorID),
			amount,
		)
	case txTypeTreasury:
		_, err = stmts.insertTreasuryTx.Exec(
			bytesToHex(id[:]),
			amount,
			bytesToHex(proposalHash),
		)
	case txTypeMint:
		_, err = stmts.insertMintTx.Exec(
			bytesToHex(id[:]),
			mintType,
			bytesToHex(assetID),
			bytesToHex(documentHash),
			newTokens,
			bytesToHex(mintKey),
		)
	}
	if err != nil {
		return err
	}

	return nil
}
