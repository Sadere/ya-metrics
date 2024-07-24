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
)

// middleware для шифрования тела запроса
type CryptoRoundTripper struct {
	KeyFilePath string
	Next        http.RoundTripper
}

func (rt *CryptoRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	// читаем файл публичного ключа
	pubKeyPEM, err := os.ReadFile(rt.KeyFilePath)
	if err != nil {
		return nil, err
	}

	// парсим данные из формата PEM
	pubKeyBlock, _ := pem.Decode(pubKeyPEM)
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// создаем ключ AES-256
	key, err := common.GenerateRandom(2 * aes.BlockSize)
	if err != nil {
		return nil, err
	}

	AESBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// создаем блочный шифр
	GCM, err := cipher.NewGCM(AESBlock)
	if err != nil {
		return nil, err
	}

	// создаём вектор инициализации
	nonce, err := common.GenerateRandom(GCM.NonceSize())
	if err != nil {
		return nil, err
	}

	// зашифровываем тело запроса
	encrypted := GCM.Seal(nonce, nonce, body, nil)

	// зашифровываем AES ключ с помощью RSA
	encryptedAESKey, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), key)
	if err != nil {
		return nil, err
	}

	// ставим заголовок с ключом AES
	r.Header.Set(common.AESKeyHeader, hex.EncodeToString(encryptedAESKey))

	// заменияем тело зашифрованным текстом
	r.Body = io.NopCloser(bytes.NewReader(encrypted))

	return rt.Next.RoundTrip(r)
}
