package crawl

import (
	"crypto/rand"
	"log"
	"math/big"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 15_4 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Version/15.4 Mobile/15E148 Safari/537.36",
	"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36",
}

func getRandomUserAgent() string {
	max := big.NewInt(int64(len(userAgents)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Printf("failed to get secure random index: %v", err)
		return userAgents[0]
	}
	return userAgents[n.Int64()]
}
