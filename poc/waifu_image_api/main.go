package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/yusufpapurcu/maid-bot/poc/waifu_image_api/internal/dto"
	gomail "gopkg.in/gomail.v2"
)

const EMAIL_PASS = "***"

func main() {

	var waifuList dto.RandomWaifuList
	for {
		waifuList = getRandomWaifu()
		if isWaifuUninque(waifuList.Images[0].URL) {
			break
		}
	}

	send(waifuList.Images[0].URL, waifuList.Images[0].Width, waifuList.Images[0].Height)
	writeWaifuUrlToCache(waifuList.Images[0].URL)
}

func getRandomWaifu() dto.RandomWaifuList {
	params := dto.QueryParameters{
		SelectedTags: []string{},
		ExcludedTags: []string{},
		IsNSFW:       false,
		Gif:          false,
		OrderBy:      "",
		Orientation:  "LANDSCAPE",
	}

	req, err := http.NewRequest(http.MethodGet, createURL(params), http.NoBody)
	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "yusuf-papurcu-maid-poc")

	fmt.Println(req.URL.String())

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Got status code %d", res.StatusCode))
	}

	var waifuList dto.RandomWaifuList
	err = json.NewDecoder(res.Body).Decode(&waifuList)
	if err != nil {
		panic(err)
	}

	return waifuList
}

func createURL(params dto.QueryParameters) string {
	res, _ := url.Parse("https://api.waifu.im/search")
	res.ForceQuery = true

	query := res.Query()
	query.Add("orientation", params.Orientation)

	res.RawQuery = query.Encode()
	return res.String()
}

func send(waifu_url string, width, height int) {
	// Configuration
	from := "yusufturhanp@gmail.com"
	pass := EMAIL_PASS
	to := "yusufpapurcu@gmail.com"

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "おはようございます Master "+time.Now().Format("01/02/2006"))
	m.SetBody("text/html", fmt.Sprintf(`<img src="%s" alt="おはようございます Master" width="%d" height="%d">`, waifu_url, width, height))

	d := gomail.NewDialer("smtp.gmail.com", 587, from, pass)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	log.Println("Successfully sended to " + to)
}

func isWaifuUninque(url string) bool {
	filename := "uniquness.cache"
	var f *os.File
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	} else {
		f, _ = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0660)
	}

	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			fmt.Println(err)
			break
		}

		if url == string(line) {
			return false
		}
	}
	return true
}

func writeWaifuUrlToCache(url string) {
	filename := "uniquness.cache"
	var f *os.File
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	} else {
		f, _ = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0660)
	}

	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(url + "\n")
	w.Flush()
}
