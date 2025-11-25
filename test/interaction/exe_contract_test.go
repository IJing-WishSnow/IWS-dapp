package interaction

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	contractAddr = "0x48Bd8C28155a382d872e4758c11b967303fEDD90"
)

func TestExeContract(t *testing.T) {
	// ============ 第一步：连接以太坊节点 ============
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/4Mb8kv8N7tWzzTDYHAkE3")
	if err != nil {
		log.Fatalf("❌ 连接节点失败: %v", err)
	}
	fmt.Println("✅ 成功连接到 Sepolia 测试网")

	// ============ 第二步：加载私钥 ============
	privateKey, err := crypto.HexToECDSA("ab99f80b034909680a1f840bd37a5f45bda536a2cc484c09dbea504914bcbbd9")
	if err != nil {
		log.Fatalf("❌ 加载私钥失败: %v", err)
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("📍 操作地址: %s\n", fromAddress.Hex())

	// ============ 第三步：手动构造 setItem 调用数据 ============
	fmt.Println("\n📝 准备调用 setItem...")

	// 计算函数选择器（函数签名的 Keccak256 哈希的前 4 字节）
	methodSignature := []byte("setItem(bytes32,bytes32)")
	methodSelector := crypto.Keccak256(methodSignature)[:4]

	// 准备参数
	var key [32]byte
	var value [32]byte
	copy(key[:], []byte("demo_save_key_no_use_abi"))
	copy(value[:], []byte("demo_save_value_no_use_abi_11111"))

	fmt.Printf("🔑 Key: %s\n", string(key[:24]))
	fmt.Printf("💎 Value: %s\n", string(value[:32]))

	// 组合调用数据：函数选择器 + 参数
	var input []byte
	input = append(input, methodSelector...)
	input = append(input, key[:]...)
	input = append(input, value[:]...)

	// ============ 第四步：构造并发送交易 ============
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("❌ 获取 nonce 失败: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取 gas price 失败: %v", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取 chain ID 失败: %v", err)
	}

	// 创建交易
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(contractAddr),
		big.NewInt(0),
		uint64(200000),
		gasPrice,
		input, // 使用手动构造的 calldata
	)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("❌ 签名交易失败: %v", err)
	}

	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("❌ 发送交易失败: %v", err)
	}

	fmt.Printf("\n✅ 交易已发送: %s\n", signedTx.Hash().Hex())
	fmt.Printf("🔍 在 Etherscan 查看: https://sepolia.etherscan.io/tx/%s\n", signedTx.Hash().Hex())

	// ============ 第五步：等待交易确认 ============
	fmt.Print("⏳ 等待交易确认")
	receipt, err := waitForReceipt4(client, signedTx.Hash())
	if err != nil {
		log.Fatalf("❌ 等待交易确认失败: %v", err)
	}

	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Printf("\n✅ 交易执行成功!\n")
		fmt.Printf("⛽ Gas 使用量: %d\n", receipt.GasUsed)
		fmt.Printf("📦 区块高度: %d\n", receipt.BlockNumber.Uint64())
	} else {
		log.Fatalf("❌ 交易执行失败! Status: %d", receipt.Status)
	}

	// ============ 第六步：手动构造查询数据 ============
	fmt.Println("\n🔍 验证写入的数据...")

	// 计算 items 函数的选择器
	itemsSignature := []byte("items(bytes32)")
	itemsSelector := crypto.Keccak256(itemsSignature)[:4]

	// 组合查询数据
	var callInput []byte
	callInput = append(callInput, itemsSelector...)
	callInput = append(callInput, key[:]...)

	// 构造调用消息
	to := common.HexToAddress(contractAddr)
	callMsg := ethereum.CallMsg{
		To:   &to,
		Data: callInput,
	}

	// ============ 第七步：调用合约查询 ============
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatalf("❌ 调用合约失败: %v", err)
	}

	// ============ 第八步：手动解析返回值 ============
	var unpacked [32]byte
	copy(unpacked[:], result)

	// 验证数据
	if unpacked == value {
		fmt.Printf("✅ 数据验证成功! 存储的值与原始值相同\n")
		fmt.Printf("📌 存储的值: %s\n", string(unpacked[:32]))
	} else {
		fmt.Printf("⚠️  数据不匹配!\n")
		fmt.Printf("   期望值: %s\n", string(value[:32]))
		fmt.Printf("   实际值: %s\n", string(unpacked[:32]))
	}

	fmt.Println("\n🎉 操作完成!")
}

// ==================== 辅助函数：等待交易确认 ====================
func waitForReceipt4(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for i := 0; i < 60; i++ {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}

		if err.Error() != "not found" {
			return nil, fmt.Errorf("查询交易收据失败: %v", err)
		}

		fmt.Printf(".")
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("交易确认超时（已等待 60 秒）")
}
