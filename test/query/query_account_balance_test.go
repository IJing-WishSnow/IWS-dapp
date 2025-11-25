package query

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestQueryAccountBalance(t *testing.T) {
	// 连接到以太坊节点
	// client, err := ethclient.Dial("https://eth-mainnet.g.alchemy.com/v2/4Mb8kv8N7tWzzTDYHAkE3")
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/4Mb8kv8N7tWzzTDYHAkE3")
	if err != nil {
		log.Fatal(err)
	}

	// 将十六进制地址转换为Address类型
	// account := common.HexToAddress("0x2281cB9267ABAF264c0A4c0dD2e414b4d68cE634")
	account := common.HexToAddress("0x8c8ab9b6178877246b224f8d745a1410c4928373")

	// 查询当前最新的余额
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance)

	// 查询特定区块高度的余额
	blockNumber := big.NewInt(5532993)
	balanceAt, err := client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balanceAt)

	// 将最新余额从wei转换为ETH
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println(ethValue)

	// 查询待处理余额（包含待处理交易）
	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println(pendingBalance)
}
