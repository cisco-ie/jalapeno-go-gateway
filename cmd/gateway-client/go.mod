module gateway-client

go 1.13

require (
	github.com/cisco-ie/jalapeno-go-gateway/pkg/apis v0.0.0-00010101000000-000000000000 // indirect
	google.golang.org/grpc v1.27.0 // indirect
)

replace github.com/cisco-ie/jalapeno-go-gateway/pkg/apis => ../../pkg/apis
