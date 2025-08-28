package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// 密码加密算法，根据原生密码，算出加密后存储的密码
func PasswordEncrypt(password string, salt string) string {
	b := []byte(fmt.Sprint(password, salt))
	var r [16]byte
	for i := 0; i < 3; i++ {
		r = md5.Sum(b)
		b = []byte(hex.EncodeToString(r[:]))
	}
	return hex.EncodeToString(r[:])
}
