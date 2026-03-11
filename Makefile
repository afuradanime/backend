run:
	go run ./cmd/api

build:
	go build -o bin/afuradanime-api ./cmd/api

run-build:
	@if [ ! -f ./bin/afuradanime-api ]; then \
		$(MAKE) build; \
	fi
	./bin/afuradanime-api

clean:
	rm -rf bin/ && rm -rf _openapi/