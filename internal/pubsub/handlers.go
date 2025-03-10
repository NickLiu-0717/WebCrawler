package pubsub

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/NickLiu-0717/crawler/config"
	"github.com/NickLiu-0717/crawler/crawl"
	"github.com/NickLiu-0717/crawler/internal/models"
	"github.com/NickLiu-0717/crawler/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

var domainLimits = make(map[string]chan struct{})
var domainLimitsMu sync.Mutex

func HandlerExtractArticle(cfg *crawl.CrawlConfig) func(msg models.Message) Acktype {
	return func(msg models.Message) Acktype {

		cfg.Config.ConcurrencyControl <- struct{}{}
		defer func() {
			<-cfg.Config.ConcurrencyControl
			// cfg.Config.Wg.Done()
		}()

		if !service.CheckArticle(msg.Url) {
			return Ack
		}

		norURL, err := crawl.NormalizeURL(msg.Url)
		if err != nil {
			fmt.Printf("couldn't normalize current URL: %v\n", err)
			return NackDiscard
		}

		gottitle, gotcontent, gotpubAt, err := service.ExtractArticles(msg.Html)
		if err != nil {
			fmt.Printf("couldn't extract title and content from HTML: %v\n", err)
			return NackDiscard
		}

		gotcatagory, err := service.ClassifyArticle(gottitle)
		if err != nil {
			fmt.Printf("couldn't classify article by title: %v\n", err)
			return NackDiscard
		}

		err = cfg.CreateArticles(norURL, gottitle, gotcontent, gotcatagory, gotpubAt.UTC())
		if err != nil {
			return NackDiscard
		}
		// fmt.Printf("create article to database: %v\n", gottitle)
		return Ack
	}
}

func HandlerExtractURLs(cfg *crawl.CrawlConfig, pubChannel *amqp.Channel) func(msg models.Message) Acktype {
	return func(msg models.Message) Acktype {

		cfg.Config.ConcurrencyControl <- struct{}{}
		defer func() {
			<-cfg.Config.ConcurrencyControl
			// cfg.Config.Wg.Done()
		}()

		cfg.Config.Mu.Lock()
		if len(cfg.Config.Pages) >= cfg.Config.MaxPages {
			fmt.Println("Reached max pages limit, stopping further crawling.")
			cfg.Config.Mu.Unlock()
			return Ack
		}
		cfg.Config.Mu.Unlock()

		norURL, err := crawl.NormalizeURL(msg.Url)
		if err != nil {
			fmt.Printf("couldn't normalize current URL: %v\n", err)
			return NackDiscard
		}

		cfg.Config.Mu.Lock()
		if _, ok := cfg.Config.Pages[norURL]; ok {
			cfg.Config.Pages[norURL]++
			cfg.Config.Mu.Unlock()
			return Ack
		}
		cfg.Config.Pages[norURL] = 1
		cfg.Config.Mu.Unlock()

		URLs, err := crawl.GetURLsFromHTML(msg.Html, msg.Url)
		if err != nil {
			fmt.Printf("couldn't get URL from HTML: %v\n", err)
			return NackDiscard
		}

		baseHost := map[string]int{} //store hostname of baseURLs
		for _, baseURL := range cfg.Config.BaseURLs {
			parsedBaseURL, err := url.Parse(baseURL)
			if err != nil {
				fmt.Printf("couldn't parse url: %v\n", err)
				continue
			}
			baseHost[parsedBaseURL.Hostname()]++
		}
		// var wg sync.WaitGroup
		for _, nextUrl := range URLs {
			// wg.Add(1)
			go func(nextUrl string) {
				// defer wg.Done()
				cfg.Config.Mu.Lock()
				if len(cfg.Config.Pages) >= cfg.Config.MaxPages {
					fmt.Println("Reached max pages limit, stopping further crawling.")
					cfg.Config.Mu.Unlock()
					return
				}
				cfg.Config.Mu.Unlock()

				parsed, err := url.Parse(nextUrl)
				if err != nil {
					fmt.Printf("couldn't parse url: %v\n", err)
					return
				}
				for _, group := range cfg.Config.RobotGroup {
					if group != nil {
						if !group.Test(parsed.Path) {
							fmt.Printf("URL %s is not allowed to crawl\n", parsed.String())
							return
						}
					}
				}

				if _, ok := baseHost[parsed.Hostname()]; ok { //chech if the url is under desired hostname
					domainChan := getDomainLimit(parsed.Hostname(), 10)
					domainChan <- struct{}{}
					defer func() { <-domainChan }()
					config.RandomSleep(2000, 5000)
					html, err := crawl.GetHTML(parsed.String())
					if err != nil {
						fmt.Printf("Error - couldn't get HTML: %v\n", err)
						return
					}

					msg := models.Message{
						Url:  parsed.String(),
						Html: html,
					}

					err = Publish(pubChannel, "crawl_page", "page.process.url", msg)
					// fmt.Printf("Publish url: %v\n", msg.Url)
					if err != nil {
						fmt.Printf("couldn't publish message: %v\n", err)
					}
				}
			}(nextUrl)
		}
		// wg.Wait()
		return Ack
	}

}

func getDomainLimit(domain string, limit int) chan struct{} {
	domainLimitsMu.Lock()
	defer domainLimitsMu.Unlock()
	if _, exists := domainLimits[domain]; !exists {
		domainLimits[domain] = make(chan struct{}, limit)
	}
	return domainLimits[domain]
}
