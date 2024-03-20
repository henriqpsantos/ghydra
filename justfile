@install: tidy
	go install

@run: tidy
	go run .

@tidy:
	go mod tidy

@build: tidy
	go build .

