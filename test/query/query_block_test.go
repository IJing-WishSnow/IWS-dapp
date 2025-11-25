package query

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestQueryBlock(t *testing.T) {

	header, err := client.HeaderByNumber(context.Background(), blockNumber) // 根据区块号获取区块头信息
	fmt.Println(header.Number.Uint64())                                     // 打印区块号：23866957
	fmt.Println(header.Time)                                                // 打印区块时间戳
	fmt.Println(header.Difficulty.Uint64())                                 // 打印挖矿难度值
	fmt.Println(header.Hash().Hex())                                        // 打印区块哈希值

	if err != nil {
		log.Fatal(err)
	}

	block, err := client.BlockByNumber(context.Background(), blockNumber) // 根据区块号获取完整区块信息
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(block.Number().Uint64())     // 打印区块号：23866957
	fmt.Println(block.Time())                // 打印区块时间戳：1763965919
	fmt.Println(block.Difficulty().Uint64()) // 打印挖矿难度值：0
	fmt.Println(block.Hash().Hex())          // 打印区块哈希值：0xbd829c3842e414aa7cab0040b5ff50a3c7a6ffacbbb653d1288fbbb9d8ba92d7
	fmt.Println(len(block.Transactions()))   // 打印交易数量：123

	count, err := client.TransactionCount(context.Background(), block.Hash()) // 通过区块哈希查询该区块中的交易总数

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(count) // 打印交易总数：123
}

// ```go
// fmt.Println(len(block.Transactions()))   // 从已获取的区块对象中直接计算交易数量：123
// fmt.Println(count) // 通过额外API调用查询该区块的交易总数：123
// ```

// **区别：**
// - `len(block.Transactions())`：从内存中的区块数据直接读取交易列表长度
// - `TransactionCount()`：向节点发起新的网络请求查询交易数量

// **结果相同，但来源不同** - 一个是从本地数据计算，一个是远程API查询。
