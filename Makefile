run:
	go run ./cmd/shortener/main.go

re:
	killall -9 shortener || true
	rm -f cmd/shortener/shortener
	go build -o cmd/shortener/shortener cmd/shortener/main.go

PHONY: run re
