package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/NickLiu-0717/crawler/internal/models"
)

func ClassifyArticle(title string) (string, error) {
	url := "http://127.0.0.1:8000/classify/"

	jsonbody, err := json.Marshal(map[string]string{"title": title})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonbody))
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var catagory models.ClassifyCatagory
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(dat, &catagory); err != nil {
		return "", err
	}
	return catagory.Catagory, nil
}
