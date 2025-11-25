package interaction

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TestSubBlock 测试通过WebSocket订阅新区块头功能
// 演示了如何实时监听以太坊网络的新区块生成
func TestSubBlock(t *testing.T) {
	// 通过WebSocket连接到以太坊Sepolia测试网络
	// WebSocket连接支持实时订阅功能，适合监听区块和交易事件
	client, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/4Mb8kv8N7tWzzTDYHAkE3")
	if err != nil {
		log.Fatal("连接以太坊节点失败:", err)
	}

	// 创建用于接收新区块头的通道
	headers := make(chan *types.Header)

	// 订阅新区块头事件
	// 当网络中有新区块产生时，区块头信息会通过这个通道发送
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal("订阅新区块头失败:", err)
	}

	// 无限循环监听新区块事件
	for {
		select {
		case err := <-sub.Err():
			// 处理订阅错误
			log.Fatal("订阅错误:", err)
		case header := <-headers:
			// 接收到新的区块头信息

			// 打印区块哈希
			fmt.Println("新区块哈希:", header.Hash().Hex())

			// 通过区块哈希获取完整的区块信息
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal("获取区块详情失败:", err)
			}

			// 打印区块详细信息
			fmt.Println("区块哈希:", block.Hash().Hex())        // 区块哈希值
			fmt.Println("区块高度:", block.Number().Uint64())   // 区块号/高度
			fmt.Println("区块时间戳:", block.Time())             // 区块时间戳（修复：直接使用uint64值）
			fmt.Println("区块随机数:", block.Nonce())            // 工作量证明随机数
			fmt.Println("交易数量:", len(block.Transactions())) // 区块中包含的交易数量

			// 可选：打印区块中的交易哈希
			for i, tx := range block.Transactions() {
				fmt.Printf("交易 %d: %s\n", i, tx.Hash().Hex())
			}

			fmt.Println("--- 新区块信息结束 ---")
		}
	}
}
