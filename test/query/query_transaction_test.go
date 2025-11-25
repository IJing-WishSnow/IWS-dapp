package query

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

func TestQueryTransaction(t *testing.T) {
	// 创建带超时的上下文，避免请求卡死
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 获取区块链网络ID - 用于交易签名验证
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal("获取链ID失败:", err)
	}
	fmt.Printf("当前网络链ID: %s\n", chainID.String())

	// 使用一个已知的区块号进行测试
	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		log.Fatal("获取区块信息失败:", err)
	}
	fmt.Printf("成功获取区块 #%d，包含 %d 笔交易\n", block.Number().Uint64(), len(block.Transactions()))

	// 遍历区块中的交易并分析交易数据
	txCount := 0
	for _, tx := range block.Transactions() {
		txCount++
		if txCount > 3 { // 限制处理3笔交易，避免超时
			fmt.Println("已处理3笔交易，跳过剩余交易...")
			break
		}

		fmt.Println("\n=== 交易分析 ===")
		fmt.Println(tx.Hash().Hex())        // 交易哈希 - 交易的唯一标识符，用于在区块链上唯一识别该交易
		fmt.Println(tx.Value().String())    // 交易金额（wei）- 转账的以太币数量，1 ETH = 10^18 wei
		fmt.Println(tx.Gas())               // Gas限制 - 交易允许消耗的最大Gas量，防止无限循环和过度消耗资源
		fmt.Println(tx.GasPrice().Uint64()) // Gas价格（wei）- 每单位Gas的价格，决定交易处理优先级
		fmt.Println(tx.Nonce())             // 发送者交易计数器 - 防止重放攻击，确保交易顺序执行
		fmt.Println(tx.Data())              // 交易附加数据 - 智能合约调用参数或备注信息，普通转账为空
		if tx.To() != nil {
			fmt.Println(tx.To().Hex()) // 接收方地址 - 资金或合约调用的目标地址，nil表示合约创建交易
		} else {
			fmt.Println("合约创建交易")
		}

		// 从交易签名中恢复发送者地址 - 使用最新签名器支持所有交易类型
		if sender, err := types.Sender(types.LatestSignerForChainID(chainID), tx); err == nil {
			fmt.Println("发送者地址:", sender.Hex()) // 交易发起方的以太坊地址
		} else {
			fmt.Printf("恢复发送者失败: %v\n", err)
			continue // 跳过此交易继续处理下一个
		}

		// 获取交易收据 - 包含交易在区块链上的执行结果和产生的事件日志
		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			fmt.Printf("获取交易收据失败: %v\n", err)
			continue
		}

		// 交易状态: 1=成功, 0=失败
		fmt.Printf("交易状态: %d ", receipt.Status)
		if receipt.Status == 1 {
			fmt.Println("(成功)")
		} else {
			fmt.Println("(失败)")
		}

		// 交易产生的事件日志 - 智能合约执行过程中触发的事件记录
		fmt.Printf("事件日志数量: %d\n", len(receipt.Logs))
		for i, log := range receipt.Logs {
			fmt.Printf("  日志%d: 合约=%s, 主题数=%d, 数据长度=%d字节\n",
				i, log.Address.Hex(), len(log.Topics), len(log.Data))
		}
	}

	fmt.Println("\n=== 区块交易查询 ===")

	// 验证区块是否存在
	block, err = client.BlockByHash(ctx, blockHash)
	if err != nil {
		fmt.Printf("区块不存在或查询失败: %v\n", err)
	} else {
		fmt.Printf("区块 #%d 验证成功\n", block.Number().Uint64())
	}

	// 获取指定区块中的交易总数
	count, err := client.TransactionCount(ctx, blockHash)
	if err != nil {
		log.Fatal("获取交易数量失败:", err)
	}
	fmt.Printf("区块中包含 %d 笔交易\n", count)

	// 遍历区块中的交易（限制前3笔）
	for idx := uint(0); idx < count && idx < 3; idx++ {
		tx, err := client.TransactionInBlock(ctx, blockHash, idx)
		if err != nil {
			log.Fatal("获取区块中交易失败:", err)
		}
		fmt.Printf("交易%d: %s\n", idx, tx.Hash().Hex())
	}

	fmt.Println("\n=== 单个交易查询 ===")

	// 通过交易哈希直接查询交易详情
	tx, isPending, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		log.Fatal("通过哈希查询交易失败:", err)
	}

	// 交易状态说明
	fmt.Printf("交易状态: ")
	if isPending {
		fmt.Println("等待确认")
		// isPending = true: 交易已广播到网络，但尚未被打包进区块
		// 还在内存池（mempool）中等待矿工处理，交易状态未确定
	} else {
		fmt.Println("已确认")
		// isPending = false: 交易已被打包进区块，有了确定的区块位置
		// 交易有了最终的执行结果（成功或失败），可以在区块浏览器中查到
	}

	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())
	fmt.Println("=== 测试完成 ===")
}
