package graphqlapi

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vvakame/til/grpc/grpc-gqlgen/todopb"
)

func MarshalListADoneFilter(v todopb.ListARequest_DoneFilter) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(v.String()))
	})
}

func UnmarshalListADoneFilter(v interface{}) (todopb.ListARequest_DoneFilter, error) {
	if tmpStr, ok := v.(string); ok {
		v, ok := todopb.ListARequest_DoneFilter_value[tmpStr]
		if !ok {
			return 0, fmt.Errorf("invalid value format]: %s", tmpStr)
		}
		return todopb.ListARequest_DoneFilter(v), nil
	}
	return 0, fmt.Errorf("unexpected value type: %T", v)
}

func MarshalListBDoneFilter(v todopb.ListBRequest_DoneFilter) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(v.String()))
	})
}

func UnmarshalListBDoneFilter(v interface{}) (todopb.ListBRequest_DoneFilter, error) {
	if tmpStr, ok := v.(string); ok {
		v, ok := todopb.ListBRequest_DoneFilter_value[tmpStr]
		if !ok {
			return 0, fmt.Errorf("invalid value format]: %s", tmpStr)
		}
		return todopb.ListBRequest_DoneFilter(v), nil
	}
	return 0, fmt.Errorf("unexpected value type: %T", v)
}
