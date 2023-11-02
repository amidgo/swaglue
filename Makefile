swaglue:
	go build -ldflags="-s -w" -o build/swaglue cmd/swaglue/main.go 
install:
	go build -ldflags="-s -w" -o ${GOPATH}/bin/swaglue cmd/swaglue/main.go 