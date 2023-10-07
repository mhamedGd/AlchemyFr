package alchemy

import "syscall/js"

var DPadUp_Pressed AlchemyEvent[int]
var DPadUp_Released AlchemyEvent[int]

var DPadDown_Pressed AlchemyEvent[int]
var DPadDown_Released AlchemyEvent[int]

var DPadLeft_Pressed AlchemyEvent[int]
var DPadLeft_Released AlchemyEvent[int]

var DPadRight_Pressed AlchemyEvent[int]
var DPadRight_Released AlchemyEvent[int]

var MainButton_Pressed AlchemyEvent[int]
var MainButton_Released AlchemyEvent[int]

var SideButton_Pressed AlchemyEvent[int]
var SideButton_Released AlchemyEvent[int]

// ---------------
func JSDpadUp(this js.Value, inputs []js.Value) interface{} {
	if inputs[0].Int() == 1 {
		DPadUp_Pressed.Invoke(0)
	} else {
		DPadUp_Released.Invoke(0)
	}
	return nil
}

// ---------------
func JSDpadDown(this js.Value, inputs []js.Value) interface{} {
	if inputs[0].Int() == 1 {
		DPadDown_Pressed.Invoke(0)
	} else {
		DPadDown_Released.Invoke(0)
	}
	return nil
}

// ---------------
func JSDpadLeft(this js.Value, inputs []js.Value) interface{} {
	if inputs[0].Int() == 1 {
		DPadLeft_Pressed.Invoke(0)
	} else {
		DPadLeft_Released.Invoke(0)
	}
	return nil
}

// ---------------
func JSDpadRight(this js.Value, inputs []js.Value) interface{} {
	if inputs[0].Int() == 1 {
		DPadRight_Pressed.Invoke(0)
	} else {
		DPadRight_Released.Invoke(0)
	}
	return nil
}

// ---------------
func JSMainButton(this js.Value, inputs []js.Value) interface{} {
	if inputs[0].Int() == 1 {
		MainButton_Pressed.Invoke(0)
	} else {
		MainButton_Released.Invoke(0)
	}
	return nil
}

// ---------------
func JSSideButton(this js.Value, inputs []js.Value) interface{} {
	if inputs[0].Int() == 1 {
		SideButton_Pressed.Invoke(0)
	} else {
		SideButton_Released.Invoke(0)
	}
	return nil
}
