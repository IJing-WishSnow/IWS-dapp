package interaction

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"testing"

	"github.com/IJing-WishSnow/dapp/test/interaction/contracts/token" // 导入通过abigen生成的ERC20合约Go绑定代码

	"github.com/ethereum/go-ethereum/accounts/abi/bind" // 提供合约调用选项，如交易发送者、Gas限制等
	"github.com/ethereum/go-ethereum/common"            // 提供以太坊地址和哈希类型处理
	"github.com/ethereum/go-ethereum/ethclient"         // 以太坊客户端，用于连接以太坊节点
)

// TestQueryBalance 测试查询ERC20代币余额及相关信息的功能
// 本测试用例演示了如何与部署在以太坊网络上的ERC20代币合约进行交互
func TestQueryBalance(t *testing.T) {
	// 连接到以太坊Sepolia测试网络
	// 使用Alchemy提供的节点服务，替换为你自己的API密钥
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/4Mb8kv8N7tWzzTDYHAkE3")
	if err != nil {
		log.Fatal("连接以太坊节点失败:", err)
	}

	// IWS代币合约地址（部署在Sepolia测试网上的自定义代币）
	tokenAddress := common.HexToAddress("0xE5aFC41736bBE96cCB912Cb2d2e6BB503979b657")

	// 创建代币合约实例，使用abigen生成的NewToken函数
	// 此实例提供了与ERC20合约交互的所有方法
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal("创建合约实例失败:", err)
	}

	// 要查询余额的钱包地址（需要替换为你实际要查询的地址）
	// address := common.HexToAddress("0x2281cB9267ABAF264c0A4c0dD2e414b4d68cE634")
	address := common.HexToAddress("0x8c8aB9B6178877246B224F8D745A1410C4928373")

	// 查询代币余额 - 返回的是最小单位的余额（基于代币的decimals）
	// 例如：如果decimals=18，返回的是以wei为单位的余额
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal("查询余额失败:", err)
	}

	// 查询代币名称 - ERC20标准可选函数
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal("查询代币名称失败:", err)
	}

	// 查询代币符号 - ERC20标准可选函数
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal("查询代币符号失败:", err)
	}

	// 查询代币精度（小数位数）- ERC20标准可选函数
	// 决定了代币可分割的最小单位，常见值为18（如ETH）
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal("查询代币精度失败:", err)
	}

	// 输出代币基本信息和原始余额
	fmt.Printf("代币名称: %s\n", name)      // 实际输出: "IWS Token"
	fmt.Printf("代币符号: %s\n", symbol)    // 实际输出: "IWS"
	fmt.Printf("小数位数: %v\n", decimals)  // 实际输出: 18
	fmt.Printf("原始余额(最小单位): %s\n", bal) // 原始余额，基于代币精度，如: "999999000000000000000000"

	// 将余额从最小单位转换为可读的十进制格式
	fbal := new(big.Float)
	fbal.SetString(bal.String()) // 将big.Int余额转换为big.Float以便进行小数运算

	// 计算转换系数：10^decimals，用于将最小单位转换为标准单位
	conversionFactor := math.Pow10(int(decimals))

	// 将余额除以转换系数得到实际金额
	// 例如：999999000000000000000000 / 10^18 = 999999.0
	value := new(big.Float).Quo(fbal, big.NewFloat(conversionFactor))

	// 输出格式化后的余额（标准单位）
	fmt.Printf("实际余额: %f\n", value) // 实际输出: 999999.000000
}
