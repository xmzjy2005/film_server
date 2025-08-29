package util

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
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

// GenerateSalt 生成 length为16的随机字符串
func GenerateSalt() (uuid string) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid = fmt.Sprintf("%X", b)
	return
}

// 解析私钥 传入私钥字符串
func ParsePriKeyBytes(buf []byte) (*rsa.PrivateKey, error) {
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	if p == nil {
		return nil, errors.New("private key parse error")
	}
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

// 解析公钥 将pem格式的公钥解析成rsa格式，供给jwt验证
func ParsePubKeyBytes(buf []byte) (*rsa.PublicKey, error) {
	p, _ := pem.Decode(buf)
	if p == nil {
		return nil, errors.New("parse publicKey Content nil")
	}
	//解析PKCS1格式的
	pubInterface, err := x509.ParsePKIXPublicKey(p.Bytes)
	if err != nil {
		return nil, errors.New("x509 parse pkcs1 格式错误")
	}
	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("断言错误，不是rsa.Publickey")
	}
	return pubKey, nil
}
