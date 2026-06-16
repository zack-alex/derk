//go:build wasm

package main

import (
	"encoding/json"
	"github.com/zack-alex/derk"
	"syscall/js"
)

func deriveAndFormat(this js.Value, p []js.Value) any {
	masterPassword := p[0].String()
	jsonSpec := p[1].String()
	var spec map[string]string
	err := json.Unmarshal([]byte(jsonSpec), &spec)
	if err != nil {
		return err.Error()
	}

	result, err := derk.DeriveAndFormat(masterPassword, spec)
	if err != nil {
		return err.Error()
	}
	return js.ValueOf(result)
}

func main() {
	js.Global().Set("deriveAndFormat", js.FuncOf(deriveAndFormat))
	<-make(chan struct{})
}
