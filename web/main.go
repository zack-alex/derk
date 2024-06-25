//go:build wasm

package main

import (
	"github.com/tekhnus/derk"
	"syscall/js"
)

func deriveAndFormat(this js.Value, p []js.Value) interface{} {
	s := p[0].String()
	result, err := derk.DeriveAndFormat(s, map[string]string{"method": "v1"})
	if err != nil {
		return err.Error()
	}
	return js.ValueOf(result)
}

func main() {
	js.Global().Set("deriveAndFormat", js.FuncOf(deriveAndFormat))
	<-make(chan struct{})
}
