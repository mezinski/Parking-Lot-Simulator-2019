build:
	go build -o parkinglot

release:
	GOOS=darwin GOARCH=amd64 go build -o releases/parkinglot-osx-amd64
	GOOS=linux GOARCH=amd64 go build -o releases/parkinglot-linux-amd64
	GOOS=windows GOARCH=amd64 go build -o releases/parkinglot-windows-amd64