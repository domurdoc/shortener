run:
	go run ./cmd/shortener/main.go

exe:
	./cmd/shortener/shortener

re:
	rm -f cmd/shortener/shortener
	go build -o cmd/shortener/shortener cmd/shortener/main.go

kill:
	killall -9 shortener || true

test: re test1 test2 test3 test4 test5 test6 test7 test8 test9 test10 test11 test12

test1: kill
	./shortenertest -test.v -test.run=^TestIteration1$$ -binary-path=cmd/shortener/shortener

test2: kill
	./shortenertest -test.v -test.run=^TestIteration2$$ -source-path=.

test3: kill
	./shortenertest -test.v -test.run=^TestIteration3$$ -source-path=.

test4: kill
	./shortenertest -test.v -test.run=^TestIteration4$$ -binary-path=cmd/shortener/shortener -server-port=8080

test5: kill
	SERVER_PORT=8080 ./shortenertest -test.v -test.run=^TestIteration5$$ -binary-path=cmd/shortener/shortener -server-port=8080

test6: kill
	./shortenertest -test.v -test.run=^TestIteration6$$ -source-path=.

test7: kill
	./shortenertest -test.v -test.run=^TestIteration7$$ -binary-path=cmd/shortener/shortener -source-path=.

test8: kill
	./shortenertest -test.v -test.run=^TestIteration8$$ -binary-path=cmd/shortener/shortener

test9: kill
	./shortenertest -test.v -test.run=^TestIteration9$$ -binary-path=cmd/shortener/shortener -source-path=. -file-storage-path=db.json

test10: kill
	./shortenertest -test.v -test.run=^TestIteration10$$ -binary-path=cmd/shortener/shortener -source-path=. -database-dsn=postgresql://domurdoc@localhost:5432/test

test11: kill
	./shortenertest -test.v -test.run=^TestIteration11$$ -binary-path=cmd/shortener/shortener -database-dsn=postgresql://domurdoc@localhost:5432/test

test12: kill
	./shortenertest -test.v -test.run=^TestIteration12$$ -binary-path=cmd/shortener/shortener -database-dsn=postgresql://domurdoc@localhost:5432/test

PHONY: run exe re test test1 test2 test3 test4 test5 test6 test7 test8 test9 test10 test11 test12
