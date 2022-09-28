package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
)

func withFile(output string, fn func(b *bytes.Buffer)) {
	out := new(bytes.Buffer)
	fmt.Printf("%s\n", output)
	fn(out)

	formatted, err := format.Source(out.Bytes())
	if err != nil {
		panic(fmt.Errorf("%s\n\n%w", out.String(), err))
	}

	f, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	must(err)
	defer func() {
		must(err)
	}()

	_, err = f.Write(formatted)
	must(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
