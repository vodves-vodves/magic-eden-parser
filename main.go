package main

import (
	"fmt"
	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	floorCollections("utility_ape")
	//popularCollections()
}

//Выводит флор определенной коллекции
func floorCollections(collectionSymbol string) {
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
}

//Выводит популярные коллекции
func popularCollections() {
	url := "https://api-mainnet.magiceden.io/popular_collections"
	query := make(map[string]string)
	result := sendRequest(url, query)
	fmt.Println(result)
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
