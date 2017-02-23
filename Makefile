.PHONY: build
build:
	go get ./...
	for arch in amd64 ; do \
		for os in darwin linux windows; do \
			GOOS=$$os GOARCH=$$arch go build -o docker-machine-driver-glesys_$$os-$$arch ./cmd/docker-machine-driver-glesys ; \
		done \
	done
	GOOS=windows GOARCH=386 go build -o docker-machine-driver-glesys_windows-i386 ./cmd/docker-machine-driver-glesys ; \
