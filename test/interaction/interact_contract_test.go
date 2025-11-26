package interaction

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/IJing-WishSnow/IWS-dapp/test/interaction/contracts/storeabi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ==================== 测试函数：连接已部署的合约并交互 ====================
func TestInteractContract(t *testing.T) {
	// ============ 第一步：连接以太坊节点 ============
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/4Mb8kv8N7tWzzTDYHAkE3")
	if err != nil {
		log.Fatalf("❌ 连接节点失败: %v", err)
	}
	fmt.Println("✅ 成功连接到 Sepolia 测试网")

	// ============ 第二步：连接已部署的合约 ============
	// 使用之前部署成功的合约地址
	contractAddr := "0x48Bd8C28155a382d872e4758c11b967303fEDD90"
	storeContract, err := storeabi.NewStoreabi(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatalf("❌ 连接合约失败: %v", err)
	}
	fmt.Printf("✅ 成功连接到合约: %s\n", contractAddr)

	// ============ 第三步：读取合约数据（不需要私钥）============
	fmt.Println("\n📖 读取合约数据...")

	// 读取 version
	version, err := storeContract.Version(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("❌ 读取 version 失败: %v", err)
	}
	fmt.Printf("📌 合约版本: %s\n", version)

	// 读取之前写入的数据
	key := [32]byte{}
	copy(key[:], "mykey")
	storedValue, err := storeContract.Items(&bind.CallOpts{}, key)
	if err != nil {
		log.Fatalf("❌ 读取数据失败: %v", err)
	}
	fmt.Printf("📌 Key 'mykey' 的值: %s\n", string(storedValue[:7]))

	// ============ 第四步：写入新数据（需要私钥）============
	fmt.Println("\n📝 写入新数据到合约...")

	// 加载私钥
	privateKey, err := crypto.HexToECDSA("ab99f80b034909680a1f840bd37a5f45bda536a2cc484c09dbea504914bcbbd9")
	if err != nil {
		log.Fatalf("❌ 加载私钥失败: %v", err)
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("📍 操作地址: %s\n", fromAddress.Hex())

	// 获取链 ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取 Chain ID 失败: %v", err)
	}

	// 创建交易签名器
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("❌ 创建交易签名器失败: %v", err)
	}

	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(100000)

	// 准备新的键值对
	newKey := [32]byte{}
	newValue := [32]byte{}
	copy(newKey[:], "newkey")
	copy(newValue[:], "newvalue")

	fmt.Printf("🔑 Key: %s\n", string(newKey[:6]))
	fmt.Printf("💎 Value: %s\n", string(newValue[:8]))

	// 发送交易
	tx, err := storeContract.SetItem(auth, newKey, newValue)
	if err != nil {
		log.Fatalf("❌ 调用 SetItem 失败: %v", err)
	}
	fmt.Printf("✅ 交易已发送: %s\n", tx.Hash().Hex())
	fmt.Printf("🔍 在 Etherscan 查看: https://sepolia.etherscan.io/tx/%s\n", tx.Hash().Hex())

	// ============ 第五步：等待交易确认 ============
	fmt.Print("⏳ 等待交易确认")
	receipt, err := waitForReceipt3(client, tx.Hash())
	if err != nil {
		log.Fatalf("❌ 等待交易确认失败: %v", err)
	}

	// ============ 第六步：检查交易结果 ============
	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Printf("\n✅ 交易执行成功!\n")
		fmt.Printf("⛽ Gas 使用量: %d\n", receipt.GasUsed)
		fmt.Printf("💰 Gas 费用: %s ETH\n", weiToEth2(new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), tx.GasPrice())))
		fmt.Printf("📦 区块高度: %d\n", receipt.BlockNumber.Uint64())
	} else {
		log.Fatalf("❌ 交易执行失败! Status: %d", receipt.Status)
	}

	// ============ 第七步：验证写入的数据 ============
	fmt.Println("\n🔍 验证新写入的数据...")
	verifyValue, err := storeContract.Items(&bind.CallOpts{}, newKey)
	if err != nil {
		log.Fatalf("❌ 读取数据失败: %v", err)
	}

	if verifyValue == newValue {
		fmt.Printf("✅ 数据验证成功! 存储的值: %s\n", string(verifyValue[:8]))
	} else {
		fmt.Printf("⚠️  数据不匹配!\n")
		fmt.Printf("   期望值: %s\n", string(newValue[:8]))
		fmt.Printf("   实际值: %s\n", string(verifyValue[:8]))
	}

	fmt.Println("\n🎉 合约交互测试完成!")
}

// ==================== 辅助函数：等待交易确认 ====================
func waitForReceipt3(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
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

// ==================== 辅助函数：Wei 转 ETH ====================
func weiToEth2(wei *big.Int) string {
	fwei := new(big.Float).SetInt(wei)
	fether := new(big.Float).Quo(fwei, big.NewFloat(1e18))
	return fether.Text('f', 6)
}
