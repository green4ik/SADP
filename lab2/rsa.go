package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"time"
)

// func main() {
// 	// Генерація RSA ключів
// 	privateKey, err := rsa.GenerateKey(rand.Reader, 8192)
// 	if err != nil {
// 		fmt.Println("Помилка при генерації ключа:", err)
// 		return
// 	}
// 	fmt.Println("Ключі сгенеровані", err)
// 	publicKey := &privateKey.PublicKey

// 	// Шифрування файлу
// 	startTime := time.Now() // Початок вимірювання часу
// 	err = encryptFile("test.txt", "encrypted", publicKey)
// 	if err != nil {
// 		fmt.Println("Помилка при шифруванні файлу:", err)
// 		return
// 	}
// 	encryptDuration := time.Since(startTime) // Час шифрування
// 	fmt.Printf("Файл успішно зашифровано за %v.\n", encryptDuration)

// 	// Розшифрування файлу
// 	startTime = time.Now() // Початок вимірювання часу
// 	err = decryptFile("encrypted", "decrypted", privateKey)
// 	if err != nil {
// 		fmt.Println("Помилка при розшифруванні файлу:", err)
// 		return
// 	}
// 	decryptDuration := time.Since(startTime) // Час розшифрування
// 	fmt.Printf("Файл успішно розшифровано за %v.\n", decryptDuration)
// }

// Функція для шифрування файлу
func encryptFileRSA(inputPath, outputPath string, publicKey *rsa.PublicKey) (error, time.Duration) {
	// Читання даних з файлу
	startTime := time.Now()
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("не вдалося прочитати файл: %w", err), 0
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
