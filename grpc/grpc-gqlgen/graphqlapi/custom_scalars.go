package graphqlapi

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/genproto/googleapis/type/date"
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

func MarshalGraphQLDateScalar(v date.Date) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(fmt.Sprintf("%s/%s/%s", v.GetYear(), v.GetMonth(), v.GetDay())))
	})
}

func UnmarshalGraphQLDateScalar(v interface{}) (date.Date, error) {
	if tmpStr, ok := v.(string); ok {
		ss := strings.SplitN(tmpStr, "/", 3)
		if len(ss) != 3 {
			return date.Date{}, errors.New("date should be YYYY/MM/DD formatted string")
		}

		y, err := strconv.Atoi(ss[0])
		if err != nil {
			return date.Date{}, errors.New("date should be YYYY/MM/DD formatted string")
		}
		m, err := strconv.Atoi(ss[1])
		if err != nil {
			return date.Date{}, errors.New("date should be YYYY/MM/DD formatted string")
		}
		d, err := strconv.Atoi(ss[2])
		if err != nil {
			return date.Date{}, errors.New("date should be YYYY/MM/DD formatted string")
		}

		return date.Date{
			Year:  int32(y),
			Month: int32(m),
			Day:   int32(d),
		}, nil
	}
	return date.Date{}, errors.New("date should be YYYY/MM/DD formatted string")
}

func MarshalGraphQLUInt64Scalar(v uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(fmt.Sprintf("%d", v)))
	})
}

func UnmarshalGraphQLUInt64Scalar(v interface{}) (uint64, error) {
	if tmpStr, ok := v.(string); ok {
		v, err := strconv.ParseUint(tmpStr, 10, 64)
		if err != nil {
			return 0, errors.New("invalid uint64 format")
		}
		return v, nil
	}
	return 0, errors.New("invalid uint64 format")
}

func MarshalGraphQLInt64Scalar(v int64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.Quote(fmt.Sprintf("%d", v)))
	})
}

func UnmarshalGraphQLInt64Scalar(v interface{}) (int64, error) {
	if tmpStr, ok := v.(string); ok {
		v, err := strconv.ParseInt(tmpStr, 10, 64)
		if err != nil {
			return 0, errors.New("invalid uint64 format")
		}
		return v, nil
	}
	return 0, errors.New("invalid uint64 format")
}
