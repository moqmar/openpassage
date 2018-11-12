build:
	cd ui && parcel build index.html
	go-bindata ui/dist/**
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static" -s' -v -o openpassage-server .
run: build
	./openpassage-server
