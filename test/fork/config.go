package fork

import "math/big"

// 测试配置
type TestConfig struct {
	RPCURL      string
	TestTimeout int // 秒
	ChainIDs    map[string]*big.Int
	TestAddresses []string
}

// 默认配置
var DefaultConfig = TestConfig{
	RPCURL:      "http://127.0.0.1:8545",
	TestTimeout: 10,
	ChainIDs: map[string]*big.Int{
		"hardhat":   big.NewInt(31337),
		"bsc_test":  big.NewInt(97),
		"bsc_main":  big.NewInt(56),
		"sepolia":   big.NewInt(11155111),
	},
	TestAddresses: []string{
		"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", // Hardhat 账户0
		"0x70997970C51812dc3A010C7d01b50e0d17dc79C8", // Hardhat 账户1
		"0xae13d989daC2f0dEbFf460aC112a837C89BAa7cd", // BSC测试网 WBNB
		"0x9Ac64Cc6e4415144C455BD8E4837Fea55603e5c3", // PancakeSwap 路由器
	},
}