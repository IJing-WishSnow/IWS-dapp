package interaction

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestTokenTransfer(t *testing.T) {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 尝试多个公共节点，增加可靠性
	var client *ethclient.Client
	var err error

	nodeURLs := []string{
		// "https://rpc.sepolia.org",
		"https://ethereum-sepolia-rpc.publicnode.com",
		// "https://sepolia.drpc.org",
	}

	for _, url := range nodeURLs {
		client, err = ethclient.Dial(url)
		if err == nil {
			fmt.Printf("成功连接到节点: %s\n", url)
			break
		}
		fmt.Printf("节点 %s 连接失败: %v\n", url, err)
	}

	if err != nil {
		log.Fatal("所有节点连接都失败")
	}
	defer client.Close()

	// 检查余额和网络状态
	balance, err := client.BalanceAt(ctx, common.HexToAddress("0x8c8aB9B6178877246B224F8D745A1410C4928373"), nil)
	if err != nil {
		log.Fatalf("网络连接测试失败: %v", err)
	}
	fmt.Printf("测试地址余额: %s ETH\n", new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18)).String())

	privateKey, err := crypto.HexToECDSA("ab99f80b034909680a1f840bd37a5f45bda536a2cc484c09dbea504914bcbbd9")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法断言类型: 公钥不是*ecdsa.PublicKey类型")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 带超时的获取 nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("获取 nonce 失败: %v", err)
	}
	fmt.Printf("当前 nonce: %d\n", nonce)

	value := big.NewInt(0)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Gas 价格: %s Gwei\n", new(big.Float).Quo(new(big.Float).SetInt(gasPrice), big.NewFloat(1e9)).String())

	toAddress := common.HexToAddress("0x2281cB9267ABAF264c0A4c0dD2e414b4d68cE634")
	tokenAddress := common.HexToAddress("0xE5aFC41736bBE96cCB912Cb2d2e6BB503979b657")

	// 构建 transfer 函数调用数据
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println("函数选择器:", hexutil.Encode(methodID))

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println("编码地址:", hexutil.Encode(paddedAddress))

	amount := new(big.Int)
	amount.SetString("1000000000000000000", 10) // 减少到 1 个代币，避免余额不足
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println("编码金额:", hexutil.Encode(paddedAmount))

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// 估算 Gas（带超时）
	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		log.Fatalf("估算 Gas 失败: %v", err)
	}
	fmt.Printf("估算 Gas: %d\n", gasLimit)

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("代币转账交易已发送: %s\n", signedTx.Hash().Hex())
}
