package db

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/tyler-smith/iexplorer/internal/db/models"
)

type BlockSort int

const (
	BlockSortHeightAsc BlockSort = iota
	BlockSortHeightDesc
)

func GetBlocks(db *sqlx.DB, limit int, offset int) ([]models.Block, error) {
	ctx := context.Background()
	blocks := []models.Block{}
	q := "SELECT * FROM blocks ORDER BY height DESC LIMIT ? OFFSET ?"
	err := db.SelectContext(ctx, &blocks, q, limit, offset)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

func GetBlock(db *sqlx.DB, id string) (models.Block, error) {
	ctx := context.Background()
	block := models.Block{}
	err := db.GetContext(ctx, &block, sqlSelectBlock, id)
	return block, err
}

func GetTransaction(db *sqlx.DB, id string) (models.Transaction, error) {
	ctx := context.Background()

	tx := models.Transaction{}
	if err := db.GetContext(ctx, &tx, sqlSelectTransaction, id); err != nil {
		return tx, err
	}

	// Get outputs and nullifiers
	outputs := []models.Output{}
	if err := db.SelectContext(ctx, &outputs, sqlSelectTransactionOutputs, id); err != nil {
		return models.Transaction{}, err
	}
	tx.Outputs = outputs

	nullifiers := []string{}
	if err := db.SelectContext(ctx, &nullifiers, sqlSelectTransactionNullifiers, id); err != nil {
		return models.Transaction{}, err
	}
	tx.Nullifiers = nullifiers

	return tx, nil
}

const (
	sqlSelectBlock = `
SELECT * FROM blocks WHERE id = ? LIMIT 1;
`

	sqlSelectTransaction = `
SELECT t.id,
       t.block_id,
       t.txo_root,
       t.type,
       t.locktime,
       t.fee,
       t.proof,
       COALESCE(coinbases.validator_id, '') AS coinbase_validator_id,
       COALESCE(coinbases.signature, '') AS coinbase_signature,
       COALESCE(coinbases.new_coins, 0) AS coinbase_new_coins,
       COALESCE(stakes.validator_id, '')    AS stake_validator_id,
       COALESCE(stakes.amount, 0) AS stake_amount,
       COALESCE(treasury_proposals.proposal_hash, '') AS proposal_hash,
       COALESCE(treasury_proposals.amount, 0) AS proposal_amount,
       COALESCE(mints.mint_type, 0) AS mint_type,
       COALESCE(mints.asset_id, '') AS mint_asset_id,
       COALESCE(mints.document_hash, '') AS mint_document_hash,
       COALESCE(mints.new_tokens, 0) AS mint_new_tokens,
       COALESCE(mints.mint_key, '') AS mint_key
FROM transactions AS t
         LEFT JOIN coinbases ON t.id = coinbases.transaction_id AND t.type = 1
         LEFT JOIN stakes ON t.id = stakes.transaction_id AND t.type = 2
         LEFT JOIN treasury_proposals ON t.id = treasury_proposals.transaction_id AND t.type = 3
         LEFT JOIN mints ON t.id = mints.transaction_id AND t.type = 4
WHERE t.id = ?
LIMIT 1;
`

	sqlSelectTransactionOutputs = `
SELECT o.output_index, o.commitment, o.ciphertext
FROM outputs AS o
WHERE transaction_id = ?;
`

	sqlSelectTransactionNullifiers = `
SELECT n.id FROM nullifiers AS n WHERE transaction_id = ?;
`
)
