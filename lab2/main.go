package main

import (
	"fmt"
	"os"

	"github.com/mymmrac/telego"
	//"github.com/mymmrac/telego/telegoapi"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {

	botToken := "7468936018:AAEzA2LP9qA2xu7awLB5NOIml1Rciz8yDlE"
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		return
	}
	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	userMode := make(map[int64]string)
	passWord := ""
	for update := range updates {
		if update.Message != nil {
			userMessage := update.Message.Text
			chatId := update.Message.Chat.ID
			if userMessage == "/setpassword" {
				// key := md5([]byte(passWord))
				// fmt.Println(passWord + " " + userMessage)
				// encryptedMessage := encryptCBC([]byte(userMessage), key[:], []byte("1111"))
				// _, err := bot.SendMessage(&telego.SendMessageParams{
				// 	ChatID: telego.ChatID{ID: chatId},
				// 	Text:   hex.EncodeToString(encryptedMessage),
				// })
				// if err != nil {
				// 	fmt.Println("Помилка відправки повідомлення:", err)
				// }
				// passWord = ""
				userMode[chatId] = "setpassword"
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   "Придумайте пароль, ви можете змінити його в будь-який час!",
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
			} else if mode, ok := userMode[chatId]; ok && mode == "setpassword" {
				userMode[chatId] = ""
				passWord = userMessage
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   "Пароль змінено!",
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
			} else if userMessage == "/getpassword" {
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   "Пароль : " + passWord,
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
			} else if userMessage == "/md5" {
				userMode[chatId] = "md5"
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   "Напиши будь-яку фразу чи слово, або надішли файл, і я її захешую за допомогою MD5 (;",
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
			} else if userMessage == "/rc5" {
				userMode[chatId] = "rc5"
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   "/rc5encrypt для шифрування та /rc5decrypt для дешифрування",
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
			} else if userMessage == "/rc5encrypt" {
				userMode[chatId] = "rc5_encrypt"
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   "Напиши що треба зашифрувати, або надішли файл до 20 мб.",
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
			} else if userMessage == "/rc5decrypt" {
				userMode[chatId] = "rc5_decrypt"
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   "Напиши що треба розшифрувати, або надішли файл до 20 мб.",
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
			} else if userMessage == "/cancel" {
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   "Режим скинут. Використовуй /md5 або /rc5 для вибору режиму.",
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
				userMode[chatId] = ""
			} else if mode, ok := userMode[chatId]; ok && mode == "md5" && update.Message.Document == nil {
				///////////////////////////////////////////////////////////////////////////////////////////// MD5
				hash := md5([]byte(userMessage))
				response := fmt.Sprintf("MD5 хеш повідомлення \"%s\": %x", userMessage, hash)

				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   response,
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
			} else if mode, ok := userMode[chatId]; ok && mode == "rc5_encrypt" && update.Message.Document == nil {
				/////////////////////////////////////////////////////////////////////////////////////////// RC5
				response := ""
				fileToSend := "encrypted"
				if passWord != "" {
					keyhash := []byte(passWord)
					err := encryptStringToFile([]byte(userMessage), keyhash[:], fileToSend)
					if err != nil {
						response = "Не вдалося зашифрувати : " + fmt.Sprintf("%v", err)
					} else {
						response = "Зашифровано у файл"
						newfile, err := os.Open(fileToSend)
						if err != nil {
							fmt.Sprintf("не вдалося відкрити файл: %v", err)
						}
						defer newfile.Close()
						document := tu.Document(
							telego.ChatID{ID: chatId},
							tu.File(newfile),
						).WithCaption("decrypted")

						msg, err := bot.SendDocument(document)
						if err != nil {
							fmt.Println(err)
							return
						}
						fmt.Println(msg.Document)
					}

				} else {
					response = "Парольна фраза відстуня, спробуйте /setpassword"
				}
				_, err1 := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   response,
				})
				if err1 != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}

			} else if mode, ok := userMode[chatId]; ok && mode == "rc5_decrypt" && update.Message.Document == nil {
				response := ""
				response = "Розшифрований текст : "
				if passWord != "" {
					keyhash := []byte(passWord)

					err := decryptFile("encrypted.txt", "decrypted.txt", keyhash[:])
					if err != nil {
						fmt.Println("Помилка дешифрування:", err)
						response = response + "Не вдалося розшифрувати"
					} else {
						response = "Розшифровано"
					}

				} else {
					response = "Парольна фраза відстуня, спробуйте /setpassword"
				}
				_, err1 := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   response,
				})
				if err1 != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
			} else if mode, ok := userMode[chatId]; ok && mode == "md5" && update.Message.Document != nil {
				fileId := update.Message.Document.FileID
				file, err := bot.GetFile(&telego.GetFileParams{FileID: fileId})
				if err != nil {
					fmt.Println("Помилка отримання файлу:", err)
					continue
				}

				fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", botToken, file.FilePath)
				localFilePath := "received_file"
				err = downloadFile(fileURL, localFilePath)
				if err != nil {
					fmt.Println("Помилка завантаження файлу:", err)
					continue
				}
				var response string
				if mode, ok := userMode[chatId]; ok && mode == "md5" {
					fileHash, err := md5FromFile(localFilePath)
					if err != nil {
						fmt.Println("Помилка обчислення MD5:", err)
						continue
					}

					response = fmt.Sprintf("MD5 хеш файлу: %x", fileHash)
				} else {
					response = "Невідома команда. Використовуй /md5 або /rc5 для вибору режиму."
				}

				_, err = bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   response,
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
				os.Remove(localFilePath)
			} else if mode, ok := userMode[chatId]; ok && mode == "rc5_encrypt" && update.Message.Document != nil {
				fileId := update.Message.Document.FileID
				file, err := bot.GetFile(&telego.GetFileParams{FileID: fileId})
				if err != nil {
					fmt.Println("Помилка отримання файлу:", err)
					continue
				}
				fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", botToken, file.FilePath)
				localFilePath := "received_file"
				err = downloadFile(fileURL, localFilePath)
				if err != nil {
					fmt.Println("Помилка завантаження файлу:", err)
					continue
				}
				keyhash := []byte(passWord)
				encryptedFilePath := "file_to_send"
				err1 := encryptFile(localFilePath, encryptedFilePath, keyhash[:])
				if err1 != nil {
					fmt.Println("Помилка обчислення rc5:", err)
					continue
				}
				newfile, err := os.Open(encryptedFilePath)
				if err != nil {
					fmt.Sprintf("не вдалося відкрити файл: %v", err)
				}
				defer newfile.Close()
				document := tu.Document(
					telego.ChatID{ID: chatId},
					tu.File(newfile),
				).WithCaption("encrypted")

				msg, err := bot.SendDocument(document)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(msg.Document)
			} else if mode, ok := userMode[chatId]; ok && mode == "rc5_decrypt" && update.Message.Document != nil {
				fileId := update.Message.Document.FileID
				file, err := bot.GetFile(&telego.GetFileParams{FileID: fileId})
				if err != nil {
					fmt.Println("Помилка отримання файлу:", err)
					continue
				}
				fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", botToken, file.FilePath)
				localFilePath := "received_file"
				err = downloadFile(fileURL, localFilePath)
				if err != nil {
					fmt.Println("Помилка завантаження файлу:", err)
					continue
				}
				keyhash := []byte(passWord)
				decryptedFilePath := "file_to_send"
				err1 := decryptFile(localFilePath, decryptedFilePath, keyhash[:])
				if err1 != nil {
					fmt.Println("Помилка обчислення rc5:", err)
					continue
				}
				newfile, err := os.Open(decryptedFilePath)
				if err != nil {
					fmt.Sprintf("не вдалося відкрити файл: %v", err)
				}
				defer newfile.Close()
				document := tu.Document(
					telego.ChatID{ID: chatId},
					tu.File(newfile),
				).WithCaption("decrypted")

				msg, err := bot.SendDocument(document)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(msg.Document)
			} else {
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   "Невідома команда. Використовуй /md5 або /rc5 для вибору режиму.",
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
			}
		}
	}
	defer bot.StopLongPolling()
}
