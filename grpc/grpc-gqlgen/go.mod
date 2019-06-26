module github.com/vvakame/til/grpc/grpc-gqlgen

go 1.12

replace github.com/vektah/gqlparser => github.com/vvakame/gqlparser v0.0.0-20190614064228-62f7407202a0

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.12.2
	github.com/99designs/gqlgen v0.9.0
	github.com/akutz/memconn v0.1.0
	github.com/favclip/golidator v2.1.1+incompatible // indirect
	github.com/favclip/ucon v2.2.1+incompatible
	github.com/golang/protobuf v1.3.1
	github.com/google/uuid v1.1.1
	github.com/google/wire v0.3.0
	github.com/grpc-ecosystem/grpc-gateway v1.9.2
	github.com/jhump/protoreflect v1.4.2
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/rakyll/statik v0.1.6
	github.com/vektah/gqlparser v1.1.2
	github.com/vvakame/sdlog v0.0.0-20190523062053-be70263e9c6c
	go.opencensus.io v0.22.0
	golang.org/x/xerrors v0.0.0-20190513163551-3ee3066db522
	google.golang.org/genproto v0.0.0-20190620144150-6af8c5fc6601
	google.golang.org/grpc v1.21.1
)
