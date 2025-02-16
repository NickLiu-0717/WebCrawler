package main

import (
	"fmt"
	"net/http"

	"github.com/temoto/robotstxt"
)

func checkRobotsTxt(baseURL string) (*robotstxt.Group, error) {
	robotsURL := baseURL + "/robots.txt"

	resp, err := http.Get(robotsURL)
	if err != nil {
		fmt.Println("No robots.txt on:", baseURL)
		return nil, nil
	}
	defer resp.Body.Close()

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("couldn't analyze robots.txt: %v", err)
	}

	return data.FindGroup("*"), nil

}
