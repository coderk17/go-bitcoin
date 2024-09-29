package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Block 表示区块链中的一个区块
type Block struct {
	Index     int
	Timestamp string
	Data      string
	PrevHash  string
	Hash      string
	Nonce     int
}

// Blockchain 是一个区块的切片，代表整个区块链
type Blockchain []Block

// 计算区块的哈希值
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.Data + block.PrevHash + strconv.Itoa(block.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 创建新的区块
func generateBlock(oldBlock Block, data string) Block {
	var newBlock Block

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Nonce = 0

	// 简化的工作量证明
	for i := 0; ; i++ {
		newBlock.Nonce = i
		newBlock.Hash = calculateHash(newBlock)
		if strings.HasPrefix(newBlock.Hash, "00") {
			break
		}
	}

	return newBlock
}

// 验证区块是否有效
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

// 保存区块链到文件
func saveBlockchain(blockchain Blockchain) error {
	data, err := json.MarshalIndent(blockchain, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("blockchain.json", data, 0644)
}

// 从文件加载区块链
func loadBlockchain() (Blockchain, error) {
	var blockchain Blockchain

	if _, err := os.Stat("blockchain.json"); os.IsNotExist(err) {
		// 如果文件不存在，返回空的区块链
		return blockchain, nil
	}

	data, err := ioutil.ReadFile("blockchain.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &blockchain)
	if err != nil {
		return nil, err
	}

	return blockchain, nil
}

func main() {
	// 尝试加载现有的区块链
	blockchain, err := loadBlockchain()
	if err != nil {
		log.Fatal("加载区块链失败:", err)
	}

	if len(blockchain) == 0 {
		// 如果区块链为空，创建创世区块
		genesisBlock := Block{0, time.Now().String(), "创世区块", "", "", 0}
		genesisBlock.Hash = calculateHash(genesisBlock)
		blockchain = append(blockchain, genesisBlock)
	}

	// 添加新区块
	blockchain = append(blockchain, generateBlock(blockchain[len(blockchain)-1], "第二个区块"))
	blockchain = append(blockchain, generateBlock(blockchain[len(blockchain)-1], "第三个区块"))

	// 保存区块链到文件
	err = saveBlockchain(blockchain)
	if err != nil {
		log.Fatal("保存区块链失败:", err)
	}

	// 打印区块链
	for _, block := range blockchain {
		fmt.Printf("索引: %d\n", block.Index)
		fmt.Printf("时间戳: %s\n", block.Timestamp)
		fmt.Printf("数据: %s\n", block.Data)
		fmt.Printf("哈希: %s\n", block.Hash)
		fmt.Printf("前一个哈希: %s\n", block.PrevHash)
		fmt.Printf("随机数: %d\n", block.Nonce)
		fmt.Println("------------------------")
	}
}