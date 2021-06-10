package mklic

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
)

func aesEncrypt(in []byte, key string) string {
	block, _ := aes.NewCipher([]byte(key))
	in = pKCS7Padding(in, block.BlockSize())
	encrypted := make([]byte, len(in))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(in); bs, be = bs+size, be+size {
		block.Encrypt(encrypted[bs:be], in[bs:be])
	}
	return base64.StdEncoding.EncodeToString(encrypted)
}

func aesDecrypt(in []byte, key string) ([]byte, error) {
	content, err := base64.StdEncoding.DecodeString(string(in))
	if err != nil {
		return nil, err
	}

	block, _ := aes.NewCipher([]byte(key))
	decrypted := make([]byte, len(content))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(content); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], content[bs:be])
	}

	return pKCS7UnPadding(decrypted), nil
}

func pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS7UnPadding(s []byte) []byte {
	length := len(s)
	padding := int(s[length-1])
	return s[:(length - padding)]
}
