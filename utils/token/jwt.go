package token

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"neptune/global"
	myerrors "neptune/utils/errors"
	"neptune/utils/rsp"
	"time"
)

type CustomClaims struct {
	UserID   uint
	UserRole string
	jwt.RegisteredClaims
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.GetHeader("Authorization")

		if token == "" {
			rsp.ErrRsp(c, myerrors.TokenInvalidErr{Err: fmt.Errorf("无效Token")})
			c.Abort()
			return
		}
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, ErrTokenExpired) {
				rsp.ErrRsp(c, myerrors.TokenInvalidErr{Err: fmt.Errorf("token已过期")})
				c.Abort()
				return
			}
			rsp.ErrRsp(c, myerrors.TokenInvalidErr{Err: fmt.Errorf("无效Token")})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

func GetClaims(c *gin.Context) *CustomClaims {
	claims, exist := c.Get("claims")
	if !exist {
		return nil
	}
	return claims.(*CustomClaims)
}

type JWT struct {
	SigningKey []byte
}

var (
	ErrTokenInvalid = errors.New("couldn't handle this token")
	ErrTokenExpired = errors.New("token is expired")
)

// NewJWT 创建一个新的jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTConfig.SigningKey),
	}
}

// CreateToken 创建一个token并使用key进行签名
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析token并返回自定义的Claims
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		log.Error("parse token err:", err)
		return nil, err
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, ErrTokenInvalid

	} else {
		return nil, ErrTokenInvalid

	}

}

// RefreshToken 更新token
//func (j *JWT) RefreshToken(tokenString string) (string, error) {
//	jwt.TimeFunc = func() time.Time {
//		return time.Unix(0, 0)
//	}
//	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return j.SigningKey, nil
//	})
//	if err != nil {
//		return "", err
//	}
//	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
//		jwt.TimeFunc = time.Now
//		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
//		return j.CreateToken(*claims)
//	}
//	return "", TokenInvalid
//}

func GenerateToken(Id int, role string) (string, error) {
	//生成token信息
	j := NewJWT()
	exp := time.Now().Add(time.Hour * time.Duration(global.ServerConfig.JWTConfig.ExpireTime)) // ExpireTime以小时为单位
	claims := CustomClaims{
		UserID:   uint(Id),
		UserRole: role,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	//生成token
	token, err := j.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}
