package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const ExchangeCrawlPageTopic = "crawl_page"
const QueueURL = "page.process.URL"
const QueueArticle = "page.process.article"
const CrawlKeyPrefix = "page.process"

func SetupRabbitMQ(conn *amqp.Connection) error {
	if err := DeclareDeadLetterSetUp(conn); err != nil {
		return fmt.Errorf("failed to declare DLX: %v", err)
	}

	_, _, err := DeclareAndBind(
		conn,
		ExchangeCrawlPageTopic,
		"topic",
		QueueURL,
		CrawlKeyPrefix+".*",
		SimpleQueueTransient,
	)
	if err != nil {
		return fmt.Errorf("couldn't declare and bind queue page.process.url: %v", err)
	}

	_, _, err = DeclareAndBind(
		conn,
		ExchangeCrawlPageTopic,
		"topic",
		QueueArticle,
		CrawlKeyPrefix+".*",
		SimpleQueueTransient,
	)
	if err != nil {
		return fmt.Errorf("couldn't declare and bind queue page.process.article: %v", err)
	}
	return nil
}
