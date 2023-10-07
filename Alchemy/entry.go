package alchemy

import (
	"syscall/js"

	"github.com/gowebapi/webapi/graphics/webgl"
	"github.com/gowebapi/webapi/graphics/webgl2"
)

var app_url string

type App struct {
	Width    int
	Height   int
	Title    string
	OnStart  func()
	OnUpdate func(float64)
	OnDraw   func()
	OnEvent  func(*AppEvent)
}

// Used to make the update function only available in the local App struct, to the whole file
var tempStart func()
var tempUpdate func(float64)
var tempDraw func()

/*
USING THE EventFunc[T] type ------- (1)

	var custom_func EventFunc[string] = func(x ...string) {
		fmt.Printf(x[0] + "\n")
	}

USING THE EventFunc[T] type ------- (1)
*/

// *** Declaring an AlchemyEvent[T] *** var event AlchemyEvent[int] ------- (2)

var currentWidth, currentHeight int
var canvas js.Value
var glRef webgl2.RenderingContext
var appRef *App

var Cam Camera2D
var Shapes ShapeBatch

func Run(_app *App) {
	appRef = _app

	matrix4x4One = NewMatrix(4, 4)
	matrix4x4One.SetAll(1.0)

	js.Global().Get("document").Set("title", _app.Title)

	canvas = js.Global().Get("document").Call("getElementById", "viewport")

	canvasContext = canvas.Call("getContext", "webgl2")
	if canvasContext.IsNull() {
		LogF("CANVAS: Failed to Get Context")
	}

	canvas.Set("width", _app.Width)
	canvas.Set("height", _app.Height)

	glRef = *webgl2.RenderingContextFromJS(canvasContext)
	/*
		//vao := glRef.CreateVertexArray()
		//glRef.BindVertexArray(vao)
		vao := canvasContext.Call("createVertexArray")
		canvasContext.Call("bindVertexArray", vao)

		//vbo := glRef.CreateBuffer()
		//glRef.BindBuffer(webgl2.ARRAY_BUFFER, vbo)

		verts := []float32{
			0.0, 0.5,
			0.5, -0.5,
			-0.5, -0.5,
		}
		//vertsData := jsconv.Float32ToJs(verts)

		jsVerts := js.Global().Get("Uint8Array").New(len(verts) * 4)
		var verticesBytes []byte
		header := (*reflect.SliceHeader)(unsafe.Pointer(&verticesBytes))
		header.Cap = cap(verts) * 4
		header.Len = len(verts) * 4
		header.Data = uintptr(unsafe.Pointer(&verts[0]))

		js.CopyBytesToJS(jsVerts, verticesBytes)

		//glRef.BufferData2(webgl2.ARRAY_BUFFER, webgl2.UnionFromJS(vertsData), webgl2.DYNAMIC_DRAW)

		//ibo := glRef.CreateBuffer()
		//glRef.BindBuffer(webgl2.ELEMENT_ARRAY_BUFFER, ibo)

		elements := []uint32{
			0, 1, 2, 0,
		}
		//elemData := jsconv.UInt32ToJs(elements)

		//glRef.BufferData2(webgl2.ELEMENT_ARRAY_BUFFER, webgl2.UnionFromJS(elemData), webgl2.DYNAMIC_DRAW)

		jsElem := js.Global().Get("Uint8Array").New(len(elements) * 4)
		var elementsBytes []byte
		headerElem := (*reflect.SliceHeader)(unsafe.Pointer(&elementsBytes))
		headerElem.Cap = cap(elements) * 4
		headerElem.Len = len(elements) * 4
		headerElem.Data = uintptr(unsafe.Pointer(&elements[0]))

		js.CopyBytesToJS(jsElem, elementsBytes)

		vbo := canvasContext.Call("createBuffer")
		canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), vbo)
		canvasContext.Call("bufferData", canvasContext.Get("ARRAY_BUFFER"), jsVerts, canvasContext.Get("STATIC_DRAW"))

		ibo := canvasContext.Call("createBuffer")
		canvasContext.Call("bindBuffer", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), ibo)
		canvasContext.Call("bufferData", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), jsElem, canvasContext.Get("STATIC_DRAW"))

		canvasContext.Call("vertexAttribPointer", 0, 2, canvasContext.Get("FLOAT"), false, 0, 0)
		canvasContext.Call("enableVertexAttribArray", 0)
		//glRef.EnableVertexAttribArray(0)
		//glRef.VertexAttribPointer(0, 3, webgl.FLOAT, false, 0, 0)

		// Enable the attribute

		glRef.BindVertexArray(&webgl2.VertexArrayObject{})
		glRef.BindBuffer(webgl2.ARRAY_BUFFER, &webgl.Buffer{})
		glRef.BindBuffer(webgl2.ELEMENT_ARRAY_BUFFER, &webgl.Buffer{})

		setupShaders(&glRef)

		//glRef.BindVertexArray(vao)
		canvasContext.Call("bindVertexArray", vao)
		//canvasContext.Call("bindBuffer", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), ibo)
		// Point an attribute to the currently bound VBO

		glRef.ClearColor(0.5, 0.5, 0.5, 0.9)
		glRef.Clear(webgl2.COLOR_BUFFER_BIT)

		glRef.DrawElements(webgl.LINES, len(elements), webgl.UNSIGNED_INT, 0)

		   	vertexShaderCode := `
		       attribute vec2 coordinates;

		       void main(void) {
		           gl_Position = vec4(coordinates,0.0, 1.0);
		       }`

		   	vertexShader := canvasContext.Call("createShader", canvasContext.Get("VERTEX_SHADER"))
		   	canvasContext.Call("shaderSource", vertexShader, vertexShaderCode)
		   	canvasContext.Call("compileShader", vertexShader)

		   	fragmentShaderCode := `
		       void main(void) {
		           gl_FragColor = vec4(0.0, 1.0, 0.0, 0.1);
		       }`

		   	fragmentShader := canvasContext.Call("createShader", canvasContext.Get("FRAGMENT_SHADER"))
		   	canvasContext.Call("shaderSource", fragmentShader, fragmentShaderCode)
		   	canvasContext.Call("compileShader", fragmentShader)

		   	shaderProgram := canvasContext.Call("createProgram")
		   	canvasContext.Call("attachShader", shaderProgram, vertexShader)
		   	canvasContext.Call("attachShader", shaderProgram, fragmentShader)
		   	canvasContext.Call("linkProgram", shaderProgram)
		   	canvasContext.Call("useProgram", shaderProgram)

		   	canvasContext.Call("bindVertexArray", vao)
		   	//canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), vbo)
		   	//canvasContext.Call("bindBuffer", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), ibo)

		   	glRef.Viewport(0, 0, _app.Width, _app.Height)

		   	canvasContext.Call("clearColor", 1.0, 0.5, 0.5, 0.9)
		   	canvasContext.Call("clear", canvasContext.Get("COLOR_BUFFER_BIT"))

		   	//canvasContext.Call("drawArrays", canvasContext.Get("TRIANGLES"), 0, 3)
		   	canvasContext.Call("drawElements", canvasContext.Get("TRIANGLES"), len(elements), canvasContext.Get("UNSIGNED_INT"), 0)
	*/

	tempStart = _app.OnStart
	tempUpdate = _app.OnUpdate
	tempDraw = _app.OnDraw

	InitInputs()

	app_url = js.Global().Get("location").Get("href").String()
	js.Global().Set("js_start", js.FuncOf(JSStart))
	js.Global().Set("js_update", js.FuncOf(JSUpdate))
	js.Global().Set("js_draw", js.FuncOf(JSDraw))

	js.Global().Set("js_dpad_up", js.FuncOf(JSDpadUp))
	js.Global().Set("js_dpad_down", js.FuncOf(JSDpadDown))
	js.Global().Set("js_dpad_left", js.FuncOf(JSDpadLeft))
	js.Global().Set("js_dpad_right", js.FuncOf(JSDpadRight))
	js.Global().Set("js_main_button", js.FuncOf(JSMainButton))
	js.Global().Set("js_side_button", js.FuncOf(JSSideButton))

	// if I put it above the "js_start" then it would take a lot of time to run
	Cam.Init(*_app)
	Shapes.Init(_app)
	Cam.Update(*_app)

	addEventListenerWindow(JS_KEYUP, func(ae *AppEvent) {
		_app.OnEvent(ae)
	})
	addEventListenerWindow(JS_KEYDOWN, func(ae *AppEvent) {
		_app.OnEvent(ae)
	})
	addEventListenerWindow(JS_MOUSEDOWN, func(ae *AppEvent) {
		_app.OnEvent(ae)
	})
	addEventListenerWindow(JS_MOUSEUP, func(ae *AppEvent) {
		_app.OnEvent(ae)
	})
	addEventListenerWindow(JS_MOUSEMOVED, func(ae *AppEvent) {
		_app.OnEvent(ae)
	})

	/*
		Using the AlchemyEvent[T] ------- (2)

			event.AddListener(print_s)
			event.Invoke(1, 20003)
			event.RemoveListener(print_s)
	*/

	select {}
	//UnuseShader()
	//custom_func("STRING") ------- (1)

}

