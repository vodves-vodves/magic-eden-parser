package main

import (
	"fmt"
	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

const token = "MTAxODI0OTE0NjMxNjEwNzc4Ng.GIWUSj.wYI8Caw2UPq43Vyz0bCcQd49L7FZ2-hKIDuy-E"

var botId string

func main() {
	start()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	//floorCollections("utilityf_ape")
	//popularCollections()
}

//запуск discord bot
func start() {
	goBot, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("1", err)
		return
	}
	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err)
		return
	}
	botId = u.ID
	goBot.AddHandler(messageHandler)
	err = goBot.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("bot started")
}

//обработчик сообщений
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == botId {
		return
	}

	if strings.HasPrefix(m.Content, "!floor") {
		//rename this
		var test string
		collSymbol := strings.TrimSpace(strings.Split(m.Content, "!floor")[1])
		fmt.Println(collSymbol)
		result := floorCollections(collSymbol)
		if len(result) == 0 {
			_, err := s.ChannelMessageSend(m.ChannelID, "Collection not found!")
			if err != nil {
				fmt.Println("err: ", err)
				return
			}
		} else {
			for _, item := range result {
				test += item
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Запрос "+fmt.Sprintf("%s#%s", m.Author.Username, m.Author.Discriminator), test, 990099))
			if err != nil { //800080
				fmt.Println("err: ", err)
				return
			}
		}
		fmt.Println("message sended")
	}
}

//Выводит флор определенной коллекции
func floorCollections(collectionSymbol string) []string {
	var result []string
	url := "https://api-mainnet.magiceden.io/rpc/getListedNFTsByQueryLite"
	query := map[string]string{
		"q": "{\"$match\":{\"collectionSymbol\":\"" + collectionSymbol + "\"},\"$sort\":{\"takerAmount\":1},\"$skip\":0,\"$limit\":10,\"status\":[]}",
	}
	json, err := sendRequest(url, query)
	if err != nil {
		fmt.Println("err: ", err)
	}
	value := gjson.Get(json, "results").Array()
	for _, item := range value {
		result = append(result, fmt.Sprintf("%s - %s $SOL - https://www.magiceden.io/item-details/%s\n\n", item.Get("title").String(), item.Get("price").String(), item.Get("mintAddress").String()))
	}
	return result
}

//Выводит популярные коллекции
func popularCollections() {
	var result []string
	url := "https://api-mainnet.magiceden.io/popular_collections"
	query := make(map[string]string)
	json, err := sendRequest(url, query)
	if err != nil {
		fmt.Println("err: ", err)
	}
	value := gjson.Get(json, "collections")
	fmt.Println(result, value)
}

//Отправка запроса
func sendRequest(url string, query map[string]string) (string, error) {
	headers := map[string]string{
		"Authority":  "stats-mainnet.magiceden.io",
		"Accept":     "*/*",
		"User-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36",
	}
	client := http.Client{Timeout: 5 * time.Second}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
	req, _ := http.NewRequest("GET", url, nil)
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
