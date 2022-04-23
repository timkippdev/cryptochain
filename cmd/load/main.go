package main

import (
	"fmt"

	"github.com/timkippdev/cryptochain/pkg/blockchain"
	"github.com/timkippdev/cryptochain/pkg/timestamp"
)

func main() {
	bc := blockchain.NewBlockchain()
	bc.AddBlock("block-initial")

	var previousTimestamp, nextTimestamp timestamp.Timestamp
	var nextBlock *blockchain.Block
	var timeDifference int64
	var averageTime int64
	times := make([]int64, 0)

	for i := 0; i < 10000; i++ {
		previousTimestamp = bc.GetLastBlock().Timestamp

		bc.AddBlock(fmt.Sprintf("block-%d", i))
		nextBlock = bc.GetLastBlock()

		nextTimestamp = nextBlock.Timestamp
		timeDifference = int64(nextTimestamp) - int64(previousTimestamp)
		times = append(times, timeDifference)

		totalTime := int64(0)
		for _, t := range times {
			totalTime += t
		}
		averageTime = totalTime / int64(len(times))

		fmt.Printf("Time to mine block: %vms | Difficulty: %d | Average time: %vms | Hash: %s\n", timeDifference, nextBlock.Difficulty, averageTime, nextBlock.Hash)
	}
}
