package alchemy

import "syscall/js"

var DPadUp_Pressed AlchemyEvent[int]
var DPadUp_Released AlchemyEvent[int]

func JSDpadUp(this js.Value, inputs []js.Value) interface{} {
	if inputs[0].Int() == 1 {
		DPadUp_Pressed.Invoke(0)
	} else {
		DPadUp_Released.Invoke(0)
	}
	return nil
}
