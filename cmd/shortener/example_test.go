package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"resty.dev/v3"

	"github.com/domurdoc/shortener/internal/app"
	"github.com/domurdoc/shortener/internal/auth"
	"github.com/domurdoc/shortener/internal/config"
	"github.com/domurdoc/shortener/internal/handler"
	"github.com/domurdoc/shortener/internal/httputil"
	"github.com/domurdoc/shortener/internal/router"
)

func getAppRouter(shortCode string) (*app.App, http.Handler) {
	cfg := config.Default()
	cfg.ServiceGeneratorConstantValue = shortCode
	a, err := app.New(cfg)
	if err != nil {
		panic(err)
	}
	h := handler.New(a)
	r := router.New(h)
	r = httputil.AddMiddlewares(
		r,
		auth.NewAuthMiddleware(a.Auth),
	)
	return a, r
}

func ExampleHandler_Shorten() {
	shortCode := "123456"
	a, r := getAppRouter(shortCode)
	defer func() { _ = a.Close(nil) }()

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := resty.New()
	resp, err := client.
		R().
		SetHeader("Content-Type", "text/plain; charset=utf-8").
		SetBody("http://maloslnieogurchiki.com").
		Post(ts.URL + "/")

	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.Header().Get("Content-Type"))
	fmt.Println(string(result))

	// Output:
	// 201
	// text/plain; charset=utf-8
	// http://localhost:8080/123456
}

func ExampleHandler_ShortenJSON() {
	shortCode := "123456"
	a, r := getAppRouter(shortCode)
	defer func() { _ = a.Close(nil) }()

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := resty.New()
	resp, err := client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"url":"http://maloslnieogurchiki.com"}`).
		Post(ts.URL + "/api/shorten")

	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.Header().Get("Content-Type"))
	fmt.Println(string(result))

	// Output:
	// 201
	// application/json
	// {"result":"http://localhost:8080/123456"}
}

func ExampleHandler_Retrieve() {
	shortCode := "123456"
	a, r := getAppRouter(shortCode)
	defer func() { _ = a.Close(nil) }()

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy())
	resp, err := client.
		R().
		SetHeader("Content-Type", "text/plain; charset=utf-8").
		SetBody("http://maloslnieogurchiki.com").
		Post(ts.URL + "/")

	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		panic(err)
	}

	resp, err = client.
		R().
		Get(ts.URL + "/" + shortCode)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.Header().Get("Location"))

	// Output:
	// 307
	// http://maloslnieogurchiki.com
}
