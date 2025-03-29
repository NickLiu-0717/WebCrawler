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

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	if dbURL == "" {
		log.Fatal("PostgreSQL environment variables must be set")
	}
	amqpString := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)
	if amqpString == "" {
		log.Fatal("RabbitMQ environment variables must be set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("secretKey must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)
	defer db.Close()

	baseURLs := []string{
		// "https://www.bbc.com/",
		"https://edition.cnn.com/",
		// "https://www.forbes.com/",
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
