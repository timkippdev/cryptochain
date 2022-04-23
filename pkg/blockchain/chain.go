package blockchain

type Blockchain struct {
	chain []*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		chain: []*Block{GenesisBlock()},
	}
}

func (bc *Blockchain) AddBlock(data string) *Block {
	b := MineBlock(bc.GetLastBlock(), data)
	bc.chain = append(bc.chain, b)
	return b
}

func (bc *Blockchain) GetChain() []*Block {
	return bc.chain
}

func (bc *Blockchain) GetLastBlock() *Block {
	if len(bc.chain) == 0 {
		return nil
	}
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) ReplaceChain(chain []*Block) {
	if len(chain) <= len(bc.chain) {
		return
	}

	if !ValidateChain(chain) {
		return
	}

	bc.chain = chain
}

func ValidateChain(chain []*Block) bool {
	if len(chain) == 0 || !IsGenesisBlock(chain[0]) {
		return false
	}

	for i := 1; i < len(chain); i++ {
		previousBlock := chain[i-1]
		previousBlockHash := previousBlock.Hash
		currentBlock := chain[i]
		if previousBlockHash != currentBlock.LastHash {
			return false
		}
		difficultyDifference := previousBlock.Difficulty - currentBlock.Difficulty
		if difficultyDifference > 1 || difficultyDifference < -1 {
			return false
		}

		if currentBlock.GenerateHash() != currentBlock.Hash {
			return false
		}
	}

	return true
}
