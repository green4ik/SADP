package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mymmrac/telego"
)

var s = []int{
	7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22,
	5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20,
	4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23,
	6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21,
}

var K = []uint32{
	0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee, 0xf57c0faf, 0x4787c62a, 0xa8304613, 0xfd469501,
	0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be, 0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821,
	0xf61e2562, 0xc040b340, 0x265e5a51, 0xe9b6c7aa, 0xd62f105d, 0x02441453, 0xd8a1e681, 0xe7d3fbc8,
	0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed, 0xa9e3e905, 0xfcefa3f8, 0x676f02d9, 0x8d2a4c8a,
	0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c, 0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70,
	0x289b7ec6, 0xeaa127fa, 0xd4ef3085, 0x04881d05, 0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665,
	0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039, 0x655b59c3, 0x8f0ccc92, 0xffeff47d, 0x85845dd1,
	0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1, 0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391,
}

func leftRotate(x uint32, c int) uint32 {
	return (x << c) | (x >> (32 - c))
}

var initValues = [4]uint32{
	0x67452301,
	0xefcdab89,
	0x98badcfe,
	0x10325476,
}

func md5(input []byte) [16]byte {
	msg := padding(input)
	var a, b, c, d uint32 = initValues[0], initValues[1], initValues[2], initValues[3]

	for i := 0; i < len(msg); i += 64 {
		M := msg[i : i+64]
		var aa, bb, cc, dd = a, b, c, d

		for j := 0; j < 64; j++ {
			var f uint32
			var g int
			switch {
			case j < 16:
				f = (b & c) | (^b & d)
				g = j
			case j < 32:
				f = (d & b) | (^d & c)
				g = (5*j + 1) % 16
			case j < 48:
				f = b ^ c ^ d
				g = (3*j + 5) % 16
			default:
				f = c ^ (b | ^d)
				g = (7 * j) % 16
			}

			f = f + a + K[j] + binary.LittleEndian.Uint32(M[4*g:4*g+4])
			a = d
			d = c
			c = b
			b = b + leftRotate(f, s[j])
		}

		a += aa
		b += bb
		c += cc
		d += dd
	}

	var hash [16]byte
	binary.LittleEndian.PutUint32(hash[0:], a)
	binary.LittleEndian.PutUint32(hash[4:], b)
	binary.LittleEndian.PutUint32(hash[8:], c)
	binary.LittleEndian.PutUint32(hash[12:], d)
	return hash
}

func padding(input []byte) []byte {
	msgLen := len(input)
	input = append(input, 0x80)

	for len(input)%64 != 56 {
		input = append(input, 0x00)
	}

	msgLenBits := uint64(msgLen) * 8
	var lenBytes [8]byte
	binary.LittleEndian.PutUint64(lenBytes[:], msgLenBits)
	input = append(input, lenBytes[:]...)

	return input
}

func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func md5FromFile(filePath string) ([16]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return [16]byte{}, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return [16]byte{}, err
	}
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)
	_, err = file.Read(buffer)
	if err != nil {
		return [16]byte{}, err
	}

	return md5(buffer), nil
}

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
