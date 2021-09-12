build:
	go build -v -ldflags="-X 'github.com/vzau/thoth/pkg/version.BuildTime=$(shell TZ='UTC' date)' -X 'github.com/vzau/thoth/pkg/version.GitCommit=$(shell git rev-parse --short HEAD)' -X 'github.com/vzau/thoth/pkg/version.Version=$(shell git describe --tags --abbrev=0 HEAD || echo dev)'" .

clean:
	rm -f thoth