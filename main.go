package main

import (
	"be-alquran-api/db"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	telegrambot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := telegrambot.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Printf("failed to connect telegram bot %v", err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := telegrambot.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	// Mengatur perintah khusus pada bot
	for update := range updates {
		if update.Message.Text == "" {
			continue
		}
		if update.Message.IsCommand() {
			msg := telegrambot.NewMessage(update.Message.Chat.ID, "")
			msg.ReplyToMessageID = update.Message.MessageID

			switch update.Message.Command() {
			case "start":
				msg.Text = "Halo, saya bot yang membantu anda untuk mendapatkan pesan suara dan teks dari Al-Qur'an, berdasarkan KEMENAG \n\nKetik /help untuk melihat perintah yang tersedia"
			case "help":
				msg.Text = "Ketik /audio + nomor surah untuk mendapatkan pesan suara, contoh: /audio 1 \n\nKetik /surah + nomor surah dan ayat untuk mendapatkan pesan teks, contoh: /surah 1:1"
			case "surah":
				msg.Text = "Masukan surah dan ayat yang diinginkan, contoh: /surah 1:1"
				if len(update.Message.CommandArguments()) > 0 {
					// check if the command arguments is valid
					arguments := update.Message.CommandArguments()
					args := strings.Split(arguments, ":")
					numSurah, _ := strconv.ParseInt(args[0], 10, 64)
					numAyah, _ := strconv.ParseInt(args[1], 10, 64)

					data, err := db.FindSurahAndAyahByNumber(numSurah, numAyah)
					if err != nil {
						text := fmt.Sprintf("Maaf, tidak terdapat surah %d ayat %d", numSurah, numAyah)
						msg.Text = text
						break
					}
					text := data.NameSurah + "\n\n" + data.Verses.Text + "\n\n[IND] " + data.Verses.TranslationID + "\n[ENG] " + data.Verses.TranslationEn
					msg.Text = text
				}
			case "audio":
				msg.Text = "Masukan surah dan ayat yang ingin di dengar, contoh: /audio 1"
				if len(update.Message.CommandArguments()) > 0 {
					// send audio
					arguments := update.Message.CommandArguments()
					stringsArgs := string(arguments)
					numSurah, _ := strconv.ParseInt(stringsArgs, 10, 64)
					data, err := db.FindAudioSurahNumber(numSurah)
					if err != nil {
						text := fmt.Sprintf("Maaf, tidak terdapat surah %d", numSurah)
						msg.Text = text
						break
					}
					caption := fmt.Sprintf("Name Reciter : %s \nNama Surah Arab : %s \nNama Surah Indonesia : %s \nNama Surah English : %s \nJumlah Ayat : %d \nPlace : %s \nType: %s", data.NameReciter, data.NameSurahAR, data.NameSurahIND, data.NameSurahENG, data.NumberOfAyah, data.Place, data.Type)
					fileName := fmt.Sprintf("%s.mp3", data.NameSurahAR)
					SendAudio(data.Audio, bot, update, caption, fileName)
				}
			default:
				msg.Text = "Maaf, perintah tidak dikenali"
			}
			// check timeout message
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}

func SendAudio(fileUrl string, bot *telegrambot.BotAPI, update telegrambot.Update, caption, fileName string) error {
	response, err := http.Get(fileUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fileBytes := telegrambot.FileBytes{
		Name:  fileName,
		Bytes: data,
	}
	document := telegrambot.NewAudio(update.Message.Chat.ID, fileBytes)
	document.Caption = caption

	if _, err := bot.Send(document); err != nil {
		log.Panic(err)
	}
	return nil
}
