package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth gin middleware，check jwt token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "Token is invalid",
			})
			c.Abort()
			return
		}

		j := NewJWT()
		// parseToken
		claims, err := j.ParseToken(token)
		if err != nil {

			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
	}
}

// JWT  the jwt main struct and includ sign key
type JWT struct {
	SigningKey []byte
}

var (
	//ErrTokenExpired Token is expired
	ErrTokenExpired error = errors.New("Token is expired")
	//ErrTokenNotValidYet Token not Valid yet
	ErrTokenNotValidYet error = errors.New("Token not Valid yet")
	//ErrTokenMalformed Token format is incorrect
	ErrTokenMalformed error = errors.New("Token format is incorrect")
	//ErrTokenInvalid   Token is invalid
	ErrTokenInvalid error = errors.New("Token is invalid")
	//SignKey secret string
	SignKey string = "newtrekWang"
)

//CustomClaims jwt playload
type CustomClaims struct {
	ID       string `json:"userId"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	jwt.RegisteredClaims
}

//NewJWT  New a JWT instance
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// GetSignKey get sign key for  anothe pkg
func GetSignKey() string {
	return SignKey
}

// SetSignKey set sign key for anothe pkg
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

// CreateToken  Create Token with claims
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	if claims.ExpiresAt == nil {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//ParseToken Parse Token from tokenString
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrTokenMalformed
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, ErrTokenNotValidYet
		default:
			return nil, ErrTokenInvalid
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

// RefreshToken Refresh Token by a jwt token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	parser := jwt.NewParser(jwt.WithTimeFunc(func() time.Time {
		return time.Unix(0, 0)
	}))
	token, err := parser.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
		return j.CreateToken(*claims)
	}
	return "", ErrTokenInvalid
}
