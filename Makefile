build:
	./set_version.sh
	go build ./cmd/pgsectest

debug:
	go build -gcflags "all=-N -l" ./cmd/pgsectest
	~/go/bin/dlv --headless --listen=:2345 --api-version=2 --accept-multiclient exec ./pgsectest ./testdata

run:
	./pgsectest -f tests.yaml

fmt:
	gofmt -w .

test: sec lint

sec:
	gosec ./...
lint:
	golangci-lint run
