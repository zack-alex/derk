//go:build wasm

package main

import (
	"github.com/zack-alex/derk"
	"syscall/js"
)

func deriveAndFormat(this js.Value, p []js.Value) interface{} {
	masterPassword := p[0].String()
	domain := p[1].String()
	username := p[2].String()
	method := p[3].String()

	result, err := derk.DeriveAndFormat(masterPassword, map[string]string{"domain": domain, "username": username, "method": method})
	if err != nil {
		return err.Error()
	}
	return js.ValueOf(result)
}

func main() {
	js.Global().Set("deriveAndFormat", js.FuncOf(deriveAndFormat))
	<-make(chan struct{})
}
