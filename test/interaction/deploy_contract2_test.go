package interaction

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/IJing-WishSnow/dapp/test/interaction/contracts/store"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestDeployContract2(t *testing.T) {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/4Mb8kv8N7tWzzTDYHAkE3")
	if err != nil {
		log.Fatalf("❌ 连接节点失败: %v", err)
	}
	fmt.Println("✅ 成功连接到 Sepolia 测试网")

	privateKey, err := crypto.HexToECDSA("ab99f80b034909680a1f840bd37a5f45bda536a2cc484c09dbea504914bcbbd9")
	if err != nil {
		log.Fatalf("❌ 加载私钥失败: %v", err)
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("📍 部署地址: %s\n", fromAddress.Hex())

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取 Chain ID 失败: %v", err)
	}
	fmt.Printf("🔗 Chain ID: %s\n", chainID.String())

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("❌ 创建交易签名器失败: %v", err)
	}

	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)

	fmt.Println("✅ 交易签名器创建成功")

	fmt.Println("\n🚀 开始部署合约...")
	versionParam := "v1.0.0"
	fmt.Printf("🔧 构造函数参数: version = %s\n", versionParam)

	address, tx, instance, err := store.DeployStore(auth, client, versionParam)
	if err != nil {
		log.Fatalf("❌ 部署合约失败: %v", err)
	}

	txHash := tx.Hash().Hex()
	fmt.Printf("\n✅ 合约部署交易已发送!\n")
	fmt.Printf("📋 交易哈希: %s\n", txHash)
	fmt.Printf("📍 预计合约地址: %s\n", address.Hex())
	fmt.Printf("🔍 在 Etherscan 查看: https://sepolia.etherscan.io/tx/%s\n\n", txHash)

	fmt.Println("⏳ 等待交易被矿工确认（约 15-30 秒）...")
	receipt, err := waitForReceipt2(client, tx.Hash())
	if err != nil {
		log.Fatalf("❌ 等待交易确认失败: %v", err)
	}

	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Printf("\n✅ 合约部署成功!\n")
		fmt.Printf("📍 合约地址: %s\n", receipt.ContractAddress.Hex())
		fmt.Printf("⛽ Gas 使用量: %d\n", receipt.GasUsed)
		fmt.Printf("💰 Gas 费用: %s ETH\n", weiToEth(new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), tx.GasPrice())))
		fmt.Printf("📦 区块高度: %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("🔍 在 Etherscan 查看合约: https://sepolia.etherscan.io/address/%s\n", receipt.ContractAddress.Hex())

		fmt.Println("\n🧪 测试合约调用...")
		testContractInteraction(instance, auth, client)
	} else {
		log.Fatalf("❌ 合约部署失败! Transaction Status: %d", receipt.Status)
	}
}

func waitForReceipt2(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
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

func weiToEth(wei *big.Int) string {
	fwei := new(big.Float).SetInt(wei)
	fether := new(big.Float).Quo(fwei, big.NewFloat(1e18))
	return fether.Text('f', 6)
}

func testContractInteraction(instance *store.Store, auth *bind.TransactOpts, client *ethclient.Client) {
	version, err := instance.Version(&bind.CallOpts{})
	if err != nil {
		fmt.Printf("⚠️  读取 version 失败: %v\n", err)
		return
	}
	fmt.Printf("📖 读取到的 version: %s\n", version)

	fmt.Println("\n📝 测试写入数据...")

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], "mykey")
	copy(value[:], "myvalue")

	fmt.Printf("🔑 Key: %s\n", string(key[:5]))
	fmt.Printf("💎 Value: %s\n", string(value[:7]))

	tx, err := instance.SetItem(auth, key, value)
	if err != nil {
		fmt.Printf("⚠️  调用 SetItem 失败: %v\n", err)
		return
	}
	fmt.Printf("✅ SetItem 交易已发送: %s\n", tx.Hash().Hex())

	fmt.Print("⏳ 等待交易确认")
	receipt, err := waitForReceipt2(client, tx.Hash())
	if err != nil {
		fmt.Printf("\n⚠️  等待交易确认失败: %v\n", err)
		return
	}

	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Printf("\n✅ SetItem 调用成功!\n")
		fmt.Printf("⛽ Gas 使用量: %d\n", receipt.GasUsed)
		fmt.Printf("📦 区块高度: %d\n", receipt.BlockNumber.Uint64())
	} else {
		fmt.Printf("\n❌ SetItem 调用失败! Status: %d\n", receipt.Status)
		return
	}

	fmt.Println("\n🔍 验证写入的数据...")
	storedValue, err := instance.Items(&bind.CallOpts{}, key)
	if err != nil {
		fmt.Printf("⚠️  读取数据失败: %v\n", err)
		return
	}

	if storedValue == value {
		fmt.Printf("✅ 数据验证成功! 存储的值: %s\n", string(storedValue[:7]))
	} else {
		fmt.Printf("⚠️  数据不匹配!\n")
		fmt.Printf("   期望值: %s\n", string(value[:7]))
		fmt.Printf("   实际值: %s\n", string(storedValue[:7]))
	}
}
