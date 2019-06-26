package graphqlapi

import (
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func MarshalGraphQLTimestampScalar(v timestamp.Timestamp) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		t, err := ptypes.Timestamp(&v)
		if err != nil {
			panic(err)
		}
		_, _ = io.WriteString(w, strconv.Quote(t.Format(time.RFC3339Nano)))
	})
}

func UnmarshalGraphQLTimestampScalar(v interface{}) (timestamp.Timestamp, error) {
	if tmpStr, ok := v.(string); ok {
		t, err := time.Parse(time.RFC3339Nano, tmpStr)
		if err != nil {
			return timestamp.Timestamp{}, err
		}
		v, err := ptypes.TimestampProto(t)
		if err != nil {
			return timestamp.Timestamp{}, err
		}
		return *v, nil
	}
	return timestamp.Timestamp{}, errors.New("time should be RFC3339 formatted string")
}
