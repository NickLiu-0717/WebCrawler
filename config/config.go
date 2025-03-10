package config

import (
	"math/rand"
	"sync"
	"time"

	database "github.com/NickLiu-0717/crawler/internal/database"
	"github.com/temoto/robotstxt"
)

type ApiConfig struct {
	Db         *database.Queries
	Port       string
	TotalPages int
	SecretKey  string
}

type Config struct {
	Pages              map[string]int
	BaseURLs           []string
	Mu                 *sync.Mutex
	ConcurrencyControl chan struct{}
	Wg                 *sync.WaitGroup
	MaxPages           int
	RobotGroup         []*robotstxt.Group
	MaxDepth           int
	Db                 *database.Queries
}

func Configure(rawBaseURLs []string, maxConcurrency, maxPages, maxDepth int) (*Config, error) {
	return &Config{
		Pages:              make(map[string]int),
		BaseURLs:           rawBaseURLs,
		Mu:                 &sync.Mutex{},
		ConcurrencyControl: make(chan struct{}, maxConcurrency),
		// Wg:                 &sync.WaitGroup{},
		MaxPages: maxPages,
		MaxDepth: maxDepth,
	}, nil
}

func RandomSleep(min, max int) {
	delay := rand.Intn(max-min) + min
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
