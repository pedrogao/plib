package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"

	"golang.org/x/crypto/scrypt"
)

// PasswordEncrypt encrypt password
func PasswordEncrypt(salt, password string) string {
	dk, _ := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	return fmt.Sprintf("%x", string(dk))
}

/*
	reference: https://eli.thegreenplace.net/2019/aes-encryption-of-files-in-go/
*/
func main() {
	// Dummy text to encrypt. Size should be multiple of AES block size (16).
	text := bytes.Repeat([]byte("i"), 96)

	key := sha256.Sum256([]byte("kitty"))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}

	// iv
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		log.Fatal(err)
	}

	// Create a new CBC mode encrypter using our AES block cipher, and use it
	// to encrypt our text.
	ciphertext := make([]byte, len(text))
	enc := cipher.NewCBCEncrypter(block, iv) // 这里的iv最好使用key
	enc.CryptBlocks(ciphertext, text)

	fmt.Println(ciphertext)
}
