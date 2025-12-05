package fork

import (
	"context"
	"math/big"
	"testing"
	"time"
)

// 测试主入口
func TestForkStatus(t *testing.T) {
	t.Run("连接测试", testConnection)
	t.Run("网络信息测试", testNetworkInfo)
	t.Run("区块信息测试", testBlockInfo)
	t.Run("账户余额测试", testAccountBalances)
	t.Run("BSC合约测试", testBSCContracts)
}

// 测试1: 连接是否成功
func testConnection(t *testing.T) {
	cfg := DefaultConfig

	// 创建客户端
	cli, err := NewClient(cfg.RPCURL, time.Duration(cfg.TestTimeout)*time.Second)
	if err != nil {
		t.Fatalf("❌ 连接失败: %v", err)
	}
	defer cli.Close()

	t.Log("✅ 连接本地节点成功")

	// 测试 ping
	blockNumber, err := cli.GetCurrentBlockNumber()
	if err != nil {
		t.Fatalf("❌ 获取区块号失败: %v", err)
	}

	t.Logf("📦 当前区块号: %d", blockNumber)
}

// 测试2: 网络信息
func testNetworkInfo(t *testing.T) {
	cli, err := NewClient(DefaultConfig.RPCURL,
		time.Duration(DefaultConfig.TestTimeout)*time.Second)
	if err != nil {
		t.Fatalf("❌ 连接失败: %v", err)
	}
	defer cli.Close()

	// 获取网络ID
	chainID, err := cli.GetNetworkID()
	if err != nil {
		t.Fatalf("❌ 获取网络ID失败: %v", err)
	}

	t.Logf("🌐 网络ID: %d", chainID)

	// 判断网络类型
	switch chainID.Uint64() {
	case DefaultConfig.ChainIDs["hardhat"].Uint64():
		t.Log("ℹ️  检测到 Hardhat 本地网络")
		// 这里可以标记测试结果为警告，但不是失败
		t.Log("⚠️  警告: 这可能是纯本地节点，未分叉到测试网")
	case DefaultConfig.ChainIDs["bsc_test"].Uint64():
		t.Log("✅ 检测到 BSC 测试网分叉")
	case DefaultConfig.ChainIDs["bsc_main"].Uint64():
		t.Log("ℹ️  检测到 BSC 主网分叉")
	case DefaultConfig.ChainIDs["sepolia"].Uint64():
		t.Log("ℹ️  检测到 Sepolia 测试网分叉")
	default:
		t.Logf("ℹ️  未知网络 (ID: %d)", chainID)
	}

	// 如果是分叉测试，期望是 BSC 测试网
	expectedChainID := DefaultConfig.ChainIDs["bsc_test"]
	if chainID.Cmp(expectedChainID) != 0 {
		t.Logf("⚠️  注意: 期望链ID %d (BSC测试网)，实际得到 %d",
			expectedChainID, chainID)
	}
}

// 测试3: 区块信息
func testBlockInfo(t *testing.T) {
	cli, err := NewClient(DefaultConfig.RPCURL,
		time.Duration(DefaultConfig.TestTimeout)*time.Second)
	if err != nil {
		t.Fatalf("❌ 连接失败: %v", err)
	}
	defer cli.Close()

	blockNumber, err := cli.GetCurrentBlockNumber()
	if err != nil {
		t.Fatalf("❌ 获取区块号失败: %v", err)
	}

	t.Logf("📊 当前区块高度: %d", blockNumber)

	// 判断区块高度
	if blockNumber == 0 {
		t.Log("⚠️  警告: 区块高度为 0，可能是全新的本地节点")
	} else if blockNumber < 1000 {
		t.Log("ℹ️  区块高度较低，可能是刚启动的节点")
	} else if blockNumber > 30000000 {
		t.Log("✅ 区块高度正常，符合 BSC 测试网范围")
	}

	// 获取 Gas 价格
	ctx := context.Background()
	gasPrice, err := cli.SuggestGasPrice(ctx)
	if err != nil {
		t.Logf("⚠️  获取 Gas 价格失败: %v", err)
	} else {
		gasPriceGwei := new(big.Int).Div(gasPrice, big.NewInt(1e9))
		t.Logf("⛽ 建议 Gas 价格: %d Gwei", gasPriceGwei)
	}
}

// 测试4: 账户余额
func testAccountBalances(t *testing.T) {
	cli, err := NewClient(DefaultConfig.RPCURL,
		time.Duration(DefaultConfig.TestTimeout)*time.Second)
	if err != nil {
		t.Fatalf("❌ 连接失败: %v", err)
	}
	defer cli.Close()

	// 测试前几个账户
	for i, addr := range DefaultConfig.TestAddresses[:2] {
		balance, err := cli.GetBalance(addr)
		if err != nil {
			t.Logf("⚠️  获取账户 %d 余额失败: %v", i, err)
			continue
		}

		ethBalance := WeiToEther(balance)
		t.Logf("💰 账户%d (%s): %s ETH",
			i,
			addr[:8]+"...",
			ethBalance.Text('f', 4))

		// 检查是否有足够的测试 ETH
		minBalance := new(big.Int).Mul(big.NewInt(1), big.NewInt(1e18)) // 1 ETH
		if balance.Cmp(minBalance) < 0 {
			t.Logf("⚠️  账户%d 余额较低: %s ETH", i, ethBalance.Text('f', 4))
		}
	}
}

// 测试5: BSC 测试网合约
func testBSCContracts(t *testing.T) {
	cli, err := NewClient(DefaultConfig.RPCURL,
		time.Duration(DefaultConfig.TestTimeout)*time.Second)
	if err != nil {
		t.Fatalf("❌ 连接失败: %v", err)
	}
	defer cli.Close()

	// 检查 BSC 测试网已知合约
	bscContracts := []struct {
		name    string
		address string
	}{
		{"WBNB", "0xae13d989daC2f0dEbFf460aC112a837C89BAa7cd"},
		{"PancakeSwap Router", "0x9Ac64Cc6e4415144C455BD8E4837Fea55603e5c3"},
	}

	contractsFound := 0
	for _, contract := range bscContracts {
		hasCode, codeSize, err := cli.HasContractCode(contract.address)
		if err != nil {
			t.Logf("⚠️  检查合约 %s 失败: %v", contract.name, err)
			continue
		}

		if hasCode {
			t.Logf("✅ 检测到 %s 合约 (代码大小: %d 字节)",
				contract.name, codeSize)
			contractsFound++
		} else {
			t.Logf("❌ 未检测到 %s 合约", contract.name)
		}
	}

	// 判断分叉是否成功
	chainID, err := cli.GetNetworkID()
	if err == nil && chainID.Cmp(DefaultConfig.ChainIDs["bsc_test"]) == 0 {
		if contractsFound >= 1 {
			t.Log("🎉 BSC 测试网分叉成功！")
		} else {
			t.Log("⚠️  检测到 BSC 测试网链ID，但未找到标准合约")
			t.Log("   可能是分叉的区块较早，这些合约还未部署")
		}
	}
}
