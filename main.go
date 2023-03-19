package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// Update recebida via webhook do Telegram
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message contém as informações da mensagem recebida
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// Chat contém as informações do chat
type Chat struct {
	Id int `json:"id"`
}

var (
	message string
)

// Faz o parse do update vindo via webhook
func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		return nil, err
	}
	return &update, nil
}

func getTelegramApi() string {
	return "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"
}

// Envia o texto para o chat do Telegram
func sendTextToTelegramChat(chatId int, text string) (string, error) {
	log.Printf("Sending %s to chat_id %d", text, chatId)

	var telegramApi = getTelegramApi()

	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		},
	)

	if err != nil {
		log.Printf("Erro ao enviar a mensagem para o telegram: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("Erro ao ler o body da resposta do telegram: %s", errRead.Error())
		return "", errRead
	}

	bodyString := string(bodyBytes)
	log.Printf("Body da resposta para o telegram: %s", bodyString)

	return bodyString, nil
}

// Handler para o webhook do Telegram
func handleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("Erro ao fazer parse do update: %s", err.Error())
		return
	}

	if update.Message.Text == "/hello" {
		message = "Hello World!"
	} else {
		message = "I don't understand"
	}

	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.Id, message)
	if errTelegram != nil {
		log.Printf("Recebi erro %s do telegram, o body é %s", errTelegram.Error(), telegramResponseBody)
	} else {
		log.Printf("Mensagem %s distribuída com sucesso para o canal %d", message, update.Message.Chat.Id)
	}
}

func main() {
	http.HandleFunc("/", handleTelegramWebhook)
	log.Fatal(http.ListenAndServe(":80", nil))
}
