package system

import (
	"errors"
	"film_server/config"
	"film_server/plugin/db"
	"film_server/plugin/util"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type UserClaims struct {
	UserID   uint   `json:"userID"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

// 生成token
func GetToken(userId uint, userName string) (string, error) {
	uc := UserClaims{
		UserID:   userId,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Issuer,
			Subject:   userName,
			Audience:  jwt.ClaimStrings{"Auth_All"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AuthTokenExpires * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-10 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        util.GenerateSalt(),
		},
	}
	priKey, err := util.ParsePriKeyBytes([]byte(config.PrivateKey))
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, uc).SignedString(priKey)
	return token, err
}

// 保存token到redis中
func SaveUserToken(token string, userId uint) error {
	return db.Rdb.Set(db.Cxt, fmt.Sprintf(config.UserTokenKey, userId), token, (config.AuthTokenExpires*7*24)*time.Hour).Err()
}

// 解析token
func ParseToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		//这边转换了公钥变成jwt rsa格式的公钥用于解析token
		pub, err := util.ParsePubKeyBytes([]byte(config.PublicKey))
		if err != nil {
			return nil, err
		}
		return pub, nil
	})
	if err != nil {
		//解析有错，且类型为token过期
		if errors.Is(err, jwt.ErrTokenExpired) {
			claims, _ := token.Claims.(*UserClaims)
			return claims, err
		}
	}
	//验证token是否有效
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}
	//解析token的claims的内容,断言为userClaims
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid userClaim type")
	}
	return claims, err
}

// 从redis中获取指定user_id的token
func GetUserTokenById(userId uint) string {
	token, err := db.Rdb.Get(db.Cxt, fmt.Sprintf(config.UserTokenKey, userId)).Result()
	if err != nil {
		log.Println(err)
		return ""
	}
	return token
}
