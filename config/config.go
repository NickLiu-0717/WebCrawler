package config

import (
	"crypto/rand"
	"log"
	"math/big"
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
	if max <= min {
		return
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		log.Printf("failed to generate secure random delay: %v", err)
		return
	}
	delay := int(n.Int64()) + min
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
