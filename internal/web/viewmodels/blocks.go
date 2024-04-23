package viewmodels

import (
	"strconv"
	"time"

	"github.com/tyler-smith/iexplorer/internal/db/models"
)

type Block struct {
	ID         string
	Height     string
	ProducerID string
	Time       string
	TxCount    string
	Size       string
	URL        string
}

type BlocksIndex struct {
	Blocks []Block
}

func NewBlocksIndex(blocks []models.Block) BlocksIndex {
	vm := BlocksIndex{}
	for _, block := range blocks {
		vm.Blocks = append(vm.Blocks, blockModelToViewModel(block))
	}
	return vm
}

type BlocksShow struct {
	Block Block
}

func NewBlocksShow(block models.Block) BlocksShow {
	return BlocksShow{
		Block: blockModelToViewModel(block),
	}
}

func blockModelToViewModel(block models.Block) Block {
	producer := block.ProducerID
	if block.Height == 0 {
		producer = "(genesis)"
	}

	return Block{
		ID:         block.ID,
		Height:     strconv.Itoa(block.Height),
		ProducerID: producer,
		Time:       time.Unix(int64(block.Timestamp), 0).Format(time.DateTime),
		TxCount:    strconv.Itoa(block.TxCount),
		Size:       strconv.Itoa(block.Size),
		URL:        "/blocks/" + block.ID,
	}
}
