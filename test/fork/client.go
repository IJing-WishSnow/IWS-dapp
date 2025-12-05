package fork

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockchainClient struct {
	*ethclient.Client
	URL     string
	ChainID *big.Int
	ctx     context.Context
	cancel  context.CancelFunc
}

// 创建新客户端
func NewClient(rpcURL string, timeout time.Duration) (*BlockchainClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("连接失败: %v", err)
	}

	// 测试连接
	_, err = client.NetworkID(ctx)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("网络连接测试失败: %v", err)
	}

	return &BlockchainClient{
		Client: client,
		URL:    rpcURL,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// 释放资源
func (c *BlockchainClient) Close() {
	if c.Client != nil {
		c.Client.Close()
	}
	if c.cancel != nil {
		c.cancel()
	}
}

// 获取当前区块号
func (c *BlockchainClient) GetCurrentBlockNumber() (uint64, error) {
	return c.BlockNumber(c.ctx)
}

// 获取网络ID
func (c *BlockchainClient) GetNetworkID() (*big.Int, error) {
	if c.ChainID != nil {
		return c.ChainID, nil
	}

	chainID, err := c.NetworkID(c.ctx)
	if err != nil {
		return nil, err
	}

	c.ChainID = chainID
	return chainID, nil
}

// 获取账户余额
func (c *BlockchainClient) GetBalance(address string) (*big.Int, error) {
	addr := common.HexToAddress(address)
	return c.BalanceAt(c.ctx, addr, nil)
}

// 检查合约代码
func (c *BlockchainClient) HasContractCode(address string) (bool, int, error) {
	addr := common.HexToAddress(address)
	code, err := c.CodeAt(c.ctx, addr, nil)
	if err != nil {
		return false, 0, err
	}
	return len(code) > 0, len(code), nil
}

// Wei 转 Ether
func WeiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
}
