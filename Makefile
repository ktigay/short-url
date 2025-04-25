SERVER_PORT=12345

build:
	go build -o ./cmd/shortener/shortener ./cmd/shortener/*.go

run-test: \
	build \
	run-test-a \
	run-test-u \
	run-test-s

run-test-u:
	go test ./...

run-test-s:
	go vet -vettool=$$(which statictest) ./...

run-test-a: \
	run-test-a1 \

run-test-a1:
	shortenertestbeta -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener

update-tpl:
	# git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
	git fetch template && git checkout template/main .github