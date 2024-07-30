package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"io"
	"net/http"
	"os"

	"github.com/Sadere/ya-metrics/internal/common"
	"github.com/gin-gonic/gin"
)

// инициализирует middleware для расшифровки тела запроса, возвращает ошибку в случае неудачи
func Decrypt(privKeyFilePath string) (gin.HandlerFunc, error) {
	// читаем файл приватного ключа
	privKeyPEM, err := os.ReadFile(privKeyFilePath)
	if err != nil {
		return nil, err
	}

	// парсим ключ из содержимого файла
	privKeyBlock, _ := pem.Decode(privKeyPEM)
	privKey, err := x509.ParsePKCS1PrivateKey(privKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// сам middleware
	return func(c *gin.Context) {
		// проверяем AES ключ из заголовка
		encAESKey := c.Request.Header.Get(common.AESKeyHeader)

		if len(encAESKey) == 0 {
			c.Next()
			return
		}

		// расшифровываем AES ключ
		keyBytes, err := hex.DecodeString(encAESKey)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		AESKey, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, keyBytes)
		if err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// создаем AES блок
		AESBlock, err := aes.NewCipher(AESKey)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		GCM, err := cipher.NewGCM(AESBlock)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// читаем тело запроса
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// вектор инициализации
		nonce := body[:GCM.NonceSize()]
		cipherText := body[GCM.NonceSize():]

		decrypted, err := GCM.Open(nil, nonce, cipherText, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// заменяем тело на расшифрованное
		c.Request.Body = io.NopCloser(bytes.NewReader(decrypted))
		c.Next()
	}, nil
}
