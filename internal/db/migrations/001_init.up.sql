CREATE TABLE IF NOT EXISTS blocks (
    id          CHAR(64) PRIMARY KEY,
    parent_id   CHAR(64) NOT NULL,
    producer_id VARCHAR NOT NULL,
    txo_root    CHAR(64) NOT NULL,

    version     INTEGER  NOT NULL,
    height      INTEGER  NOT NULL,
    timestamp   INTEGER  NOT NULL,
    size        INTEGER  NOT NULL,
    tx_count    INTEGER  NOT NULL
    );

CREATE TABLE IF NOT EXISTS transactions (
    id       CHAR(64) PRIMARY KEY,
    block_id CHAR(64),
    txo_root CHAR(64),
    type    SMALLINT NOT NULL,
    locktime INTEGER,
    fee      BIGINT NOT NULL DEFAULT (0),
    proof    TEXT   NOT NULL
    );

CREATE TABLE IF NOT EXISTS outputs (
    transaction_id CHAR(64)        NOT NULL,
    output_index   INTEGER         NOT NULL,
    commitment     CHAR(64)        NOT NULL,
    ciphertext     TEXT            NOT NULL,
    PRIMARY KEY (transaction_id, output_index),
    FOREIGN KEY (transaction_id) REFERENCES transactions (id)
    );

CREATE TABLE IF NOT EXISTS nullifiers (
    id   CHAR(64) PRIMARY KEY,
    transaction_id CHAR(64),
    FOREIGN KEY (transaction_id) REFERENCES transactions (id)
    );

CREATE TABLE IF NOT EXISTS coinbases (
    transaction_id CHAR(64) PRIMARY KEY,
    validator_id   CHAR(64) NOT NULL,
    signature      CHAR(64) NOT NULL,
    new_coins      BIGINT   NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions (id)
    );

CREATE TABLE IF NOT EXISTS stakes (
    transaction_id CHAR(64) PRIMARY KEY,
    validator_id   CHAR(64) NOT NULL,
    amount         BIGINT   NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions (id)
    );

CREATE TABLE IF NOT EXISTS treasury_proposals (
    transaction_id CHAR(64) PRIMARY KEY,
    proposal_hash  TEXT   NOT NULL,
    amount         BIGINT NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions (id)
    );

CREATE TABLE IF NOT EXISTS mints (
    transaction_id CHAR(64) PRIMARY KEY,
    mint_type      SMALLINT NOT NULL,
    asset_id       CHAR(64) NOT NULL,
    document_hash  CHAR(64) NOT NULL,
    new_tokens     BIGINT   NOT NULL,
    mint_key       CHAR(64) NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions (id)
    );

CREATE INDEX IF NOT EXISTS idx_transactions_block_id ON transactions (block_id);
CREATE INDEX IF NOT EXISTS idx_outputs_transaction_id ON outputs (transaction_id);
CREATE INDEX IF NOT EXISTS idx_nullifiers_transaction_id ON nullifiers (transaction_id);
CREATE INDEX IF NOT EXISTS idx_coinbase_transaction_id ON coinbases (transaction_id);
CREATE INDEX IF NOT EXISTS idx_stake_transaction_id ON stakes (transaction_id);
CREATE INDEX IF NOT EXISTS idx_treasury_transaction_id ON treasury_proposals (transaction_id);
CREATE INDEX IF NOT EXISTS idx_mint_transaction_id ON mints (transaction_id);
