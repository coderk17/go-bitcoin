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
	"time"
	"math/big"
)

// Block 表示区块链中的一个区块
type Block struct {
	Index     int
	Timestamp string
	Data      string
	PrevHash  string
	Hash      string
	Nonce     int
	Difficulty int
}

// Blockchain 是一个区块的切片，代表整个区块链
type Blockchain []Block

// 计算区块的哈希值
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.Data + block.PrevHash + strconv.Itoa(block.Nonce) + strconv.Itoa(block.Difficulty)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 创建新的区块
func generateBlock(oldBlock Block, data string) Block {
	var newBlock Block

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = time.Now().Format(time.RFC3339)
	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = adjustDifficulty(oldBlock, newBlock)
	newBlock.Nonce = 0

	target := calculateTarget(newBlock.Difficulty)

	startTime := time.Now()
	for {
		newBlock.Hash = calculateHash(newBlock)
		hashInt, _ := new(big.Int).SetString(newBlock.Hash, 16)
		if hashInt.Cmp(target) == -1 {
			break
		}
		newBlock.Nonce++
	}
	endTime := time.Now()

	fmt.Printf("挖矿耗时: %v\n", endTime.Sub(startTime))

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
	target := calculateTarget(newBlock.Difficulty)
	hashInt, _ := new(big.Int).SetString(newBlock.Hash, 16)
	if hashInt.Cmp(target) != -1 {
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

// 添加新的函数来计算目标哈希
func calculateTarget(bits int) *big.Int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-bits))
	return target
}

// 添加新的函数来调整难度
func adjustDifficulty(oldBlock, newBlock Block) int {
	if oldBlock.Index == 0 {
		return oldBlock.Difficulty // 对于创世区块，保持初始难度
	}

	expectedTime := targetBlockTime
	actualTime, _ := time.Parse(time.RFC3339, newBlock.Timestamp)
	oldTime, _ := time.Parse(time.RFC3339, oldBlock.Timestamp)
	timeDiff := actualTime.Sub(oldTime)

	if timeDiff < expectedTime/2 {
		return oldBlock.Difficulty + 1
	} else if timeDiff > expectedTime*2 {
		return oldBlock.Difficulty - 1
	}
	return oldBlock.Difficulty
}

const targetBlockTime = 3 * time.Second // 目标出块时间为3秒

func main() {
	// 尝试加载现有的区块链
	blockchain, err := loadBlockchain()
	if err != nil {
		log.Fatal("加载区块链失败:", err)
	}

	if len(blockchain) == 0 {
		// 如果区块链为空，创建创世区块
		initialDifficulty := 20 // 设置初始难度，可以根据需要调整
		genesisBlock := Block{0, time.Now().Format(time.RFC3339), "创世区块", "", "", 0, initialDifficulty}
		genesisBlock.Hash = calculateHash(genesisBlock)
		blockchain = append(blockchain, genesisBlock)
	}

	// 添加新区块
	for i := 0; i < 20; i++ {
		startTime := time.Now()
		newBlock := generateBlock(blockchain[len(blockchain)-1], fmt.Sprintf("区块 #%d", len(blockchain)+1))
		blockchain = append(blockchain, newBlock)
		endTime := time.Now()
		miningTime := endTime.Sub(startTime)
		fmt.Printf("添加了新区块 #%d，难度：%d，挖矿时间：%v\n", newBlock.Index, newBlock.Difficulty, miningTime)
	}

	// 保存区块链到文件
	err = saveBlockchain(blockchain)
	if err != nil {
		log.Fatal("保存区块链失败:", err)
	}
}