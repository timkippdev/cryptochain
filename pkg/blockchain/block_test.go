package blockchain

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timkippdev/cryptochain/pkg/timestamp"
)

func Test_AdjustDifficulty_RaisesWhenTooFast(t *testing.T) {
	b := &Block{
		Difficulty: 3,
		Timestamp:  timestamp.Now(),
	}

	ts := b.Timestamp.Add(mineRateInMilliseconds - 100)
	assert.Equal(t, b.Difficulty+1, AdjustDifficulty(b, ts))
}

func Test_AdjustDifficulty_LowersWhenTooSlow(t *testing.T) {
	b := &Block{
		Difficulty: 3,
		Timestamp:  timestamp.Now(),
	}
	ts := b.Timestamp.Add(mineRateInMilliseconds + 100)
	assert.Equal(t, b.Difficulty-1, AdjustDifficulty(b, ts))
}

func Test_AdjustDifficulty_NeverGoesBelowOne(t *testing.T) {
	b := &Block{
		Difficulty: 0,
		Timestamp:  timestamp.Now(),
	}

	ts := b.Timestamp.Add(mineRateInMilliseconds + 100)
	assert.Equal(t, 1, AdjustDifficulty(b, ts))
}

func Test_GenerateHash(t *testing.T) {
	b := Block{
		Data:       "test",
		Difficulty: 2,
		LastHash:   "lh",
		Nonce:      33,
	}
	expectedHash := Hash(b.Timestamp.String(), fmt.Sprintf("%v", b.Nonce), fmt.Sprintf("%v", b.Difficulty), b.LastHash, fmt.Sprintf("%v", b.Data))
	assert.Equal(t, expectedHash, b.GenerateHash())
}

func Test_GenesisBlock(t *testing.T) {
	b := GenesisBlock()

	assert.Equal(t, genesisData, b.Data)
	assert.Equal(t, genesisDifficulty, b.Difficulty)
	assert.NotNil(t, b.Hash)
	assert.Equal(t, "", b.LastHash)
	assert.Equal(t, genesisNonce, b.Nonce)
	assert.Equal(t, genesisTimestamp, b.Timestamp)
}

func Test_IsGenesisBlock_InvalidData(t *testing.T) {
	gb := GenesisBlock()
	gb.Data = "not-genesis-data"
	assert.False(t, IsGenesisBlock(gb))
}

func Test_IsGenesisBlock_InvalidHash(t *testing.T) {
	gb := GenesisBlock()
	gb.Hash = "not-genesis-hash"
	assert.False(t, IsGenesisBlock(gb))
}

func Test_IsGenesisBlock_InvalidTimestamp(t *testing.T) {
	gb := GenesisBlock()
	gb.Timestamp = timestamp.Now().Add(-10)
	assert.False(t, IsGenesisBlock(gb))
}

func Test_IsGenesisBlock_Sanity(t *testing.T) {
	gb := GenesisBlock()
	assert.True(t, IsGenesisBlock(gb))
}

func Test_MineBlock(t *testing.T) {
	lastBlock := GenesisBlock()
	d := "test-data"
	minedBlock := MineBlock(lastBlock, d)

	assert.Equal(t, d, minedBlock.Data)
	assert.Contains(t, []int{lastBlock.Difficulty + 1, lastBlock.Difficulty - 1}, minedBlock.Difficulty)
	assert.NotEmpty(t, minedBlock.Hash)
	assert.Equal(t, strings.Repeat("0", minedBlock.Difficulty), minedBlock.Hash[:minedBlock.Difficulty])
	assert.Equal(t, lastBlock.Hash, minedBlock.LastHash)
	assert.GreaterOrEqual(t, minedBlock.Nonce, int64(0))
	assert.Greater(t, int64(minedBlock.Timestamp), int64(0))
}
