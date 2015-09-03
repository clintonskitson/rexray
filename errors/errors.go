package errors

import (
	"bytes"
	goerrors "errors"
	"fmt"
)

type Fields map[string]interface{}

func New(message string) error {
	return goerrors.New(message)
}

func Newf(format string, a ...interface{}) error {
	return goerrors.New(fmt.Sprintf(format, a))
}

func WithError(message string, inner error) error {
	return WithFieldsE(nil, message, inner)
}

func WithField(key string, val interface{}, message string) error {
	return WithFields(map[string]interface{}{key: val}, message)
}

func WithFieldE(key string, val interface{}, message string, inner error) error {
	return WithFieldsE(Fields{key: val}, message, inner)
}

func WithFields(fields map[string]interface{}, message string) error {
	return WithFieldsE(fields, message, nil)
}

func WithFieldsE(
	fields map[string]interface{}, message string, inner error) error {

	if fields == nil {
		fields = Fields{}
	}

	if inner != nil {
		fields["inner"] = inner
	}

	x := 0
	l := len(fields)

	var b bytes.Buffer
	for k, v := range fields {
		if x < l-1 {
			b.WriteString(fmt.Sprintf("%s=%v,", k, v))
		} else {
			b.WriteString(fmt.Sprintf("%s=%v", k, v))
		}
		x = x + 1
	}

	return Newf("%s %s", message, b.String())
}
