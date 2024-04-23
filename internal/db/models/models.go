package models

type Block struct {
	ID         string `db:"id"`
	ParentID   string `db:"parent_id"`
	ProducerID string `db:"producer_id"`
	TXORoot    string `db:"txo_root"`
	Version    int    `db:"version"`
	Height     int    `db:"height"`
	Timestamp  int    `db:"timestamp"`
	Size       int    `db:"size"`
	TxCount    int    `db:"tx_count"`
}

type Transaction struct {
	ID       string `db:"id"`
	Type     int    `db:"type"`
	BlockID  string `db:"block_id"`
	TXORoot  string `db:"txo_root"`
	Locktime int    `db:"locktime"`
	Fee      uint64 `db:"fee"`
	Proof    string `db:"proof"`

	// Maybe
	Outputs           []Output `json:"outputs,omitempty"`
	Nullifiers        []string `json:"nullifiers,omitempty"`
	*Coinbase         `json:"coinbase,omitempty"`
	*Stake            `json:"stake,omitempty"`
	*TreasuryProposal `json:"treasury_proposal,omitempty"`
	*Mint             `json:"mint,omitempty"`
}

type Output struct {
	//TransactionID string `db:"transaction_id"`
	OutputIndex int    `db:"output_index"`
	Commitment  string `db:"commitment"`
	Ciphertext  string `db:"ciphertext"`
}

//type Nullifier struct {
//	NullifierID   string `db:"nullifier_id"`
//	TransactionID string `db:"transaction_id"`
//}

type Coinbase struct {
	TransactionID string `db:"transaction_id"`
	ValidatorID   string `db:"coinbase_validator_id"`
	Signature     string `db:"coinbase_signature"`
	NewCoins      uint64 `db:"coinbase_new_coins"`
}

type Stake struct {
	TransactionID string `db:"transaction_id"`
	ValidatorID   string `db:"stake_validator_id"`
	Amount        uint64 `db:"stake_amount"`
}

type TreasuryProposal struct {
	TransactionID string `db:"transaction_id"`
	ProposalHash  string `db:"proposal_hash"`
	Amount        uint64 `db:"proposal_amount"`
}

type Mint struct {
	TransactionID string `db:"transaction_id"`
	MintType      int    `db:"mint_type"`
	AssetID       string `db:"mint_asset_id"`
	DocumentHash  string `db:"mint_document_hash"`
	NewTokens     uint64 `db:"mint_new_tokens"`
	MintKey       string `db:"mint_key"`
}

type CoinbaseTransaction struct {
	Transaction
	Coinbase
}

type StakeTransaction struct {
	Transaction
	Stake
}

type TreasuryTransaction struct {
	Transaction
	TreasuryProposal
}

type MintTransaction struct {
	Transaction
	Mint
}
