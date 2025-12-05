package fork

import (
	"context"
	"testing"
	"time"
)

// 测试网络响应速度
func TestNetworkLatency(t *testing.T) {
	cli, err := NewClient(DefaultConfig.RPCURL,
		time.Duration(DefaultConfig.TestTimeout)*time.Second)
	if err != nil {
		t.Fatalf("❌ 连接失败: %v", err)
	}
	defer cli.Close()

	start := time.Now()

	// 执行多个请求测试延迟
	operations := []struct {
		name string
		fn   func() error
	}{
		{
			name: "获取区块号",
			fn: func() error {
				_, err := cli.GetCurrentBlockNumber()
				return err
			},
		},
		{
			name: "获取网络ID",
			fn: func() error {
				_, err := cli.GetNetworkID()
				return err
			},
		},
		{
			name: "获取Gas价格",
			fn: func() error {
				_, err := cli.SuggestGasPrice(context.Background())
				return err
			},
		},
	}

	for _, op := range operations {
		opStart := time.Now()
		if err := op.fn(); err != nil {
			t.Logf("⚠️  %s 失败: %v", op.name, err)
		} else {
			duration := time.Since(opStart)
			t.Logf("⏱️  %s 耗时: %v", op.name, duration)

			// 检查是否超时（超过1秒为慢）
			if duration > time.Second {
				t.Logf("⚠️  %s 响应较慢", op.name)
			}
		}
	}

	totalTime := time.Since(start)
	t.Logf("🕐 总测试耗时: %v", totalTime)
}

// 测试交易功能
func TestTransactionCapability(t *testing.T) {
	cli, err := NewClient(DefaultConfig.RPCURL,
		time.Duration(DefaultConfig.TestTimeout)*time.Second)
	if err != nil {
		t.Fatalf("❌ 连接失败: %v", err)
	}
	defer cli.Close()

	ctx := context.Background()

	// 获取链ID
	chainID, err := cli.GetNetworkID()
	if err != nil {
		t.Fatalf("❌ 获取链ID失败: %v", err)
	}

	t.Logf("🔗 链ID: %d", chainID)

	// 检查是否支持 EIP-1559
	latestBlock, err := cli.BlockByNumber(ctx, nil)
	if err != nil {
		t.Logf("⚠️  获取最新区块失败: %v", err)
		return
	}

	if latestBlock.BaseFee() != nil {
		t.Log("✅ 网络支持 EIP-1559 (基础费用)")
	} else {
		t.Log("ℹ️  网络不支持 EIP-1559")
	}
}

// 测试同步状态
func TestSyncStatus(t *testing.T) {
	cli, err := NewClient(DefaultConfig.RPCURL,
		time.Duration(DefaultConfig.TestTimeout)*time.Second)
	if err != nil {
		t.Fatalf("❌ 连接失败: %v", err)
	}
	defer cli.Close()

	// 获取同步状态
	syncProgress, err := cli.SyncProgress(context.Background())
	if err != nil {
		t.Logf("⚠️  获取同步状态失败: %v", err)
		return
	}

	if syncProgress != nil {
		percentage := float64(syncProgress.CurrentBlock) / float64(syncProgress.HighestBlock) * 100
		t.Logf("🔄 节点正在同步: %.2f%% 完成", percentage)
		t.Logf("   当前区块: %d / 最高区块: %d",
			syncProgress.CurrentBlock, syncProgress.HighestBlock)
	} else {
		t.Log("✅ 节点已完全同步")
	}
}
