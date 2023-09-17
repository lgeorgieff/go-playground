package main

import (
	"context"
	"os"

	"github.com/pkg/errors"
)

type contextKey string

const streamKey contextKey = "stream_key"

func sayHello(ctx context.Context, name string) error {
	value := ctx.Value(streamKey)
	if value == nil {
		return errors.New("missing target file")
	}
	f, ok := value.(*os.File)
	if !ok {
		return errors.Errorf("%T target is not a file", value)
	}
	f.WriteString("Hello, " + name + "!\n")

	return nil
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, streamKey, os.Stdout)
	if err := sayHello(ctx, "John"); err != nil {
		panic(err)
	}
}
