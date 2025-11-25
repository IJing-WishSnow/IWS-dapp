package interaction

import (
	"context"
	_ "embed" // embed 包用于编译时嵌入文件
	"fmt"
	"log"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ==================== 编译时嵌入文件 ====================
// 使用 //go:embed 指令将文件内容嵌入到变量中
// 注意：文件必须在同目录下，文件名要完全匹配

//go:embed contracts/store/Store_sol_Store.abi
var contractABIJSON string // ABI 文件内容（JSON 格式）

//go:embed contracts/store/Store_sol_Store.bin
var contractBytecodeHex string // 字节码文件内容（十六进制字符串）

// ==================== 测试函数：部署合约 ====================
func TestDeployContract1(t *testing.T) {
	// ============ 第一步：连接以太坊节点 ============
	// 使用 Alchemy 提供的 Sepolia 测试网节点
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/4Mb8kv8N7tWzzTDYHAkE3")
	if err != nil {
		log.Fatalf("❌ 连接节点失败: %v", err)
	}
	fmt.Println("✅ 成功连接到 Sepolia 测试网")

	// ============ 第二步：加载私钥 ============
	// 从十六进制字符串加载私钥（注意：生产环境应从环境变量或密钥管理服务读取）
	privateKey, err := crypto.HexToECDSA("ab99f80b034909680a1f840bd37a5f45bda536a2cc484c09dbea504914bcbbd9")
	if err != nil {
		log.Fatalf("❌ 加载私钥失败: %v", err)
	}

	// 从私钥推导出公钥和地址
	publicKey := privateKey.PublicKey
	fromAddress := crypto.PubkeyToAddress(publicKey)
	fmt.Printf("📍 部署地址: %s\n", fromAddress.Hex())

	// ============ 第三步：获取账户 Nonce ============
	// Nonce 是账户发送交易的序号，用于防止重放攻击
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("❌ 获取 Nonce 失败: %v", err)
	}
	fmt.Printf("🔢 当前 Nonce: %d\n", nonce)

	// ============ 第四步：获取当前 Gas 价格 ============
	// Gas Price 是执行交易时愿意支付的每单位 Gas 的价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取 Gas Price 失败: %v", err)
	}
	fmt.Printf("⛽ 建议 Gas Price: %s wei\n", gasPrice.String())

	// ============ 第五步：解析 ABI ============
	// ABI（Application Binary Interface）描述了合约的函数签名
	contractABI, err := abi.JSON(strings.NewReader(contractABIJSON))
	if err != nil {
		log.Fatalf("❌ 解析 ABI 失败: %v", err)
	}
	fmt.Println("✅ ABI 解析成功")

	// ============ 第六步：准备字节码 ============
	// 从嵌入的十六进制字符串转换为字节数组
	// common.FromHex 会自动处理 "0x" 前缀
	bytecode := common.FromHex(contractBytecodeHex)
	fmt.Printf("📦 字节码长度: %d 字节\n", len(bytecode))

	// ============ 第七步：编码构造函数参数 ============
	// 根据你的 ABI，构造函数需要一个 string 类型的 _version 参数
	// abi.Pack("", ...) 中第一个参数为空字符串表示编码构造函数
	versionParam := "v1.0.0" // 合约版本号
	encodedArgs, err := contractABI.Pack("", versionParam)
	if err != nil {
		log.Fatalf("❌ 编码构造函数参数失败: %v", err)
	}
	fmt.Printf("🔧 构造函数参数: version = %s\n", versionParam)

	// ============ 第八步：拼接完整的合约部署数据 ============
	// 部署数据 = 合约字节码 + 编码后的构造函数参数
	data := append(bytecode, encodedArgs...)
	fmt.Printf("📤 完整部署数据长度: %d 字节\n", len(data))

	// ============ 第九步：创建合约部署交易 ============
	// NewContractCreation 创建一个合约部署交易（to 地址为 nil）
	// 参数：nonce, value(转账金额), gasLimit, gasPrice, data
	gasLimit := uint64(3000000) // Gas 上限，设置为 300 万
	value := big.NewInt(0)      // 不向合约转账 ETH
	tx := types.NewContractCreation(nonce, value, gasLimit, gasPrice, data)
	fmt.Println("✅ 交易对象创建成功")

	// ============ 第十步：签名交易 ============
	// 获取链 ID（Sepolia 测试网的 Chain ID 是 11155111）
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取 Chain ID 失败: %v", err)
	}
	fmt.Printf("🔗 Chain ID: %s\n", chainID.String())

	// 使用 EIP-155 签名算法对交易进行签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("❌ 签名交易失败: %v", err)
	}
	fmt.Println("✅ 交易签名成功")

	// ============ 第十一步：发送交易到网络 ============
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("❌ 发送交易失败: %v", err)
	}

	// 打印交易哈希
	txHash := signedTx.Hash().Hex()
	fmt.Printf("\n🚀 合约部署交易已发送!\n")
	fmt.Printf("📋 交易哈希: %s\n", txHash)
	fmt.Printf("🔍 在 Etherscan 查看: https://sepolia.etherscan.io/tx/%s\n\n", txHash)

	// ============ 第十二步：等待交易确认 ============
	fmt.Println("⏳ 等待交易被矿工确认（约 15-30 秒）...")
	receipt, err := waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatalf("❌ 等待交易确认失败: %v", err)
	}

	// ============ 第十三步：检查部署结果 ============
	// Status = 1 表示交易成功，Status = 0 表示交易失败
	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Printf("\n✅ 合约部署成功!\n")
		fmt.Printf("📍 合约地址: %s\n", receipt.ContractAddress.Hex())
		fmt.Printf("⛽ Gas 使用量: %d\n", receipt.GasUsed)
		fmt.Printf("📦 区块高度: %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("🔍 在 Etherscan 查看合约: https://sepolia.etherscan.io/address/%s\n", receipt.ContractAddress.Hex())
	} else {
		log.Fatalf("❌ 合约部署失败! Transaction Status: %d", receipt.Status)
	}
}

// ==================== 辅助函数：等待交易确认 ====================
// 功能：轮询查询交易收据，直到交易被打包进区块
// 参数：client - 以太坊客户端，txHash - 交易哈希
// 返回：交易收据或错误
func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	// 最多等待 60 次（每次 1 秒，总共约 1 分钟）
	for i := 0; i < 60; i++ {
		// 尝试获取交易收据
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			// 收据存在，说明交易已被确认
			return receipt, nil
		}

		// 如果错误不是 "not found"，则返回错误
		if err.Error() != "not found" {
			return nil, fmt.Errorf("查询交易收据失败: %v", err)
		}

		// 打印等待进度
		fmt.Printf(".")
		time.Sleep(1 * time.Second)
	}

	// 超时未确认
	return nil, fmt.Errorf("交易确认超时（已等待 60 秒）")
}
