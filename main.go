package main

import (
	"fmt"
	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/bwmarrin/discordgo"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"time"
)

const token = "MTAxODI0OTE0NjMxNjEwNzc4Ng.GIWUSj.wYI8Caw2UPq43Vyz0bCcQd49L7FZ2-hKIDuy-E"

var botId string

func main() {
	//start()
	floorCollections("utility_ape")
	//popularCollections()
}
func start() {
	goBot, err := discordgo.New("Bot" + token)
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

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == botId {
		return
	}
	result := floorCollections("utility_ape")
	if m.Content == "!test" {
		for _, item := range result {
			s.ChannelMessageSend(m.ChannelID, item)
		}
	}
}

//Выводит флор определенной коллекции
func floorCollections(collectionSymbol string) []string {
	var result []string
	url := "https://api-mainnet.magiceden.io/rpc/getListedNFTsByQueryLite"
	query := map[string]string{
		"q": "{\"$match\":{\"collectionSymbol\":\"" + collectionSymbol + "\"},\"$sort\":{\"takerAmount\":1},\"$skip\":0,\"$limit\":20,\"status\":[]}",
	}
	json := sendRequest(url, query)
	value := gjson.Get(json, "results").Array()
	for _, item := range value {
		result = append(result, fmt.Sprintf("%s - %s $SOL - https://www.magiceden.io/item-details/%s\n", item.Get("title").String(), item.Get("price").String(), item.Get("mintAddress").String()))
	}
	fmt.Println(result)
	return result
}

//Выводит популярные коллекции
func popularCollections() {
	var result []string
	url := "https://api-mainnet.magiceden.io/popular_collections"
	query := make(map[string]string)
	json := sendRequest(url, query)
	value := gjson.Get(json, "collections")
	fmt.Println(result, value)
}

//Отправка запроса
func sendRequest(url string, query map[string]string) string {
	headers := map[string]string{
		"Authority":  "stats-mainnet.magiceden.io",
		"Accept":     "*/*",
		"User-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36",
	}
	client := http.Client{Timeout: 1 * time.Second}
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
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
