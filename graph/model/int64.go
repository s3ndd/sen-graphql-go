package model

import (
	"errors"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalInt64(i int64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		b := strconv.AppendInt(nil, i, 10)
		_, _ = w.Write(b)
	})
}

func UnmarshalInt64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return i, nil
		}
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	}
	return 0, errors.New("not Int64")
}
