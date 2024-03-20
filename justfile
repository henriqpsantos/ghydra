@deploy: tidy build
	mv "./ghydra" ~/.local/bin

@run: tidy
	go run .

@tidy:
	go mod tidy

@build: tidy
	go build .

