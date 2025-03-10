package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/NickLiu-0717/crawler/config"
	"github.com/NickLiu-0717/crawler/crawl"
	"github.com/NickLiu-0717/crawler/handler"
	database "github.com/NickLiu-0717/crawler/internal/database"
	"github.com/NickLiu-0717/crawler/internal/models"
	"github.com/NickLiu-0717/crawler/internal/pubsub"
	"github.com/NickLiu-0717/crawler/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

const maxConcurrency = 100
const maxPages = 100
const maxDepth = 5

// var newsPages []string = []string{
// 	"https://www.bbc.com/",
// 	"https://edition.cnn.com/",
// 	"https://www.forbes.com/",
// 	"https://news.ebc.net.tw",
// 	"https://news.pts.org.tw/",
// }

func main() {

	var err error

	// cmdArg := os.Args

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("secretKey must be set")
	}
	amqpString := os.Getenv("AMQP_CONNECT_STRING")
	if amqpString == "" {
		log.Fatal("secretKey must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)
	defer db.Close()

	// switch len(cmdArg) {
	// case 1:
	// 	fmt.Println("Error - need arguments: website maxConcurrency maxPages")
	// 	os.Exit(1)
	// case 2:
	// 	maxConcurrency = defaultMaxConcurrency
	// 	maxPages = defaultMaxPages
	// 	maxDepth = defaultMaxDepth
	// case 3:
	// 	maxConcurrency, err = strconv.Atoi(cmdArg[2])
	// 	if err != nil {
	// 		fmt.Println("Error - couldn't convert input string to integer")
	// 		return
	// 	}
	// 	maxPages = defaultMaxPages
	// 	maxDepth = defaultMaxDepth
	// case 4:
	// 	maxConcurrency, err = strconv.Atoi(cmdArg[2])
	// 	if err != nil {
	// 		fmt.Println("Error - couldn't convert input string to integer")
	// 		return
	// 	}
	// 	maxPages, err = strconv.Atoi(cmdArg[3])
	// 	if err != nil {
	// 		fmt.Println("Error - couldn't convert input string to integer")
	// 		return
	// 	}
	// 	maxDepth = defaultMaxDepth
	// case 5:
	// 	maxConcurrency, err = strconv.Atoi(cmdArg[2])
	// 	if err != nil {
	// 		fmt.Println("Error - couldn't convert input string to integer")
	// 		return
	// 	}
	// 	maxPages, err = strconv.Atoi(cmdArg[3])
	// 	if err != nil {
	// 		fmt.Println("Error - couldn't convert input string to integer")
	// 		return
	// 	}
	// 	maxDepth, err = strconv.Atoi(cmdArg[4])
	// 	if err != nil {
	// 		fmt.Println("Error - couldn't convert input string to integer")
	// 		return
	// 	}
	// default:
	// 	fmt.Println("Error - too many arguments provided")
	// 	os.Exit(1)
	// }

	baseURLs := []string{
		// "https://www.bbc.com/",
		"https://edition.cnn.com/",
		"https://www.forbes.com/",
		"https://news.ebc.net.tw",
		"https://news.pts.org.tw/",
	}

	conf, err := config.Configure(baseURLs, maxConcurrency, maxPages, maxDepth)
	if err != nil {
		fmt.Printf("Error - couldn't configure: %v\n", err)
		return
	}
	cfg := crawl.CrawlConfig{Config: conf}
	cfg.Config.Db = dbQueries

	//check if the base URLs contain robots.txt
	for _, bURL := range cfg.Config.BaseURLs {
		group, err := service.CheckRobotsTxt(bURL)
		if err != nil {
			fmt.Printf("Error - couldn't check robots.txt: %v\n", err)
			return
		}
		cfg.Config.RobotGroup = append(cfg.Config.RobotGroup, group)
	}

	apiConfig := &config.ApiConfig{
		Db:        dbQueries,
		Port:      port,
		SecretKey: secretKey,
	}
	apicfg := handler.Handler{Config: apiConfig}
	totalPages, err := apicfg.GetTotalPages(5)
	if err != nil {
		fmt.Printf("Error - couldn't get total pages: %v\n", err)
	}
	apicfg.Config.TotalPages = totalPages

	// if apicfg.Config.TotalPages < 10 {
	// 	// cfg.Config.Wg.Add(1)
	// 	fmt.Println("Start crawling...")
	// 	go cfg.CrawlPage(cfg.Config.BaseURL.String(), 1)
	// 	// cfg.Config.Wg.Wait()
	// }

	// printPages(cfg.pages, strings.TrimSuffix(cfg.baseURL.String(), "/"))
	// for key, article := range cfg.articles {
	// 	fmt.Printf("From: %s, Title: %s\n", key, article.title)
	// }

	connc, err := amqp.Dial(amqpString)
	if err != nil {
		fmt.Printf("Error - couldn't connect to RabbitMQ: %v\n", err)
		return
	}
	if err = pubsub.SetupRabbitMQ(connc); err != nil {
		fmt.Printf("Error - couldn't set up exchanges and queues: %v\n", err)
		return
	}
	defer connc.Close()

	ch, err := connc.Channel()
	if err != nil {
		fmt.Printf("Error - couldn't open channel: %v\n", err)
		return
	}

	for _, baseURL := range cfg.Config.BaseURLs {
		go func(baseURL string) {
			parsedURL, err := url.Parse(baseURL)
			if err != nil {
				fmt.Printf("Error - couldn't parse url %s: %v\n", baseURL, err)
				return
			}
			html, err := crawl.GetHTML(parsedURL.String())
			if err != nil {
				fmt.Printf("Error - couldn't get HTML: %v\n", err)
				return
			}
			msg := models.Message{
				Url:  parsedURL.String(),
				Html: html,
			}
			err = pubsub.Publish(ch, pubsub.ExchangeCrawlPageTopic, pubsub.QueueURL, msg)
			if err != nil {
				fmt.Printf("couldn't publish message: %v\n", err)
				return
			}
		}(baseURL)
	}

	err = pubsub.Subscribe(
		connc,
		pubsub.ExchangeCrawlPageTopic,
		"topic",
		pubsub.QueueURL,
		pubsub.CrawlKeyPrefix+".*",
		pubsub.SimpleQueueTransient,
		pubsub.HandlerExtractURLs(&cfg, ch),
	)
	if err != nil {
		fmt.Printf("Error - couldn't subscribe to page.process.URL: %v\n", err)
	}

	err = pubsub.Subscribe(
		connc,
		pubsub.ExchangeCrawlPageTopic,
		"topic",
		pubsub.QueueArticle,
		pubsub.CrawlKeyPrefix+".*",
		pubsub.SimpleQueueTransient,
		pubsub.HandlerExtractArticle(&cfg),
	)
	if err != nil {
		fmt.Printf("Error - couldn't subscribe to page.process.article: %v\n", err)
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + apicfg.Config.Port,
		Handler: mux,
	}

	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusFound)
	})
	mux.Handle("/app/", fileServer)
	mux.HandleFunc("POST /api/reset", apicfg.HandlerReset)
	mux.HandleFunc("POST /api/signup", apicfg.HandlerSignup)
	mux.HandleFunc("POST /api/login", apicfg.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", apicfg.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", apicfg.HandlerRevoke)
	mux.HandleFunc("GET /api/articles", apicfg.HandlerGetArticles)
	mux.HandleFunc("GET /api/articles/{articleId}", apicfg.HandlerGetArticleFromID)
	mux.HandleFunc("GET /api/categories/{category}/articles", apicfg.HandlerGetCategoryArticles)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
