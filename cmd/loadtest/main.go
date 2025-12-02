package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"resty.dev/v3"

	"github.com/domurdoc/shortener/internal/utils"
)

type Options struct {
	Addr        string
	Requests    int
	Concurrency int
	Duration    time.Duration
}

func main() {
	options := Options{}
	wg := sync.WaitGroup{}

	flag.StringVar(&options.Addr, "a", "http://localhost:8080", "shortener server address")
	flag.IntVar(&options.Concurrency, "c", 50, "number of crazy users")
	flag.DurationVar(&options.Duration, "d", time.Duration(60*time.Second), "black friday time")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), options.Duration)
	defer cancel()

	wg.Add(options.Concurrency)
	for range options.Concurrency {
		crazyID := utils.GenerateRandomString(utils.ALPHA, 4)
		go CrazyUser(ctx, &wg, options.Addr, crazyID)
	}

	wg.Wait()
}

func CrazyUser(ctx context.Context, wg *sync.WaitGroup, baseURL string, crazyID string) {
	defer wg.Done()

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("crazy %s: failed to create cookiejar: %v", crazyID, err)
		return
	}

	c := resty.
		New().
		SetContext(ctx).
		SetBaseURL(baseURL).
		SetCookieJar(jar).
		SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetTransport(&http.Transport{MaxConnsPerHost: 1})
	defer c.Close()

	var (
		counter     int
		originalURL string
	)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		counter++
		originalURL = fmt.Sprintf("https://malosolnieogurchiki.ru/%s/%d/", crazyID, counter)

		shortURL, err := shorten(c, originalURL)
		if err != nil {
			log.Printf("crazy %s: shorten failed: %v", crazyID, err)
			continue
		}
		err = retrieve(c, shortURL)
		if err != nil {
			log.Printf("crazy %s: retrieve failed: %v", crazyID, err)
		}
	}
}

func shorten(c *resty.Client, originalURL string) (string, error) {
	r, err := c.R().SetHeader("Content-Type", "text/plain").SetBody(originalURL).Post("/")
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	if r.StatusCode() != http.StatusCreated && r.StatusCode() != http.StatusConflict {
		return "", fmt.Errorf("request failed: status code = %d", r.StatusCode())
	}
	shortURL, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(shortURL), nil
}

func retrieve(c *resty.Client, shortURL string) error {
	r, err := c.R().Get(shortURL)
	if err != nil {
		return err
	}
	if r.StatusCode() != http.StatusTemporaryRedirect {
		return fmt.Errorf("request failed: status code = %d", r.StatusCode())
	}
	return nil
}
