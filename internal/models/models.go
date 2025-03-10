package models

import (
	"time"

	"github.com/google/uuid"
)

//	type urlCount struct {
//		url   string
//		count int
//	}

type Article struct {
	ID           uuid.UUID `json:"id"`
	URL          string    `json:"url"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Catagory     string    `json:"catagory"`
	ImageUrl     string    `json:"image_url"`
	Created_at   time.Time `json:"created_at"`
	Published_at time.Time `json:"published_at"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type UserInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Message struct {
	Url  string
	Html string
}

type ClassifyCatagory struct {
	Catagory string `json:"catagory"`
}

// func printPages(pages map[string]int, baseURL string) {
// 	fmt.Println("=============================")
// 	fmt.Println("REPORT for", baseURL)
// 	fmt.Println("=============================")
// 	var urlCounts []urlCount
// 	for url, count := range pages {
// 		urlCounts = append(urlCounts, urlCount{url, count})
// 	}

// 	// Sort the slice (you can use sort.Slice)
// 	sort.Slice(urlCounts, func(i, j int) bool {
// 		return urlCounts[i].count > urlCounts[j].count // '>' for descending order
// 	})

// 	// Print sorted results
// 	for _, uc := range urlCounts {
// 		if uc.count > 1 {
// 			fmt.Printf("Found %d internal links to %s\n", uc.count, uc.url)
// 		} else {
// 			fmt.Printf("Found %d internal link to %s\n", uc.count, uc.url)
// 		}
// 		// fmt.Printf("%s: %d\n", uc.url, uc.count)
// 	}
// }
