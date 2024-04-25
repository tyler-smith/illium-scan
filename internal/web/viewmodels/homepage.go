package viewmodels

import "github.com/tyler-smith/iexplorer/internal/db/models"

type Homepage struct {
	Blocks []Block
	Stakes []Stake
}

func NewHomepage(blocks []models.Block, stakes []models.Stake) Homepage {
	modelBlocks := make([]Block, len(blocks))
	for i, block := range blocks {
		modelBlocks[i] = blockModelToViewModel(block)
	}

	modelStakes := make([]Stake, len(stakes))
	for i, stake := range stakes {
		modelStakes[i] = stakeModelToViewModel(stake)
	}

	return Homepage{
		Blocks: modelBlocks,
		Stakes: modelStakes,
	}
}

func stakeModelToViewModel(stake models.Stake) Stake {
	return Stake{
		ValidatorID: stake.ValidatorID,
		Amount:      stake.Amount,
	}
}
