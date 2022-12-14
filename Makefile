VERSION := 0.0.1
SHORTVERSION := 0.0.1

.PHONY: mac windows pi linux

windows:
	GOOS=windows GOARCH=386 go build -ldflags="-w -s -H windowsgui" -a -o bin/godeliver.exe main.go

linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -linkmode internal" -a -o bin/godeliver main.go

pi:
	GOOS=linux GOARCH=arm go build -ldflags="-w -s -linkmode internal" -a -o bin/godeliver_pi main.go

mac:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-w" -a -o bin/godeliver_mac main.go

all: linux pi mac

clean:
	rm -rf bin/godeliver*
