package functions

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.assistagro.com/back/back.auth.go/internal/apierrors"

	"golang.org/x/crypto/pbkdf2"
)

func CreateRandomString(len int) string {
	if len <= 0 {
		len = 80
	}
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFJHIJKLMNOPQRSTUVWXXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func PasswordEncode(password string, salt string, iterations int) (string, error) {
	if strings.TrimSpace(password) == "" {
		return "", errors.New("no password string provided")
	}
	if strings.TrimSpace(salt) == "" {
		salt = CreateRandomString(16)
	}

	if strings.Contains(salt, "$") {
		return "", errors.New("salt contains dollar sign ($)")
	}
	if iterations <= 0 {
		iterations = 260000
	}

	hash := pbkdf2.Key([]byte(password), []byte(salt), iterations, sha256.Size, sha256.New)

	hexHash := hex.EncodeToString(hash)
	return fmt.Sprintf("pbkdf2:sha256:%d$%s$%s", iterations, salt, hexHash), nil
}

func NewInternalServerError(c *gin.Context, message string) {
	log.Error(message)
	c.String(http.StatusInternalServerError, apierrors.ErrInternalServerError)
}
