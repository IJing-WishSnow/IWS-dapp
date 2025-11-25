package query

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

func TestQueryBlockReceipts(t *testing.T) {
	// 添加超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 通过区块哈希获取整个区块的所有交易收据
	receiptByHash, err := client.BlockReceipts(ctx, rpc.BlockNumberOrHashWithHash(blockHash, false))
	if err != nil {
		log.Fatal("通过哈希获取区块收据失败:", err)
	}

	// 通过区块号获取整个区块的所有交易收据（修复类型转换）
	receiptsByNum, err := client.BlockReceipts(ctx, rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumber.Int64())))
	if err != nil {
		log.Fatal("通过区块号获取区块收据失败:", err)
	}

	// 验证两种方式获取的收据是否相同
	if len(receiptByHash) > 0 && len(receiptsByNum) > 0 {
		fmt.Println("两种方式获取的收据是否相同:", receiptByHash[0].TxHash == receiptsByNum[0].TxHash)
	}

	// 遍历区块中的交易收据（限制数量避免超时）
	for i, receipt := range receiptByHash {
		if i >= 3 { // 只处理前3个收据
			break
		}
		fmt.Printf("\n=== 收据 %d ===\n", i)
		fmt.Println("交易状态:", receipt.Status)                // 1=成功, 0=失败
		fmt.Println("事件日志数量:", len(receipt.Logs))           // 交易产生的事件日志数组
		fmt.Println("交易哈希:", receipt.TxHash.Hex())          // 交易哈希标识符
		fmt.Println("交易索引:", receipt.TransactionIndex)      // 交易在区块中的位置索引
		fmt.Println("合约地址:", receipt.ContractAddress.Hex()) // 合约创建交易的地址
	}

	// 通过交易哈希单独查询交易收据
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		log.Fatal("通过交易哈希获取收据失败:", err)
	}

	// 显示单个交易的收据详情
	fmt.Println("\n=== 单独交易收据 ===")
	fmt.Println("交易状态:", receipt.Status)
	fmt.Println("事件日志数量:", len(receipt.Logs))
	fmt.Println("交易哈希:", receipt.TxHash.Hex())
	fmt.Println("交易索引:", receipt.TransactionIndex)
	fmt.Println("合约地址:", receipt.ContractAddress.Hex())
}
