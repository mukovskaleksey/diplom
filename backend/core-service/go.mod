module core-service

go 1.24.7

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/lib/pq v1.12.3
	golang.org/x/crypto v0.47.0
	google.golang.org/grpc v1.80.0
	gopkg.in/yaml.v3 v3.0.1
	proto v0.0.0
)

require (
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace proto => ../proto
