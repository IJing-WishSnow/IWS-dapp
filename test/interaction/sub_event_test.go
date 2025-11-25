package interaction

import (
	"context"
	_ "embed" // 使用 embed 包嵌入文件
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ==================== 测试函数：实时监听 ERC20 合约事件 ====================
func TestSubscribeERC20Events(t *testing.T) {
	// ============ 第一步：连接以太坊 WebSocket 节点 ============
	fmt.Println("🔌 正在连接 WebSocket...")
	client, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/4Mb8kv8N7tWzzTDYHAkE3")
	if err != nil {
		log.Fatalf("❌ 连接 WebSocket 节点失败: %v", err)
	}
	fmt.Println("✅ 成功连接到 Sepolia 测试网 WebSocket")

	// ============ 第二步：设置 ERC20 合约地址 ============
	contractAddress := common.HexToAddress("0xE5aFC41736bBE96cCB912Cb2d2e6BB503979b657")
	fmt.Printf("📍 合约地址: %s\n", contractAddress.Hex())

	// ============ 第三步：解析 ERC20 合约 ABI ============
	contractABI, err := abi.JSON(strings.NewReader(erc20ABIJSON))
	if err != nil {
		log.Fatalf("❌ 解析 ABI 失败: %v", err)
	}
	fmt.Println("✅ ABI 解析成功")

	// ============ 第四步：计算事件签名哈希 ============
	transferEventSignature := []byte("Transfer(address,address,uint256)")
	transferEventHash := crypto.Keccak256Hash(transferEventSignature)
	fmt.Printf("🔑 Transfer 事件签名哈希: %s\n", transferEventHash.Hex())

	approvalEventSignature := []byte("Approval(address,address,uint256)")
	approvalEventHash := crypto.Keccak256Hash(approvalEventSignature)
	fmt.Printf("🔑 Approval 事件签名哈希: %s\n", approvalEventHash.Hex())

	// ============ 第五步：创建事件订阅过滤器 ============
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}
	fmt.Println("✅ 过滤器设置完成")

	// ============ 第六步：创建事件日志通道并订阅 ============
	logs := make(chan types.Log)
	fmt.Println("📡 正在订阅事件日志...")

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("❌ 订阅事件失败: %v", err)
	}
	defer sub.Unsubscribe()

	fmt.Println("✅ 事件订阅成功，开始实时监听...")
	// fmt.Println("⏳ 等待新事件产生（测试将在30秒后自动结束）...")
	fmt.Println("⏳ 等待新事件产生（按 Ctrl+C 手动停止测试）...")

	// ============ 第七步：实时监听事件 ============
	// 创建手动停止通道
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("❌ 订阅错误: %v", err)
		case vLog := <-logs:
			fmt.Printf("\n🎉 收到新事件! 时间: %s\n", time.Now().Format("15:04:05"))
			processRealtimeEvent(vLog, contractABI, transferEventHash, approvalEventHash)
		case <-stopChan:
			fmt.Println("\n⏰ 收到停止信号，测试结束")
			return
		}
	}
	// timeout := time.After(30 * time.Second) // 30秒后自动结束测试

	// for {
	// 	select {
	// 	case err := <-sub.Err():
	// 		log.Fatalf("❌ 订阅错误: %v", err)
	// 	case vLog := <-logs:
	// 		fmt.Printf("\n🎉 收到新事件! 时间: %s\n", time.Now().Format("15:04:05"))
	// 		processRealtimeEvent(vLog, contractABI, transferEventHash, approvalEventHash)
	// 	case <-timeout:
	// 		fmt.Println("\n⏰ 测试时间结束，停止监听")
	// 		return
	// 	}
	// }
}

// ==================== 处理实时事件函数 ====================
func processRealtimeEvent(vLog types.Log, contractABI abi.ABI, transferEventHash common.Hash, approvalEventHash common.Hash) {
	// 检查事件类型
	if len(vLog.Topics) > 0 {
		eventSignature := vLog.Topics[0]

		// ============ 使用带标签的 switch 语句处理不同事件类型 ============
		switch eventSignature {
		case transferEventHash:
			processTransferEvent(vLog, contractABI)
		case approvalEventHash:
			processApprovalEvent(vLog, contractABI)
		default:
			processUnknownEvent(vLog)
		}
	}
}

