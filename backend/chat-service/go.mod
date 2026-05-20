module chat-service

go 1.25.0

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/lib/pq v1.12.3
	google.golang.org/grpc v1.80.0
	gopkg.in/yaml.v3 v3.0.1
	proto v0.0.0
)

require (
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260406210006-6f92a3bedf2d // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace proto => ../proto
