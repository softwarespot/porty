package logging

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONLogger struct {
	enc *json.Encoder
}

func NewJSONLogger() *JSONLogger {
	return &JSONLogger{
		enc: json.NewEncoder(os.Stdout),
	}
}

func (jl *JSONLogger) Log(v any) {
	if err := jl.enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "skipped logging entry due to error: %+v\n", err)
	}
}
