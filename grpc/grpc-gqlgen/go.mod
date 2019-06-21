module github.com/vvakame/til/grpc/grpc-gqlgen

go 1.12

require (
	cloud.google.com/go v0.40.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.12.1
	github.com/99designs/gqlgen v0.9.0
	github.com/agnivade/levenshtein v1.0.2 // indirect
	github.com/akutz/memconn v0.1.0
	github.com/aws/aws-sdk-go v1.20.0 // indirect
	github.com/favclip/golidator v2.1.1+incompatible // indirect
	github.com/favclip/ucon v2.2.1+incompatible
	github.com/golang/protobuf v1.3.1
	github.com/google/uuid v1.1.1
	github.com/google/wire v0.2.2
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.1
	github.com/jhump/protoreflect v1.4.2
	github.com/rakyll/statik v0.1.6
	github.com/vektah/gqlparser v1.1.2
	github.com/vvakame/sdlog v0.0.0-20190523062053-be70263e9c6c
	go.opencensus.io v0.22.0
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980 // indirect
	golang.org/x/sys v0.0.0-20190613124609-5ed2794edfdc // indirect
	golang.org/x/tools v0.0.0-20190613204242-ed0dc450797f // indirect
	golang.org/x/xerrors v0.0.0-20190513163551-3ee3066db522
	google.golang.org/appengine v1.6.1 // indirect
	google.golang.org/genproto v0.0.0-20190611190212-a7e196e89fd3
	google.golang.org/grpc v1.21.1
)

replace github.com/vektah/gqlparser => github.com/vvakame/gqlparser v0.0.0-20190614064228-62f7407202a0
