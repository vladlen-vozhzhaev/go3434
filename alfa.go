package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
)

const openExchangeAppKey = "ad1c3a856f32428db34c670bfefe84c1"
const giphyAPIKey = "yogxacUoXdKZX0wkD0depV72uwdbz3gj"

func getRub(date string) (float64, error) {
	endPoint := "latest.json"
	if date != "" {
		endPoint = "historical/" + date + ".json"
	}
	spec := "https://openexchangerates.org/api/" + endPoint + "?app_id=" + openExchangeAppKey
	resp, err := http.Get(spec)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return 0, err
	}
	rates, ok := parsed["rates"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("no rates found")
	}

	rub, ok := rates["RUB"].(float64)
	if !ok {
		return 0, fmt.Errorf("no RUB found")
	}

	return rub, err
}

func getGifURL(q string) (string, error) {
	randIdx := rand.IntN(24)
	url := "https://api.giphy.com/v1/gifs/search?api_key=" + giphyAPIKey + "&q=" + q
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var parsed struct {
		Data []struct {
			ID string `json:"id"`
		}
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", err
	}
	imageId := parsed.Data[randIdx].ID
	return "https://i.giphy.com/" + imageId + ".gif", err
}
