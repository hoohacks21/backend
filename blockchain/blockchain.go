package github.com/hoohacks21/blockchain

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

func (b *Blockchain) RegisterTask(task Task) bool {
	task.UID = strings.ToLower(task.UID)
	task.TaskID = strings.ToLower(task.TaskID)
	b.PendingTasks = append(b.PendingTasks, task)
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (b *Blockchain) RegisterNode(node string) bool {
	if !contains(b.NetworkNodes, node) {
		b.NetworkNodes = append(b.NetworkNodes, node)
	}
	return true
}

func (b *Blockchain) CreateNewBlock(nonce int, previousBlockHash string, hash string) Block {
	newBlock := Block{
		Index:     len(b.Chain) + 1,
		Tasks:      b.PendingTasks,
		Timestamp: time.Now().UnixNano(),
		Nonce:     nonce,
		Hash:      hash, PreviousBlockHash: previousBlockHash}

	b.PendingTasks = Tasks{}
	b.Chain = append(b.Chain, newBlock)
	return newBlock
}

func (b *Blockchain) GetLastBlock() Block {
	return b.Chain[len(b.Chain)-1]
}

func (b *Blockchain) HashBlock(previousBlockHash string, currentBlockData string, nonce int) string {
	h := sha256.New()
	strToHash := previousBlockHash + currentBlockData + strconv.Itoa(nonce)
	h.Write([]byte(strToHash))
	hashed := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return hashed
}

func (b *Blockchain) ProofOfWork(previousBlockHash string, currentBlockData string) int {
	nonce := -1
	inputFmt := ""
	for inputFmt != "00000000" {
		nonce = nonce + 1
		hash := b.HashBlock(previousBlockHash, currentBlockData, nonce)
		inputFmt = hash[0:8]
	}
	return nonce
}

func (b *Blockchain) CheckNewBlockHash(newBlock Block) bool {
	lastBlock := b.GetLastBlock()
	correctHash := lastBlock.Hash == newBlock.PreviousBlockHash
	correctIndex := (lastBlock.Index + 1) == newBlock.Index

	return (correctHash && correctIndex)
}

func (b *Blockchain) ChainIsValid() bool {
	i := 1
	for i < len(b.Chain) {
		currentBlock := b.Chain[i]
		prevBlock := b.Chain[i-1]
		currentBlockData := BlockData{Index: strconv.Itoa(prevBlock.Index - 1), Tasks: currentBlock.Tasks}
		currentBlockDataAsByteArray, _ := json.Marshal(currentBlockData)
		currentBlockDataAsStr := base64.URLEncoding.EncodeToString(currentBlockDataAsByteArray)
		blockHash := b.HashBlock(prevBlock.Hash, currentBlockDataAsStr, currentBlock.Nonce)

		if blockHash[0:8] != "00000000" {
			return false
		}

		if currentBlock.PreviousBlockHash != prevBlock.Hash {
			return false
		}

		i = i + 1
	}

	genesisBlock := b.Chain[0]
	correctNonce := genesisBlock.Nonce == 100
	correctPreviousBlockHash := genesisBlock.PreviousBlockHash == "0"
	correctHash := genesisBlock.Hash == "0"
	correctTasks := len(genesisBlock.Tasks) == 0

	return (correctNonce && correctPreviousBlockHash && correctHash && correctTasks)
}


func (b *Blockchain) GetTasksForUser(UID string) Tasks {
	tasks := Tasks{}
	i := 0
	chainLength := len(b.Chain)
	for i < chainLength {
		block := b.Chain[i]
		tasksInBlock := block.Tasks
		j := 0
		tasksLength := len(tasksInBlock)
		for j < tasksLength {
			task := tasksInBlock[j]
			if task.UID == UID {
				tasks = append(matchTasks, task)
			}
			j = j + 1
		}
		i = i + 1
	}
	return tasks
}