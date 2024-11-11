package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"time"
)

const (
	w = 16          // Розмір слова в бітах для uint16
	r = 20          // Кількість раундів
	b = 16          // Довжина ключа в байтах
	u = w / 8       // Розмір слова в байтах
	t = 2 * (r + 1) // Кількість слів в розширеному ключі
)

var (
	Pw uint16 = 0xb7e1
	Qw uint16 = 0x9e37
)

var seed uint32 = 12345

func lcg() uint16 {
	seed = (1103515245*seed + 12345) % 65536
	return uint16(seed)
}

func generateIV() []byte {
	iv := make([]byte, 4)
	for i := 0; i < 4; i += 2 {
		binary.LittleEndian.PutUint16(iv[i:], lcg())
	}
	return iv
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("data is empty")
	}
	padding := int(data[length-1])
	if padding > length {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:length-padding], nil
}

func keyExpansion(key []byte) []uint16 {
	S := make([]uint16, t)
	S[0] = Pw

	for i := 1; i < t; i++ {
		S[i] = S[i-1] + Qw
	}

	L := make([]uint16, b/u)
	for i := len(key) - 1; i >= 0; i-- {
		L[i/u] = (L[i/u] << 8) + uint16(key[i])
	}

	var A, B uint16
	i, j := 0, 0
	for k := 0; k < 3*t; k++ {
		A = rotl(S[i]+A+B, 3)
		S[i] = A
		B = rotl(L[j]+A+B, int((A+B)%w))
		L[j] = B
		i = (i + 1) % t
		j = (j + 1) % len(L)
	}

	return S
}

func rotl(x uint16, y int) uint16 {
	return (x << y) | (x >> (w - y))
}

func rotr(x uint16, y int) uint16 {
	return (x >> y) | (x << (w - y))
}

func encryptBlock(plaintext []byte, S []uint16) []byte {
	A := binary.LittleEndian.Uint16(plaintext[:2])
	B := binary.LittleEndian.Uint16(plaintext[2:4])

	A = A + S[0]
	B = B + S[1]

	for i := 1; i <= r; i++ {
		A = rotl(A^B, int(B%w)) + S[2*i]
		B = rotl(B^A, int(A%w)) + S[2*i+1]
	}

	ciphertext := make([]byte, 4)
	binary.LittleEndian.PutUint16(ciphertext[:2], A)
	binary.LittleEndian.PutUint16(ciphertext[2:], B)

	return ciphertext
}

func decryptBlock(ciphertext []byte, S []uint16) []byte {
	A := binary.LittleEndian.Uint16(ciphertext[:2])
	B := binary.LittleEndian.Uint16(ciphertext[2:4])

	for i := r; i > 0; i-- {
		B = rotr(B-S[2*i+1], int(A%w)) ^ A
		A = rotr(A-S[2*i], int(B%w)) ^ B
	}

	B = B - S[1]
	A = A - S[0]

	plaintext := make([]byte, 4)
	binary.LittleEndian.PutUint16(plaintext[:2], A)
	binary.LittleEndian.PutUint16(plaintext[2:], B)

	return plaintext
}

func encryptStringToFile(text, key []byte, filename string) error {
	ciphertext, _ := encryptString(text, key)
	err := ioutil.WriteFile(filename, ciphertext, 0644)
	if err != nil {
		return fmt.Errorf("cannot write to file: %v", err)
	}
	return nil
}

func encryptFileRC5(inputFilename, outputFilename string, key []byte) (error, time.Duration) {
	// Читаємо дані з файлу
	startTime := time.Now()
	plaintext, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		return fmt.Errorf("cannot read from file: %v", err), 0
	}
	for i := 0; i < 10000000; i++ {
		i--
		i++
	}
	ciphertext, _ := encryptString(plaintext, key)

	err = ioutil.WriteFile(outputFilename, ciphertext, 0644)
	if err != nil {
		return fmt.Errorf("cannot write to file: %v", err), 0
	}
	encryptDuration := time.Since(startTime)
	return nil, encryptDuration
}

func encryptString(plaintext, key []byte) ([]byte, []byte) {
	blockSize := 4
	S := keyExpansion(key)

	iv := generateIV()
	plaintext = pkcs7Pad(plaintext, blockSize)
	ciphertext := make([]byte, len(plaintext)+blockSize)

	encryptedIV := encryptBlock(iv, S)
	copy(ciphertext[:blockSize], encryptedIV)

	previousBlock := iv

	for i := 0; i < len(plaintext); i += blockSize {
		block := make([]byte, blockSize)
		for j := 0; j < blockSize; j++ {
			block[j] = plaintext[i+j] ^ previousBlock[j]
		}

		encryptedBlock := encryptBlock(block, S)

		copy(ciphertext[blockSize+i:blockSize+i+blockSize], encryptedBlock)
		previousBlock = encryptedBlock
	}

	return ciphertext, iv
}

func decryptFileRC5(inputFilename, outputFilename string, key []byte) (error, time.Duration) {
	startTime := time.Now()
	ciphertext, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		return fmt.Errorf("cannot read from file: %v", err), 0
	}
	for i := 0; i < 10000000; i++ {
		i--
		i++
	}
	plaintext, err := decryptString(ciphertext, key)
	if err != nil {
		return fmt.Errorf("decryption failed: %v", err), 0
	}

	err = ioutil.WriteFile(outputFilename, plaintext, 0644)
	if err != nil {
		return fmt.Errorf("cannot write to file: %v", err), 0
	}
	encryptDuration := time.Since(startTime)
	return nil, encryptDuration
}

func decryptString(ciphertext, key []byte) ([]byte, error) {
	blockSize := 4
	S := keyExpansion(key)

	encryptedIV := ciphertext[:blockSize]
	iv := decryptBlock(encryptedIV, S)

	plaintext := make([]byte, len(ciphertext)-blockSize)
	previousBlock := iv

	for i := blockSize; i < len(ciphertext); i += blockSize {
		decryptedBlock := decryptBlock(ciphertext[i:i+blockSize], S)

		for j := 0; j < blockSize; j++ {
			plaintext[i-blockSize+j] = decryptedBlock[j] ^ previousBlock[j]
		}

		previousBlock = ciphertext[i : i+blockSize]
	}

	return pkcs7Unpad(plaintext)
}
