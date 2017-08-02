package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

// Aes is data encryptor
type Aes struct {
	block cipher.Block
}

// NewAes initialize new instance of Aes.
func NewAes(b cipher.Block) *Aes {
	return &Aes{b}
}

// Encrypt performs an encryption using AES/CBC/PKCS7
// with a random IV prepended using AES/ECB/None.
func (c *Aes) Encrypt(src []byte) ([]byte, error) {
	// get a random IV and ECB encrypt it
	iv := make([]byte, aes.BlockSize)
	_, err := rand.Read(iv)
	if err != nil {
		return nil, err
	}

	encryptedIv := make([]byte, aes.BlockSize)
	c.encryptBlocks(encryptedIv, iv)

	// pad it, copy the IV to the first 16 bytes and encrypt the rest with CBC
	encrypted := c.pad(src)
	copy(encrypted, encryptedIv)
	cipher.NewCBCEncrypter(c.block, iv).CryptBlocks(encrypted[aes.BlockSize:], encrypted[aes.BlockSize:])

	return encrypted, nil
}

// Decrypts data from the reader using AES/CBC/PKCS7 with an IV
// prepended using AES/ECB/None. The src slice may not be used anymore.
func (c *Aes) Decrypt(src []byte) []byte {
	iv := src[:aes.BlockSize]
	c.decryptBlocks(iv, iv)

	data := src[aes.BlockSize:]
	cipher.NewCBCDecrypter(c.block, iv).CryptBlocks(data, data)

	return c.unpad(data)
}


func (c *Aes) pad(src []byte) []byte {
	missing := aes.BlockSize - (len(src) % aes.BlockSize)
	newSize := len(src) + aes.BlockSize + missing
	dest := make([]byte, newSize, newSize)
	copy(dest[aes.BlockSize:], src)

	padding := byte(missing)
	for i := newSize - missing; i < newSize; i++ {
		dest[i] = padding
	}

	return dest
}

func (c *Aes) unpad(src []byte) []byte {
	padLen := src[len(src)-1]

	return src[:len(src)-int(padLen)]
}

func (c *Aes) encryptBlocks(dst, src []byte) {
	for len(src) > 0 {
		c.block.Encrypt(dst, src[:c.block.BlockSize()])
		src = src[c.block.BlockSize():]
		dst = dst[c.block.BlockSize():]
	}
}

func (c *Aes) decryptBlocks(dst, src []byte) {
	for len(src) > 0 {
		c.block.Decrypt(dst, src[:c.block.BlockSize()])
		src = src[c.block.BlockSize():]
		dst = dst[c.block.BlockSize():]
	}
}

