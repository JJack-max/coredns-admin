package service

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"time"
)

func MakeSecret(password string) string {
	/* Generate a new, cryptographically secure salt */
	salt := genSalt(6)

	sha := sha1.New()
	_, _ = sha.Write([]byte(password))
	_, _ = sha.Write(salt)
	return base64.StdEncoding.EncodeToString(sha.Sum(salt))
}

func genSalt(n int) []byte {
	b := make([]byte, n)
	if _n, err := rand.Read(b); err != nil || _n != n {
		nano := time.Now().UnixNano()
		for i := 0; i < n; i++ {
			b[i] = byte(nano >> uint(i&0b1111111))
		}
	}
	return b
}

func VerifySecret(secret, password string) bool {

	bs, err := base64.StdEncoding.DecodeString(secret)
	if err != nil || len(bs) < 20 {
		return false
	}
	salt := make([]byte, len(bs)-20)
	copy(salt, bs[:len(bs)-20])
	testHash := sha1.New()
	_, _ = testHash.Write([]byte(password))
	_, _ = testHash.Write(salt)
	return bytes.Equal(bs, testHash.Sum(salt))

}
