@deploy: tidy build
	mv "./go-runner" ~/.local/bin

@run: tidy
	go run .

@tidy:
	go mod tidy

@build: tidy
	go build .

