package util

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
)

// 保存图片文件到目录
func SaveOnlineFile(url, dir string) (path string, err error) {
	//请求结构体
	r := &RequestInfo{Uri: url}
	ApiGet(r)

	if len(r.Resp) <= 0 {
		err = errors.New("寻找图片连接地址失败：" + url)
		return
	}

	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		//不存在则创建目录
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return
		}
	}

	//拼接文件并且存储
	fileName := filepath.Join(dir, filepath.Base(url))
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, err = writer.Write(r.Resp)
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
	return filepath.Base(fileName), nil

}
