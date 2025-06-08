package model

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalTimestamp(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		fmt.Fprintf(w, `"%s"`, t.Format(time.RFC3339Nano))
	})
}

func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	if s, ok := v.(string); ok {
		return time.Parse(time.RFC3339Nano, s)
	}
	return time.Time{}, errors.New("timestamp format error")
}
