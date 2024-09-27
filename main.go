package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

func main() {
	// 创建创世区块
	genesisBlock := Block{0, time.Now().String(), "Genesis Block", "", "", 0}
	genesisBlock.Hash = calculateHash(genesisBlock)

	blockchain := make(Blockchain, 0)
	blockchain = append(blockchain, genesisBlock)

	// 添加新区块
	blockchain = append(blockchain, generateBlock(blockchain[len(blockchain)-1], "Second Block"))
	blockchain = append(blockchain, generateBlock(blockchain[len(blockchain)-1], "Third Block"))

	// 打印区块链
	for _, block := range blockchain {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", block.Timestamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("PrevHash: %s\n", block.PrevHash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Println("------------------------")
	}
}