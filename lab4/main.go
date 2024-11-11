package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func main() {

	privateKeyPath := "G:\\lpnu\\7sem\\Software and Data Protection\\keys\\private_key.pem"
	publicKeyPath := "G:\\lpnu\\7sem\\Software and Data Protection\\keys\\public_key.pem"

	privatekey, publicKey, err := loadKeys(privateKeyPath, publicKeyPath)
	if err != nil {
		fmt.Println("Помилка при завантаженні ключів:", err)
		return
	}
	keyhash := md5([]byte("passWord"))
	err, time := encryptFileRC5("test.txt", "encrypted_rc5.txt", keyhash[:])
	if err != nil {
		fmt.Printf("Помилка при шифруванні rc5 %v", err)
		return
	}
	fmt.Printf("Файл зашифровано rc5 за %v", time)

	err, time = encryptFileRSA("test.txt", "ecnrypted_rsa.txt", publicKey)
	if err != nil {
		fmt.Printf("Помилка при шифруванні rsa %v", err)
		return
	}
	fmt.Printf("\nФайл зашифровано rsa за %v", time)

	err, time = decryptFileRC5("encrypted_rc5.txt", "decrypted_rc5.txt", keyhash[:])
	if err != nil {
		fmt.Printf("Помилка при дешифруванні rc5 %v", err)
		return
	}
	fmt.Printf("\nФайл розшифровано rc5 за %v", time)

	err, time = decryptFileRSA("ecnrypted_rsa.txt", "decrypted_rsa.txt", privatekey)
	if err != nil {
		fmt.Printf("Помилка при дешифруванні rsa %v", err)
		return
	}
	fmt.Printf("\nФайл розшифровано rsa за %v", time)
}
func loadKeys(privateKeyPath, publicKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// Завантаження приватного ключа
	privKeyPEM, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("не вдалося прочитати приватний ключ: %w", err)
	}

	privBlock, _ := pem.Decode(privKeyPEM)
	if privBlock == nil || privBlock.Type != "PRIVATE KEY" {
		return nil, nil, fmt.Errorf("не вдалося декодувати приватний ключ")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("не вдалося розпарсити приватний ключ: %w", err)
	}

	// Завантаження публічного ключа
	pubKeyPEM, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("не вдалося прочитати публічний ключ: %w", err)
	}

	pubBlock, _ := pem.Decode(pubKeyPEM)
	if pubBlock == nil || pubBlock.Type != "PUBLIC KEY" {
		return nil, nil, fmt.Errorf("не вдалося декодувати публічний ключ")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("не вдалося розпарсити публічний ключ: %w", err)
	}

	return privateKey, publicKey, nil
}
