package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timkippdev/cryptochain/pkg/timestamp"
)

func Test_NewBlockchain(t *testing.T) {
	g := GenesisBlock()
	bc := NewBlockchain()
	if assert.Len(t, bc.chain, 1) {
		assert.Equal(t, bc.chain[0].Hash, g.Hash)
	}
}

func Test_AddBlock(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("foo")
	assert.Len(t, bc.chain, 2)
	assert.Equal(t, "foo", bc.GetLastBlock().Data)

	bc.AddBlock("bar")
	assert.Len(t, bc.chain, 3)
	assert.Equal(t, "bar", bc.GetLastBlock().Data)
}

func Test_ReplaceChain_NewChainIsLonger_InvalidChain(t *testing.T) {
	oldChain := NewBlockchain()
	oldChainRef := oldChain.chain

	newChain := NewBlockchain()
	newChain.AddBlock("foo")

	// simulate invalid block
	newChain.GetLastBlock().Data = "not-foo"

	oldChain.ReplaceChain(newChain.chain)
	assert.Equal(t, oldChainRef, oldChain.chain)
}

func Test_ReplaceChain_NewChainIsLonger_ValidChain(t *testing.T) {
	oldChain := NewBlockchain()
	newChain := NewBlockchain()
	newChain.AddBlock("foo")

	oldChain.ReplaceChain(newChain.chain)
	assert.Equal(t, newChain.chain, oldChain.chain)
}

func Test_ReplaceChain_NewChainIsNotLonger(t *testing.T) {
	oldChain := NewBlockchain()
	oldChain.AddBlock("foo")
	oldChain.AddBlock("bar")
	oldChainRef := oldChain.chain

	newChain := NewBlockchain()

	oldChain.ReplaceChain(newChain.chain)
	assert.Equal(t, oldChainRef, oldChain.chain)
}

func TestValidateChain_VerifyGenesisBlock(t *testing.T) {
	bc := NewBlockchain()
	bc.chain[0].Data = "not-genesis-data"
	assert.False(t, ValidateChain(bc.chain))
}

func TestValidateChain_InvalidLastHashReference(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("foo")
	bc.AddBlock("bar")
	bc.chain[2].LastHash = "not-valid-hash"

	assert.False(t, ValidateChain(bc.chain))
}

func TestValidateChain_BlockContainsInvalidField(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("foo")
	bc.AddBlock("bar")
	bc.chain[2].Data = "not-valid-data"

	assert.False(t, ValidateChain(bc.chain))
}

func TestValidateChain_BlockContainsDifficultyJump_PositiveDirection(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("foo")

	lb := bc.GetLastBlock()
	corruptBlock := newBlock(timestamp.Now(), "corrupt", lb.Hash, lb.Difficulty+2, 0)
	bc.chain = append(bc.chain, corruptBlock)
	assert.False(t, ValidateChain(bc.chain))
}

func TestValidateChain_BlockContainsDifficultyJump_NegativeDirection(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("foo")

	lb := bc.GetLastBlock()
	corruptBlock := newBlock(timestamp.Now(), "corrupt", lb.Hash, lb.Difficulty-2, 0)
	bc.chain = append(bc.chain, corruptBlock)
	assert.False(t, ValidateChain(bc.chain))
}

func TestValidateChain_AllBlocksAreValid(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("foo")
	bc.AddBlock("bar")

	assert.True(t, ValidateChain(bc.chain))
}
