package models

type TxType int

const (
	TxTypeStandard TxType = 0
	TxTypeCoinbase TxType = 1
	TxTypeStake    TxType = 2
	TxTypeTreasury TxType = 3
	TxTypeMint     TxType = 4
)
