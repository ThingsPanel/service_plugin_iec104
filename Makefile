init:
	go install github.com/mitchellh/gox@latest
windows:
	gox -osarch="windows/amd64" -ldflags "-s -w" -output=./bin/iec104-windows-amd64 watersql/cmd/server
linux:
	gox -osarch="linux/amd64" -ldflags "-s -w" -output=./bin/iec104-linux-amd64 watersql/cmd/server