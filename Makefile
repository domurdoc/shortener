BIN = cmd/shortener/shortener
PORT = 8080
DSN = postgresql://domurdoc@localhost:5432/test?sslmode=disable
FILE = db.json
DIR = .
MAIN = cmd/shortener/main.go
LOAD = cmd/loadtest/main.go
PPROF_PORT = 8888
PROFILE_FILE = profiles/result.pprof
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

load:
	go run ${LOAD}

prof:
	curl -s -v -o ${PROFILE_FILE} http://localhost:${PPROF_PORT}/debug/pprof/heap?seconds=30

showprof:
	go tool pprof -http=":9090" -seconds=30 ${PROFILE_FILE}

test: re test1 test2 test3 test4 test5 test6 test7 test8 test9 test10 test11 test12 test13 test14 test15 test16 test17 test18

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

test16: kill re
	./${TESTBIN} -test.v -test.run=^TestIteration16$$ -source-path=.

test17: kill re
	./${TESTBIN} -test.v -test.run=^TestIteration17$$ -source-path=.

test18: kill re
	./${TESTBIN} -test.v -test.run=^TestIteration18$$ -source-path=.

PHONY: run exe re kill m mm md load prof showprof test test1 test2 test3 test4 test5 test6 test7 test8 test9 test10 test11 test12 test13 test14 test15 test16 test17 test18
