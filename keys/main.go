package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	// Генерація RSA ключів
	privateKey, err := rsa.GenerateKey(rand.Reader, 16384)
	if err != nil {
		fmt.Println("Помилка при генерації ключа:", err)
		return
	}
	fmt.Println("Ключі сгенеровані", err)
	publicKey := &privateKey.PublicKey

	// Запис приватного ключа в файл
	err = writePrivateKeyToFile("private_key.pem", privateKey)
	if err != nil {
		fmt.Println("Помилка при запису приватного ключа:", err)
		return
	}
	fmt.Println("Приватний ключ успішно записаний у файл 'private_key.pem'.")

	// Запис публічного ключа в файл
	err = writePublicKeyToFile("public_key.pem", publicKey)
	if err != nil {
		fmt.Println("Помилка при запису публічного ключа:", err)
		return
	}
	fmt.Println("Публічний ключ успішно записаний у файл 'public_key.pem'.")
}

// Функція для запису приватного ключа у файл
func writePrivateKeyToFile(filename string, key *rsa.PrivateKey) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("не вдалося створити файл: %w", err)
	}
	defer file.Close()

	// Кодування приватного ключа у PEM
	err = pem.Encode(file, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err != nil {
		return fmt.Errorf("не вдалося кодувати приватний ключ: %w", err)
	}

	return nil
}

// Функція для запису публічного ключа у файл
func writePublicKeyToFile(filename string, key *rsa.PublicKey) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("не вдалося створити файл: %w", err)
	}
	defer file.Close()

	// Кодування публічного ключа у PEM
	publicKeyBytes := x509.MarshalPKCS1PublicKey(key)
	err = pem.Encode(file, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	if err != nil {
		return fmt.Errorf("не вдалося кодувати публічний ключ: %w", err)
	}

	return nil
}
