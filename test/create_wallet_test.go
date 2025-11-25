package test

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

func TestCreateWallet(t *testing.T) {
	// 生成ECDSA私钥 - 椭圆曲线数字签名算法，用于以太坊账户体系
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	// 将私钥转换为字节格式并编码为16进制字符串
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 去掉'0x'前缀，显示原始私钥

	// 从私钥推导出对应的公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法断言类型: 公钥不是*ecdsa.PublicKey类型")
	}

	// 将公钥转换为字节格式并编码为16进制字符串
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("来自公钥:", hexutil.Encode(publicKeyBytes)[4:]) // 去掉'0x04'前缀，显示原始公钥

	// 通过公钥生成以太坊地址 - 这是标准的以太坊地址生成流程
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("以太坊地址:", address)

	// 手动计算地址: 对公钥进行Keccak-256哈希计算，然后取后20字节
	hash := sha3.NewLegacyKeccak256() // Keccak-256是以太坊使用的哈希算法
	hash.Write(publicKeyBytes[1:])    // 跳过ECDSAPub编码的第一个字节(0x04)
	fmt.Println("完整哈希:", hexutil.Encode(hash.Sum(nil)[:]))
	fmt.Println("地址哈希:", hexutil.Encode(hash.Sum(nil)[12:])) // 原长32位，截去前12位，保留后20位作为地址
	fmt.Println("校验和地址:", common.BytesToAddress(hash.Sum(nil)[12:]).Hex())

	fmt.Println("以太坊地址 == 校验和地址:", address == common.BytesToAddress(hash.Sum(nil)[12:]).Hex())

}
