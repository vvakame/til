package graphqlapi

import (
	"encoding/json"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func MarshalGraphQLTimeScalar(v timestamp.Timestamp) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		b, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		_, _ = w.Write(b)
	})
}

func UnmarshalGraphQLTimeScalar(v interface{}) (timestamp.Timestamp, error) {
	panic("not implemented")
}
