package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	database "github.com/NickLiu-0717/crawler/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const defaultMaxConcurrency = 5
const defaultMaxPages = 10
const defaultMaxDepth = 5

// var newsPages []string = []string{
// 	"https://www.bbc.com/",
// 	"https://edition.cnn.com/",
// 	"https://www.forbes.com/",
// 	"https://news.ebc.net.tw",
// 	"https://news.ltn.com.tw",
// 	"https://news.pts.org.tw/",
// }

func main() {
	var maxConcurrency int
	var maxPages int
	var maxDepth int
	var err error

	cmdArg := os.Args

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)

	switch len(cmdArg) {
	case 1:
		fmt.Println("Error - need arguments: website maxConcurrency maxPages")
		os.Exit(1)
	case 2:
		maxConcurrency = defaultMaxConcurrency
		maxPages = defaultMaxPages
		maxDepth = defaultMaxDepth
	case 3:
		maxConcurrency, err = strconv.Atoi(cmdArg[2])
		if err != nil {
			fmt.Println("Error - couldn't convert input string to integer")
			return
		}
		maxPages = defaultMaxPages
		maxDepth = defaultMaxDepth
	case 4:
		maxConcurrency, err = strconv.Atoi(cmdArg[2])
		if err != nil {
			fmt.Println("Error - couldn't convert input string to integer")
			return
		}
		maxPages, err = strconv.Atoi(cmdArg[3])
		if err != nil {
			fmt.Println("Error - couldn't convert input string to integer")
			return
		}
		maxDepth = defaultMaxDepth
	case 5:
		maxConcurrency, err = strconv.Atoi(cmdArg[2])
		if err != nil {
			fmt.Println("Error - couldn't convert input string to integer")
			return
		}
		maxPages, err = strconv.Atoi(cmdArg[3])
		if err != nil {
			fmt.Println("Error - couldn't convert input string to integer")
			return
		}
		maxDepth, err = strconv.Atoi(cmdArg[4])
		if err != nil {
			fmt.Println("Error - couldn't convert input string to integer")
			return
		}
	default:
		fmt.Println("Error - too many arguments provided")
		os.Exit(1)
	}

	cfg, err := configure(cmdArg[1], maxConcurrency, maxPages, maxDepth)
	if err != nil {
		fmt.Printf("Error - couldn't configure: %v\n", err)
		return
	}
	cfg.db = dbQueries

	group, err := checkRobotsTxt(cfg.baseURL.String())
	if err != nil {
		fmt.Printf("Error - couldn't check robots.txt: %v\n", err)
		return
	}
	cfg.robotGroup = group

	cfg.wg.Add(1)
	go cfg.crawlPage(cfg.baseURL.String(), 1)
	cfg.wg.Wait()

	// printPages(cfg.pages, strings.TrimSuffix(cfg.baseURL.String(), "/"))
	for key, article := range cfg.articles {
		fmt.Printf("From: %s, Title: %s\n", key, article.title)
	}
}