/*
--------- (2)

	func print_s(s ...int) {
		fmt.Println(s[1])
	}

--------- (2)
*/

func JSStart(this js.Value, inputs []js.Value) interface{} {
	tempStart()
	return nil
}

func JSUpdate(this js.Value, inputs []js.Value) interface{} {
	currentWidth = canvas.Get("width").Int()
	currentHeight = canvas.Get("height").Int()
	tempUpdate(inputs[0].Float())
	updateInput()
	Cam.Update(*appRef)
	return nil
}

var Red float64 = 0.0

func JSDraw(this js.Value, inputs []js.Value) interface{} {
	glRef.Viewport(0, 0, currentWidth, currentHeight)
	glRef.ClearColor(float32(Red), 0.0, 0.0, 1.0)
	glRef.Clear(webgl2.COLOR_BUFFER_BIT)

	//Shapes.DrawLine(NewVector2f(0.0, 0.0), NewVector2f(2.5, 0.5), RGBA8{255, 255, 0, 255})
	tempDraw()
	Shapes.Render(&Cam)
	return nil
}

func setupShaders(gl *webgl2.RenderingContext) *webgl.Program {
	// Vertex shader source code
	vertCode := `
	attribute vec3 coordinates;
	attribute vec4 colors;
	void main(void) {
		gl_Position = vec4(coordinates, 1.0);
	}`

	// Create a vertex shader object
	vShader := gl.CreateShader(webgl.VERTEX_SHADER)

	// Attach vertex shader source code
	gl.ShaderSource(vShader, vertCode)

	// Compile the vertex shader
	gl.CompileShader(vShader)

	//fragment shader source code
	fragCode := `
	void main(void) {
		gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
	}`

	// Create fragment shader object
	fShader := gl.CreateShader(webgl.FRAGMENT_SHADER)

	// Attach fragment shader source code
	gl.ShaderSource(fShader, fragCode)

	// Compile the fragmentt shader
	gl.CompileShader(fShader)

	// Create a shader program object to store
	// the combined shader program
	prog := gl.CreateProgram()

	// Attach a vertex shader
	gl.AttachShader(prog, vShader)

	// Attach a fragment shader
	gl.AttachShader(prog, fShader)

	// Link both the programs
	gl.LinkProgram(prog)

	// Use the combined shader program object
	gl.UseProgram(prog)

	return prog
}
