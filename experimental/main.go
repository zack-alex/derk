// main.go
package main

import (
	"syscall/js"
)

func invertText(this js.Value, p []js.Value) interface{} {
	input := p[0].String()
	result := reverse(input)
	return js.ValueOf(result)
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r) + "!!!"
}

func main() {
	js.Global().Set("invertText", js.FuncOf(invertText))
	<-make(chan struct{})
}
