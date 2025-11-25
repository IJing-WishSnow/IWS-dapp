package interaction

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestETHTransfer(t *testing.T) {
	// 连接到以太坊测试网络节点
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/0c994b2c7e7d4226bde6c128e2d2f2c1")
	if err != nil {
		log.Fatalf("无法连接以太坊节点: %v", err) // 连接失败直接终止程序
	}

	// 从16进制字符串加载私钥（测试环境使用，生产环境需安全存储）
	privateKey, err := crypto.HexToECDSA("ab99f80b034909680a1f840bd37a5f45bda536a2cc484c09dbea504914bcbbd9")
	if err != nil {
		log.Fatal(err)
	}

	// 从私钥推导出公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法断言类型: 公钥不是*ecdsa.PublicKey类型")
	}

	// 从公钥生成发送者地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取账户的待处理交易序号，防止重放攻击
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// 设置转账金额：1 ETH = 10^18 wei
	value := big.NewInt(20000000000000000)

	// 设置Gas限制：普通ETH转账的标准Gas用量
	gasLimit := uint64(21000)

	// 获取网络推荐的Gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 设置接收者地址
	toAddress := common.HexToAddress("0x2281cB9267ABAF264c0A4c0dD2e414b4d68cE634")

	// 普通ETH转账无附加数据
	var data []byte

	// 创建交易对象
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	// 获取当前网络的链ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 使用私钥对交易进行签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 将签名后的交易广播到网络
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	// 输出交易哈希，可用于在区块浏览器中查询交易状态
	fmt.Printf("交易已发送: %s", signedTx.Hash().Hex())
}
