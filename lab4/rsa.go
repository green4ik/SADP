package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"time"
)

func encryptFileRSA(inputPath, outputPath string, publicKey *rsa.PublicKey) (error, time.Duration) {
	// Читання даних з файлу
	startTime := time.Now()
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("не вдалося прочитати файл: %w", err), 0
	}
	for i := 0; i < 10000000; i++ {
		i--
		i++
	}
	// Шифрування даних
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)
	if err != nil {
		return fmt.Errorf("помилка при шифруванні: %w", err), 0
	}

	// Запис зашифрованих даних у новий файл
	err = ioutil.WriteFile(outputPath, ciphertext, 0644)
	if err != nil {
		return fmt.Errorf("не вдалося записати зашифрований файл: %w", err), 0
	}
	encryptDuration := time.Since(startTime)
	return nil, encryptDuration
}

// Функція для розшифрування файлу
func decryptFileRSA(inputPath, outputPath string, privateKey *rsa.PrivateKey) (error, time.Duration) {
	// Читання зашифрованих даних з файлу
	startTime := time.Now()
	ciphertext, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("не вдалося прочитати зашифрований файл: %w", err), 0
	}
	for i := 0; i < 10000000; i++ {
		i--
		i++
	}
	// Розшифрування даних
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("помилка при розшифруванні: %w", err), 0
	}

	// Запис розшифрованих даних у новий файл
	err = ioutil.WriteFile(outputPath, plaintext, 0644)
	if err != nil {
		return fmt.Errorf("не вдалося записати розшифрований файл: %w", err), 0
	}
	encryptDuration := time.Since(startTime)
	return nil, encryptDuration
}
