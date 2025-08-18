run:
	go run ./cmd/shortener/main.go

exe:
	./cmd/shortener/shortener

re:
	killall -9 shortener || true
	rm -f cmd/shortener/shortener
	go build -o cmd/shortener/shortener cmd/shortener/main.go

test:
	./shortenertestbeta -test.v -test.run=^TestIteration1$$ -binary-path=cmd/shortener/shortener
	./shortenertestbeta -test.v -test.run=^TestIteration2$$ -source-path=.
	./shortenertestbeta -test.v -test.run=^TestIteration3$$ -source-path=.
	./shortenertestbeta -test.v -test.run=^TestIteration4$$ -binary-path=cmd/shortener/shortener -server-port=8080

PHONY: run exe re test
