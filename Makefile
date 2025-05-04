SERVER_PORT=12345

build:
	go build -o ./cmd/shortener/shortener ./cmd/shortener/*.go

run-test: \
	build \
	run-test-a \
	run-test-u \
	run-test-s \
	run-lint

run-test-u:
	go test ./...

run-test-s:
	go vet -vettool=$$(which statictest) ./...

run-test-a: \
	run-test-a1 \
	run-test-a2 \
	run-test-a3 \
	run-test-a4 \
	run-test-a5

run-test-a1:
	shortenertestbeta -test.v -test.run=^TestIteration1$$ -binary-path=cmd/shortener/shortener
run-test-a2:
	shortenertestbeta -test.v -test.run=^TestIteration2$$ -source-path=.
run-test-a3:
	shortenertestbeta -test.v -test.run=^TestIteration3$$ -source-path=.
run-test-a4:
	shortenertestbeta -test.v -test.run=^TestIteration4$$ -binary-path=cmd/shortener/shortener -server-port=$(SERVER_PORT)
run-test-a5:
	shortenertestbeta -test.v -test.run=^TestIteration5$$ -binary-path=cmd/shortener/shortener -server-port=$(SERVER_PORT)

update-tpl:
	# git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
	git fetch template && git checkout template/main .github

run-lint:
	golangci-lint run