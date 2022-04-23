package blockchain

import (
	"fmt"
	"strings"
	"time"

	"github.com/timkippdev/cryptochain/pkg/timestamp"
)

var (
	genesisData            string
	genesisDifficulty      int
	genesisNonce           int64
	genesisTimestamp       timestamp.Timestamp
	mineRateInMilliseconds int64
)

func init() {
	ts, _ := time.Parse(time.RFC3339, "2021-06-03T00:00:00Z")
	genesisData = "genesis"
	genesisDifficulty = 3
	genesisNonce = int64(0)
	genesisTimestamp = timestamp.FromTime(ts.UTC())
	mineRateInMilliseconds = 750
}

type Block struct {
	Data       interface{}         `json:"data"`
	Difficulty int                 `json:"difficulty"`
	Hash       string              `json:"hash"`
	LastHash   string              `json:"lastHash"`
	Nonce      int64               `json:"nonce"`
	Timestamp  timestamp.Timestamp `json:"timestamp"`
}

func (b *Block) GenerateHash() string {
	data := fmt.Sprintf("%v", b.Data)
	difficulty := fmt.Sprintf("%v", b.Difficulty)
	nonce := fmt.Sprintf("%v", b.Nonce)
	ts := b.Timestamp
	lastHash := b.LastHash

	return Hash(ts.String(), lastHash, data, nonce, difficulty)
}

func AdjustDifficulty(block *Block, timestamp timestamp.Timestamp) int {
	if block.Difficulty < 1 {
		return 1
	}

	var newDifficulty int
	if int64(timestamp-block.Timestamp) > mineRateInMilliseconds {
		newDifficulty = block.Difficulty - 1
	} else {
		newDifficulty = block.Difficulty + 1
	}

	if newDifficulty < 1 {
		newDifficulty = 1
	}

	return newDifficulty
}

func GenesisBlock() *Block {
	return newBlock(genesisTimestamp, genesisData, "", genesisDifficulty, genesisNonce)
}

func MineBlock(lastBlock *Block, data interface{}) *Block {
	ts := timestamp.Now()
	nonce := int64(0)

	b := newBlock(ts, data, lastBlock.Hash, lastBlock.Difficulty, nonce)

	for {
		// TODO: convert to doing Proof of Work comparing binary instead of current hex values
		if b.Hash[:b.Difficulty] == strings.Repeat("0", b.Difficulty) {
			break
		}

		nonce = nonce + 1
		ts := timestamp.Now()

		b.Timestamp = ts
		b.Difficulty = AdjustDifficulty(lastBlock, ts)
		b.Nonce = nonce
		b.Hash = b.GenerateHash()
	}

	return b
}

func IsGenesisBlock(block *Block) bool {
	genesisBlock := GenesisBlock()
	if block.Data != genesisBlock.Data || block.Hash != genesisBlock.Hash || block.Timestamp != genesisBlock.Timestamp {
		return false
	}
	return true
}

func newBlock(timestamp timestamp.Timestamp, data interface{}, lastHash string, difficulty int, nonce int64) *Block {
	b := &Block{
		Data:       data,
		Difficulty: difficulty,
		LastHash:   lastHash,
		Nonce:      nonce,
		Timestamp:  timestamp,
	}
	b.Hash = b.GenerateHash()
	return b
}
