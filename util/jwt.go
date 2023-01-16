package util

import (
	"errors"
	"github.com/CCpro10/micro_douyin/conf"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	//token秘钥
	mySecret []byte

	//token过期时间
	tokenExpireDuration time.Duration
)

func InitJWTVal() {
	tokenExpireDuration = time.Duration(int64(conf.Config.Jwt.TokenExpireDuration) * int64(time.Hour))
	mySecret = []byte(conf.Config.Jwt.Secret)
}

// MyClaims自定义声明结构体并内嵌jwt.StandardClaims
type MyClaims struct {
	UserId int64 //表示用户业务ID
	jwt.StandardClaims
}

//生成token,传入ID,userId,生成JWTString和err
func GenerateToken(userId int64) (string, error) {
	// 创建一个自己的声明/请求
	c := MyClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(), // 过期时间
			Issuer:    "github.com/CCpro10/micro_douyin",          // 签发人
			Subject:   "user token",
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的秘钥签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// 解析token，返回一个包含信息的用户声明
func parseToken(token string) (*MyClaims, error) {
	// 通过(tokenStruct,请求结构,返回秘钥的一个回调函数)这三个参数,返回一个token结构体
	tokenStruct, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}

	// 校验token,token有效则返回myClaims请求
	if myClaims, ok := tokenStruct.Claims.(*MyClaims); ok && tokenStruct.Valid {
		return myClaims, nil
	}
	//token无效，返回错误
	return nil, errors.New("invalid tokenStruct")
}

//解析query中的token,token有效会返回解析出的userId
//如果token为空,返回ErrNoAuth
//如果token无效或过期,返回ErrWrongAuth
func ParseToken(token string) (userId int64, err error) {
	if token == "" {
		return 0, ErrNoAuth
	}

	claim, err := parseToken(token)
	if err != nil {
		return 0, ErrWrongAuth
	}

	return claim.UserId, nil
}
