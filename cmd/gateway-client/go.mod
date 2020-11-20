module gateway-client

go 1.13

require (
	github.com/cisco-ie/jalapeno-go-gateway/pkg/apis v0.0.0-00010101000000-000000000000 // indirect
	github.com/cisco-ie/jalapeno-go-gateway/pkg/bgpclient v0.0.0-00010101000000-000000000000 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/osrg/gobgp v2.0.0+incompatible // indirect
	google.golang.org/grpc v1.27.0 // indirect
)

replace (
	github.com/cisco-ie/jalapeno-go-gateway/pkg/apis => ../../pkg/apis
	github.com/cisco-ie/jalapeno-go-gateway/pkg/bgpclient => ../../pkg/bgpclient
)
