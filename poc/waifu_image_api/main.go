package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/yusufpapurcu/maid-bot/poc/waifu_image_api/internal/dto"
	gomail "gopkg.in/gomail.v2"
)

const EMAIL_PASS = "***"

func main() {
	params := QueryParameters{
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

	send(waifuList.Images[0].URL, waifuList.Images[0].Width, waifuList.Images[0].Height)

}

func createURL(params QueryParameters) string {
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

type QueryParameters struct {
	SelectedTags []string
	ExcludedTags []string
	IsNSFW       bool
	Gif          bool
	OrderBy      string
	Orientation  string
	Many         bool
}
