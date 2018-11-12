build:
	go-bindata ui/**
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static" -s' -v -o openpassage-server .
run: build
	./openpassage-server
