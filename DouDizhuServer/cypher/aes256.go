package cypher

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// 返回值为(加密后的数据，iv，错误)
func AesEncryptCBC(plaintext, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("aes.NewCipher failed: %w", err)
	}

	// 生成随机 IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, fmt.Errorf("iv generation failed: %w", err)
	}

	// 应用 PKCS7 填充
	paddedPlaintext := PKCS7Pad(plaintext, aes.BlockSize)

	// 确保数据长度是块大小的倍数
	if len(paddedPlaintext)%aes.BlockSize != 0 {
		return nil, nil, errors.New("padded data length not multiple of block size")
	}

	ciphertext := make([]byte, len(paddedPlaintext))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	return ciphertext, iv, nil
}

func AesDecryptCBC(ciphertext, iv, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher failed: %w", err)
	}

	// 检查 IV 长度
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid IV length: %d (expected %d)", len(iv), aes.BlockSize)
	}

	// 确保加密数据长度是块大小的倍数
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext length not multiple of block size")
	}

	// 创建解密器
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密数据（原地修改）
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	// 去除 PKCS7 填充
	unpadded, err := PKCS7Unpad(decrypted)
	if err != nil {
		return nil, fmt.Errorf("PKCS7Unpad failed: %w", err)
	}

	return unpadded, nil
}

func PKCS7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func PKCS7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("pkcs7: data is empty")
	}

	// 获取填充长度
	padding := int(data[len(data)-1])

	// 验证填充长度有效性
	if padding < 1 || padding > aes.BlockSize {
		return nil, errors.New("pkcs7: invalid padding size")
	}

	// 检查数据长度是否足够
	if len(data) < padding {
		return nil, errors.New("pkcs7: data length less than padding")
	}

	// 验证所有填充字节是否一致
	for i := len(data) - padding; i < len(data); i++ {
		if int(data[i]) != padding {
			return nil, errors.New("pkcs7: invalid padding bytes")
		}
	}

	return data[:len(data)-padding], nil
}