// ==================== 处理 Transfer 事件函数 ====================
func processTransferEvent(vLog types.Log, contractABI abi.ABI) {
	fmt.Println("💰 检测到 Transfer 事件")

	// 解析 Transfer 事件参数
	event := struct {
		From  common.Address
		To    common.Address
		Value *big.Int
	}{}

	err := contractABI.UnpackIntoInterface(&event, "Transfer", vLog.Data)
	if err != nil {
		log.Printf("⚠️  解析 Transfer 事件数据失败: %v", err)
		return
	}

	// 从 Topics 中获取 indexed 参数
	if len(vLog.Topics) >= 3 {
		event.From = common.BytesToAddress(vLog.Topics[1].Bytes())
		event.To = common.BytesToAddress(vLog.Topics[2].Bytes())
	}

	// ============ 完整的事件日志信息输出 ============
	fmt.Println("📋 === 实时事件日志信息 ===")

	// 1. 基础信息
	fmt.Println("📍 基础信息:")
	fmt.Printf("   📍 合约地址: %s\n", vLog.Address.Hex())
	fmt.Printf("   📦 区块哈希: %s\n", vLog.BlockHash.Hex())
	fmt.Printf("   🔢 区块高度: %d\n", vLog.BlockNumber)
	fmt.Printf("   📋 交易哈希: %s\n", vLog.TxHash.Hex())
	fmt.Printf("   📊 日志索引: %d\n", vLog.Index)
	fmt.Printf("   🔍 交易索引: %d\n", vLog.TxIndex)
	if vLog.Removed {
		fmt.Printf("   ⚠️  日志状态: 已移除（由于链重组）\n")
	} else {
		fmt.Printf("   ✅ 日志状态: 有效\n")
	}

	// 2. Topics 详细信息
	fmt.Println("\n🔖 Topics 详细信息:")
	fmt.Printf("   📊 Topics 数量: %d\n", len(vLog.Topics))
	for i, topic := range vLog.Topics {
		switch i {
		case 0:
			fmt.Printf("   🔑 Topic[%d] (事件签名): %s\n", i, topic.Hex())
			fmt.Printf("       📝 含义: Transfer(address,address,uint256) 的 Keccak256 哈希\n")
		case 1:
			fmt.Printf("   👤 Topic[%d] (发送方): %s\n", i, topic.Hex())
			fmt.Printf("       📝 解析地址: %s\n", common.BytesToAddress(topic.Bytes()).Hex())
		case 2:
			fmt.Printf("   👥 Topic[%d] (接收方): %s\n", i, topic.Hex())
			fmt.Printf("       📝 解析地址: %s\n", common.BytesToAddress(topic.Bytes()).Hex())
		default:
			fmt.Printf("   ❓ Topic[%d] (未知): %s\n", i, topic.Hex())
		}
	}

	// 3. Data 字段详细信息
	fmt.Println("\n📄 Data 字段详细信息:")
	fmt.Printf("   📏 Data 长度: %d 字节\n", len(vLog.Data))
	fmt.Printf("   🔢 原始数据: %s\n", common.Bytes2Hex(vLog.Data))

	// 4. 事件参数汇总
	fmt.Println("\n📊 事件参数汇总:")
	fmt.Printf("   👤 发送方 (from): %s\n", event.From.Hex())
	fmt.Printf("   👥 接收方 (to): %s\n", event.To.Hex())
	fmt.Printf("   💸 转账金额 (value): %s\n", event.Value.String())

	// 格式化金额显示
	formattedValue := new(big.Float).SetInt(event.Value)
	formattedValue.Quo(formattedValue, big.NewFloat(1e18))
	fmt.Printf("   🎯 格式化金额: %s 代币\n", formattedValue.Text('f', 6))

	// 5. 相关链接（用于调试）
	fmt.Println("\n🔗 相关链接:")
	fmt.Printf("   🌐 Etherscan 交易: https://sepolia.etherscan.io/tx/%s\n", vLog.TxHash.Hex())
	fmt.Printf("   📦 Etherscan 区块: https://sepolia.etherscan.io/block/%d\n", vLog.BlockNumber)
	fmt.Printf("   🏢 Etherscan 合约: https://sepolia.etherscan.io/address/%s\n", vLog.Address.Hex())

	fmt.Println("🎉 === 实时事件日志信息输出完成 ===")
}

// ==================== 处理 Approval 事件函数 ====================
func processApprovalEvent(vLog types.Log, contractABI abi.ABI) {
	fmt.Println("✅ 检测到 Approval 事件")

	event := struct {
		Owner   common.Address
		Spender common.Address
		Value   *big.Int
	}{}

	err := contractABI.UnpackIntoInterface(&event, "Approval", vLog.Data)
	if err != nil {
		log.Printf("⚠️  解析 Approval 事件失败: %v", err)
		return
	}

	if len(vLog.Topics) >= 3 {
		event.Owner = common.BytesToAddress(vLog.Topics[1].Bytes())
		event.Spender = common.BytesToAddress(vLog.Topics[2].Bytes())
	}

	fmt.Println("📋 === 实时 Approval 事件信息 ===")
	fmt.Printf("👤 授权方: %s\n", event.Owner.Hex())
	fmt.Printf("👥 被授权方: %s\n", event.Spender.Hex())
	fmt.Printf("💸 授权金额: %s\n", event.Value.String())
	fmt.Println("🎉 === Approval 事件信息输出完成 ===")
}

// ==================== 处理未知事件函数 ====================
func processUnknownEvent(vLog types.Log) {
	fmt.Printf("❓ 未知事件类型，签名: %s\n", vLog.Topics[0].Hex())
	fmt.Printf("📦 区块高度: %d\n", vLog.BlockNumber)
	fmt.Printf("📋 交易哈希: %s\n", vLog.TxHash.Hex())
	fmt.Println("📋 事件 Topics:")
	for j, topic := range vLog.Topics {
		fmt.Printf("  Topic[%d]: %s\n", j, topic.Hex())
	}
	fmt.Println()
}

// 除了从查询事件和订阅事件能够获得合约事件，还可以从交易收据（TransactionReceipt）的 Logs 字段获取合约事件数据。
