package main

import (
	"fmt"

	"github.com/mymmrac/telego"
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

	for update := range updates {
		if update.Message != nil {
			userMessage := update.Message.Text
			chatId := update.Message.Chat.ID

			if update.Message.Document != nil {
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

				fileHash, err := md5FromFile(localFilePath)
				if err != nil {
					fmt.Println("Помилка обчислення MD5:", err)
					continue
				}

				response := fmt.Sprintf("MD5 хеш файлу: %x", fileHash)
				_, err = bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   response,
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}

				//os.Remove(localFilePath)
			} else if userMessage != "/start" {
				hash := md5([]byte(userMessage))
				response := fmt.Sprintf("MD5 хеш повідомлення \"%s\": %x", userMessage, hash)

				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   response,
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
			} else {
				_, err := bot.SendMessage(&telego.SendMessageParams{
					ChatID: telego.ChatID{ID: chatId},
					Text:   "Напиши будь-яку фразу чи слово, або надішли файл, і я її захешую (;",
				})
				if err != nil {
					fmt.Println("Помилка відправки повідомлення:", err)
				}
			}
		}
	}
	defer bot.StopLongPolling()
}
