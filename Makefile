BIN = cmd/shortener/shortener
PORT = 8080
DSN = postgresql://domurdoc@localhost:5432/test?sslmode=disable
FILE = db.json
DIR = .
MAIN = cmd/shortener/main.go
MNAME = unnamed
TESTBIN = shortenertest
WIPEDBBIN = wipedb

run:
	go run ${MAIN} -d ${DSN}

exe:
	./${BIN}

re:
	rm -f ${BIN}
	go build -o ${BIN} ${MAIN}

kill:
	killall -9 shortener || true

mm:
	migrate create -ext sql -dir ./migrations -seq ${MNAME}

m:
	migrate -database "${DSN}" -path ./migrations up

md:
	migrate -database "${DSN}" -path ./migrations down 1

test: re test1 test2 test3 test4 test5 test6 test7 test8 test9 test10 test11 test12 test13 test14

test1: kill
	./${TESTBIN} -test.v -test.run=^TestIteration1$$ -binary-path=${BIN}

test2: kill
	./${TESTBIN} -test.v -test.run=^TestIteration2$$ -source-path=${DIR}

test3: kill
	./${TESTBIN} -test.v -test.run=^TestIteration3$$ -source-path=${DIR}

test4: kill
	./${TESTBIN} -test.v -test.run=^TestIteration4$$ -binary-path=${BIN} -server-port=${PORT}

test5: kill
	./${TESTBIN} -test.v -test.run=^TestIteration5$$ -binary-path=${BIN} -server-port=${PORT}

test6: kill
	./${TESTBIN} -test.v -test.run=^TestIteration6$$ -source-path=${DIR}

test7: kill
	./${TESTBIN} -test.v -test.run=^TestIteration7$$ -binary-path=${BIN} -source-path=${DIR}

test8: kill
	./${TESTBIN} -test.v -test.run=^TestIteration8$$ -binary-path=${BIN}

test9: kill
	./${TESTBIN} -test.v -test.run=^TestIteration9$$ -binary-path=${BIN} -source-path=${DIR} -file-storage-path=${FILE}

test10: kill
	./${TESTBIN} -test.v -test.run=^TestIteration10$$ -binary-path=${BIN} -source-path=${DIR} -database-dsn=${DSN}

test11: kill
	./${TESTBIN} -test.v -test.run=^TestIteration11$$ -binary-path=${BIN} -database-dsn=${DSN}

test12: kill
	./${TESTBIN} -test.v -test.run=^TestIteration12$$ -binary-path=${BIN} -database-dsn=${DSN}

test13: kill
	./${TESTBIN} -test.v -test.run=^TestIteration13$$ -binary-path=${BIN} -database-dsn=${DSN}

test14: kill re
	./${TESTBIN} -test.v -test.run=^TestIteration14$$ -binary-path=${BIN} -database-dsn=${DSN}

test15: kill re
	./${TESTBIN} -test.v -test.run=^TestIteration15$$ -binary-path=${BIN} -database-dsn=${DSN}

PHONY: run exe re kill m mm md test test1 test2 test3 test4 test5 test6 test7 test8 test9 test10 test11 test12 test13 test14 test15
