run:
	go run ./cmd/shortener/main.go

exe:
	./cmd/shortener/shortener

re:
	killall -9 shortener || true
	rm -f cmd/shortener/shortener
	go build -o cmd/shortener/shortener cmd/shortener/main.go

test: test1 test2 test3 test4 test5

test1:
	./shortenertestbeta -test.v -test.run=^TestIteration1$$ -binary-path=cmd/shortener/shortener

test2:
	./shortenertestbeta -test.v -test.run=^TestIteration2$$ -source-path=.

test3:
	./shortenertestbeta -test.v -test.run=^TestIteration3$$ -source-path=.

test4:
	./shortenertestbeta -test.v -test.run=^TestIteration4$$ -binary-path=cmd/shortener/shortener -server-port=8080

test5:
	SERVER_PORT=8080 ./shortenertestbeta -test.v -test.run=^TestIteration5$$ -binary-path=cmd/shortener/shortener -server-port=8080

PHONY: run exe re test test1 test2 test3 test4 test5
