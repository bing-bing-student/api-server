package utils

import (
	"blog/global"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type UserClaims struct {
	UserID uint64
	jwt.RegisteredClaims
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func GenerateToken(userClaims *UserClaims) (*TokenDetails, error) {
	tokenDetails := &TokenDetails{}
	mySigningKey := []byte(global.CONFIG.JWTConfig.SigningKey)

	// 设置访问和刷新令牌的过期时间
	// 3小时后过期
	tokenDetails.AtExpires = time.Now().Add(time.Hour * 1).Unix()
	tokenDetails.AccessUuid = uuid.New().String()

	// 7天后过期
	tokenDetails.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokenDetails.RefreshUuid = uuid.New().String()

	// 创建访问令牌的声明
	accessTokenClaims := &UserClaims{
		UserID: userClaims.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(tokenDetails.AtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建访问令牌
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(mySigningKey)
	if err != nil {
		return nil, err
	}
	tokenDetails.AccessToken = accessTokenString

	// 创建刷新令牌的声明
	refreshTokenClaims := &UserClaims{
		UserID: userClaims.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(tokenDetails.RtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建刷新令牌
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(mySigningKey)
	if err != nil {
		return nil, err
	}
	tokenDetails.RefreshToken = refreshTokenString

	return tokenDetails, nil
}

// ParseAccessToken 解析并验证访问令牌
func ParseAccessToken(tokenString string) (*UserClaims, error) {
	// 获取全局签名
	mySigningKey := []byte(global.CONFIG.JWTConfig.SigningKey)
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, errors.New("accessToken is null")
	}

	if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("accessToken is invalid")
	}
}

// ParseRefreshToken 解析并验证刷新令牌
func ParseRefreshToken(tokenString string) (*UserClaims, error) {
	mySigningKey := []byte(global.CONFIG.JWTConfig.SigningKey)
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, errors.New("refreshToken is null")
	}

	if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("refreshToken is invalid")
	}
}
