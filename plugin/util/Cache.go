package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const indexFIle = "cache/index.json"

// 索引
type actionIndex map[string]string

var (
	index actionIndex
	once  sync.Once
	mu    sync.Mutex
)

func loadIndex() {
	once.Do(func() {
		index = make(actionIndex)
		data, err := os.ReadFile(indexFIle)
		if err != nil {
			return
		}
		if err := json.Unmarshal(data, &index); err != nil {
			fmt.Println("未知错误，索引文件无法解析", err)
		}
	})
}

// 保存到索引文件
func saveIndex() error {
	mu.Lock()
	defer mu.Unlock()
	data, err := json.MarshalIndent(index, "", "	")
	if err != nil {
		return err
	}
	return os.WriteFile(indexFIle, data, 0644)
}

// 将url变成hash
func getSafeFilename(action string) string {
	hash := md5.Sum([]byte(action))
	hashStr := hex.EncodeToString(hash[:])
	//加载索引
	loadIndex()
	mu.Lock()
	defer mu.Unlock()
	eAction, e := index[hashStr]
	if !e || eAction != action {
		index[hashStr] = action
		go saveIndex()
	}
	return hashStr + ".bin"
}

// 从action.bin文件读取文件内的二进制内容，并且返回
func CacheGet(action string) ([]byte, error) {
	fileName := filepath.Join("cache", getSafeFilename(action))
	//读取文件
	return os.ReadFile(fileName)
}

// 将resp保存到action.bin文件内，如果有相同文件名，则覆盖
func CacheSave(action string, resp []byte) error {
	filename := filepath.Join("cache", getSafeFilename(action))
	return os.WriteFile(filename, resp, 0644)
}
