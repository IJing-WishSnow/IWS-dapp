package query

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"testing"
)

func TestGetBlock(t *testing.T) {
	// 通过区块高度获取哈希
	blockNumber := big.NewInt(23866957)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
}
