package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func main() {
	// plaintext values to encrypt
	plaintextValues := map[string]string{
		"ISE_USER":     "root",
		"ISE_PASSWORD": "M4rb73HalLs",
		"ISE_SERVER":   "19.14.250.23",
		"SERVER_PORT":  ":8080",
	}

	// 256-bit key (32 bytes)
	key := []byte("01234567890123456789012345678901")

	// create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating cipher block:", err)
		return
	}

	// create a new GCM cipher
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Error creating GCM cipher:", err)
		return
	}

	// generate a random nonce (12 bytes)
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("Error generating nonce:", err)
		return
	}

	// create a new .env file
	envFile := ".env"
	f, err := os.Create(envFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	// write the encrypted values to the .env file
	for k, v := range plaintextValues {
		encrypted := aesgcm.Seal(nil, nonce, []byte(v), nil)
		encoded := base64.StdEncoding.EncodeToString(append(nonce, encrypted...))
		_, err := f.WriteString(fmt.Sprintf("%s=%s\n", k, encoded))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	fmt.Println("Encrypted values written to .env file:", envFile)
}
