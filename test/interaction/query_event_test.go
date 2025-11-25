package interaction

import (
	"context"
	_ "embed" // 使用 embed 包嵌入文件
	"fmt"
	"log"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types" // 添加缺失的导入
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ==================== 测试函数：查询 ERC20 合约事件 ====================
func TestQueryEvent(t *testing.T) {
	// ============ 第一步：连接以太坊节点 ============
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/4Mb8kv8N7tWzzTDYHAkE3")
	if err != nil {
		log.Fatalf("❌ 连接节点失败: %v", err)
	}
	fmt.Println("✅ 成功连接到 Sepolia 测试网")

	// ============ 第二步：设置合约地址 ============
	contractAddress := common.HexToAddress("0xE5aFC41736bBE96cCB912Cb2d2e6BB503979b657")
	fmt.Printf("📍 合约地址: %s\n", contractAddress.Hex())

	// ============ 第三步：解析 ERC20 合约 ABI ============
	contractABI, err := abi.JSON(strings.NewReader(erc20ABIJSON))
	if err != nil {
		log.Fatalf("❌ 解析 ABI 失败: %v", err)
	}
	fmt.Println("✅ ABI 解析成功")

	// ============ 第四步：自动搜索包含事件的区块范围 ============
	fmt.Println("🔍 开始自动搜索包含事件的区块范围...")
	foundEvents := findEventsInRange(client, contractAddress, contractABI)

	if !foundEvents {
		fmt.Println("❌ 在搜索范围内未找到任何事件")
		return
	}

	// ============ 第五步：完成查询 ============
	fmt.Println("\n🎉 事件查询完成!")
}

// ==================== 自动搜索事件函数 ====================
func findEventsInRange(client *ethclient.Client, contractAddress common.Address, contractABI abi.ABI) bool {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Printf("❌ 获取最新区块失败: %v", err)
		return false
	}

	fmt.Printf("📦 当前最新区块: %d\n", header.Number.Uint64())
	fmt.Printf("🎯 开始搜索，最大搜索范围: 1000 个区块\n")

	// 计算 Transfer 事件签名哈希
	transferEventSignature := []byte("Transfer(address,address,uint256)")
	transferEventHash := crypto.Keccak256Hash(transferEventSignature)

	// 从当前区块往前搜索，每次查询10个区块
	for offset := uint64(0); offset < 1000; offset += 10 {
		// 计算当前查询的区块范围
		currentToBlock := new(big.Int).Sub(header.Number, big.NewInt(int64(offset)))
		currentFromBlock := new(big.Int).Sub(currentToBlock, big.NewInt(9))

		// 确保 FromBlock 不小于 0
		if currentFromBlock.Sign() < 0 {
			currentFromBlock = big.NewInt(0)
		}

		// 如果 FromBlock 大于 ToBlock，说明已经搜索完所有区块
		if currentFromBlock.Cmp(currentToBlock) > 0 {
			fmt.Println("🔚 已搜索到创世区块，停止搜索")
			break
		}

		fmt.Printf("\n🔄 搜索进度: 偏移 %d 区块 (搜索范围: %d ~ %d)\n",
			offset, currentFromBlock.Uint64(), currentToBlock.Uint64())

		query := ethereum.FilterQuery{
			FromBlock: currentFromBlock,
			ToBlock:   currentToBlock,
			Addresses: []common.Address{contractAddress},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Printf("⚠️  查询区块 %d ~ %d 失败: %v",
				currentFromBlock.Uint64(), currentToBlock.Uint64(), err)
			continue
		}

		if len(logs) > 0 {
			fmt.Printf("🎯 在区块 %d ~ %d 中找到 %d 个事件\n",
				currentFromBlock.Uint64(), currentToBlock.Uint64(), len(logs))

			// 处理找到的所有事件
			processEvents(logs, contractABI, transferEventHash)
			return true
		} else {
			fmt.Printf("📭 区块 %d ~ %d 中没有事件\n",
				currentFromBlock.Uint64(), currentToBlock.Uint64())
		}
	}

	return false
}

// ==================== 处理事件函数 ====================
func processEvents(logs []types.Log, contractABI abi.ABI, transferEventHash common.Hash) {
	fmt.Printf("\n📊 开始处理 %d 个事件...\n", len(logs))

	for i, vLog := range logs {
		fmt.Printf("\n=== 事件 #%d ===\n", i+1)

		// 检查事件类型
		if len(vLog.Topics) > 0 {
			eventSignature := vLog.Topics[0]

			// ============ 处理 Transfer 事件 ============
			if eventSignature == transferEventHash {
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
					continue
				}

				// 从 Topics 中获取 indexed 参数
				if len(vLog.Topics) >= 3 {
					event.From = common.BytesToAddress(vLog.Topics[1].Bytes())
					event.To = common.BytesToAddress(vLog.Topics[2].Bytes())
				}

				// ============ 完整的事件日志信息输出 ============
				fmt.Println("📋 === 完整事件日志信息 ===")

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

				// 解析 Data 字段的详细结构
				if len(vLog.Data) >= 32 {
					valueBytes := vLog.Data[:32]
					value := new(big.Int).SetBytes(valueBytes)
					fmt.Printf("   💰 解析金额: %s (原始值)\n", value.String())

					// 假设代币有 18 位小数（常见情况）
					decimalValue := new(big.Float).SetInt(value)
					decimalValue.Quo(decimalValue, big.NewFloat(1e18))
					fmt.Printf("   💎 格式化金额: %s 代币\n", decimalValue.Text('f', 6))
				}

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

				// 6. 原始日志结构（用于高级调试）
				fmt.Println("\n🔧 原始日志结构（调试用）:")
				fmt.Printf("   📍 Address: %s\n", vLog.Address.Hex())
				fmt.Printf("   🏷️  BlockHash: %s\n", vLog.BlockHash.Hex())
				fmt.Printf("   🔢 BlockNumber: %d\n", vLog.BlockNumber)
				fmt.Printf("   📋 TxHash: %s\n", vLog.TxHash.Hex())
				fmt.Printf("   📊 TxIndex: %d\n", vLog.TxIndex)
				fmt.Printf("   🔍 Index: %d\n", vLog.Index)
				fmt.Printf("   ❌ Removed: %t\n", vLog.Removed)

				fmt.Println("🎉 === 事件日志信息输出完成 ===")

			} else {
				// ============ 处理其他类型事件 ============
				fmt.Printf("❓ 未知事件类型，签名: %s\n", eventSignature.Hex())

				// 尝试识别其他常见 ERC20 事件
				approvalEventSignature := []byte("Approval(address,address,uint256)")
				approvalEventHash := crypto.Keccak256Hash(approvalEventSignature)

				if eventSignature == approvalEventHash {
					fmt.Println("✅ 检测到 Approval 事件")

					event := struct {
						Owner   common.Address
						Spender common.Address
						Value   *big.Int
					}{}

					err := contractABI.UnpackIntoInterface(&event, "Approval", vLog.Data)
					if err != nil {
						log.Printf("⚠️  解析 Approval 事件失败: %v", err)
						continue
					}

					if len(vLog.Topics) >= 3 {
						event.Owner = common.BytesToAddress(vLog.Topics[1].Bytes())
						event.Spender = common.BytesToAddress(vLog.Topics[2].Bytes())
					}

					fmt.Printf("👤 授权方: %s\n", event.Owner.Hex())
					fmt.Printf("👥 被授权方: %s\n", event.Spender.Hex())
					fmt.Printf("💸 授权金额: %s\n", event.Value.String())
				}
			}
		}
	}
}
