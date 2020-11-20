compile-gateway:
	CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a -ldflags '-extldflags "-static"' -o ../bin/jalapeno-gateway ./jalapeno-gateway.go
