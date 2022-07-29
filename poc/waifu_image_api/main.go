package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/yusufpapurcu/maid-bot/poc/waifu_image_api/internal/dto"
)

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

	fmt.Println(waifuList.Images[0].URL)

}

func createURL(params QueryParameters) string {
	res, _ := url.Parse("https://api.waifu.im/random")
	res.ForceQuery = true

	query := res.Query()
	query.Add("orientation", params.Orientation)

	res.RawQuery = query.Encode()
	return res.String()
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
