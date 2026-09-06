package logger

import (
	"net/http"

	"go.uber.org/zap/zapcore"
)

type ZapHeader http.Header

func (h *ZapHeader) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range *h {
		enc.AddArray(k, zapcore.ArrayMarshalerFunc(func(ae zapcore.ArrayEncoder) error {
			for _, vv := range v {
				ae.AppendString(vv)
			}
			return nil
		}))
	}
	return nil
}
