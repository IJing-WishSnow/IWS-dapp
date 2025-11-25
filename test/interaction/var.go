package interaction

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	client      *ethclient.Client
	blockNumber = big.NewInt(23866957) // 全局区块号
	blockHash   = common.HexToHash("0x62a45449d23bc26e6a16970345ac132e5f88f8bc198d9757005670db8aa8d7d0")
	txHash      = common.HexToHash("0x25d95c09ff74fccfdb8eca54ad8d50e1d62eabe920402e735499b06769eca59a")
)

// 初始化时连接，失败则直接退出程序
func init() {
	var err error
	// 主网
	client, err = ethclient.Dial("https://eth-mainnet.g.alchemy.com/v2/4Mb8kv8N7tWzzTDYHAkE3")
	if err != nil {
		log.Fatalf("无法连接以太坊节点: %v", err) // 连接失败直接终止程序
	}
}
