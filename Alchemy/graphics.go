package alchemy

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
	"syscall/js"
	"unsafe"

	"github.com/gowebapi/webapi/graphics/webgl"
	"github.com/gowebapi/webapi/graphics/webgl2"
)

var canvasContext js.Value

func CreateVertexArray() webgl.Buffer {
	return *webgl.BufferFromJS(canvasContext.Call("createVertexArray"))

}

/*
	func BindVertexArray(buf webgl2.Buffer) {
		canvasContext.Call("createVertexArray", js.Value(buf))
	}

	func BindAttribLocation(program webgl2.Program, index int, name string) {
		canvasContext.Call("bindAttribLocation", js.Value(program), index, name)
	}
*/
type RGBA8 struct {
	r, g, b, a uint8
}

func NewRGBA8(r, g, b, a uint8) RGBA8 {
	return RGBA8{r, g, b, a}
}

type Vertex struct {
	Coordinates Vector2f
	Color       RGBA8
}

const vertexByteSize uintptr = unsafe.Sizeof(Vertex{})

type ShapeBatch struct {
	VAO              *webgl2.VertexArrayObject
	VBO, IBO         *webgl.Buffer
	Vertices         []Vertex
	Indices          []int32
	NumberOfElements int
	Shader           ShaderProgram
	LineWidth        float32
}

func (_shapesB *ShapeBatch) Init(_app *App) {
	_shapesB.LineWidth = 50

	_shapesB.Vertices = make([]Vertex, 0)
	_shapesB.Indices = make([]int32, 0)

	_shapesB.VAO = glRef.CreateVertexArray()
	glRef.BindVertexArray(_shapesB.VAO)

	_shapesB.VBO = glRef.CreateBuffer()
	glRef.BindBuffer(webgl.ARRAY_BUFFER, _shapesB.VBO)

	/*
		glRef.EnableVertexAttribArray(0)
		glRef.VertexAttribPointer(0, 2, webgl.FLOAT, false, int(unsafe.Sizeof(Vertex{})), int(unsafe.Offsetof(Vertex{}.Coordinates)))
			glRef.EnableVertexAttribArray(1)
			glRef.VertexAttribPointer(1, 4, webgl.UNSIGNED_BYTE, true, int(unsafe.Sizeof(Vertex{})), int(unsafe.Offsetof(Vertex{}.Color)))
	*/

	glRef.EnableVertexAttribArray(0)
	glRef.VertexAttribPointer(0, 2, webgl.FLOAT, false, 12, 0)
	glRef.EnableVertexAttribArray(1)
	glRef.VertexAttribPointer(1, 4, webgl.UNSIGNED_BYTE, true, 12, 8)

	_shapesB.IBO = glRef.CreateBuffer()
	glRef.BindBuffer(webgl.ELEMENT_ARRAY_BUFFER, _shapesB.IBO)

	//glRef.BindBuffer(webgl.ARRAY_BUFFER, nil)
	//glRef.BindBuffer(webgl.ELEMENT_ARRAY_BUFFER, nil)

	glRef.DisableVertexAttribArray(0)
	glRef.DisableVertexAttribArray(1)

	vertexShader := `#version 300 es

	precision mediump float;

	in vec2 coordinates;
	in vec4 colors;

	out vec4 vertex_FragColor;

	uniform mat4 projection_matrix;
	uniform mat4 view_matrix;

	void main(void) {
		vec4 global_position = vec4(0.0);
		global_position = view_matrix * vec4(coordinates, 0.0, 1.0);
		global_position.z = 0.0;
		global_position.w = 1.0;		
		gl_Position = global_position;
		
		
		vertex_FragColor = colors;
	}`

	fragmentShader := `#version 300 es

	precision mediump float;

	in vec4 vertex_FragColor;

	out vec4 fragColor;
	void main(void) {
		fragColor = vertex_FragColor;
	}`

	_shapesB.Shader.ParseShader(vertexShader, fragmentShader)
	//_shapesB.Shader.ParseShaderFromFile("shapes.shader")
	_shapesB.Shader.CreateShaderProgram()
	_shapesB.Shader.AddAttribute("coordinates")
	//_shapesB.Shader.AddAttribute("vertexColor")
	glRef.Viewport(0, 0, _app.Width, _app.Height)

	// LogF("%v", glRef.GetAttribLocation(_shapesB.Shader.ShaderProgramID, "vertexColor"))
}

func (_sp *ShapeBatch) DrawLine(_from, _to Vector2f, _color RGBA8) {
	//var offset Vector2f = NewVector2f(_sp.LineWidth/2.0, 0.0)
	var offset Vector2f

	offset = _to.Subtract(_from)
	offset = offset.Perpendicular()
	offset = offset.Normalize()
	offset = offset.Scale(_sp.LineWidth)

	vertsSize := len(_sp.Vertices)
	_sp.Vertices = append(_sp.Vertices, Vertex{Coordinates: _from.Subtract(offset), Color: _color})
	_sp.Vertices = append(_sp.Vertices, Vertex{Coordinates: _to.Subtract(offset), Color: _color})
	_sp.Vertices = append(_sp.Vertices, Vertex{Coordinates: _from.Add(offset), Color: _color})
	_sp.Vertices = append(_sp.Vertices, Vertex{Coordinates: _to.Add(offset), Color: _color})

	_sp.Indices = append(_sp.Indices, int32(vertsSize))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+1))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+2))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+2))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+1))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+3))
}

func (_sp *ShapeBatch) DrawRect(_center, _dimensions Vector2f, _color RGBA8) {
	offsetX := _dimensions.Scale(0.5).X
	offsetY := _dimensions.Scale(0.5).Y

	lineDifference := _sp.LineWidth

	_sp.DrawLine(NewVector2f(-offsetX-lineDifference, offsetY).Add(_center), NewVector2f(offsetX+lineDifference, offsetY).Add(_center), _color)
	_sp.DrawLine(NewVector2f(-offsetX-lineDifference, -offsetY).Add(_center), NewVector2f(offsetX+lineDifference, -offsetY).Add(_center), _color)
	_sp.DrawLine(NewVector2f(offsetX, offsetY+lineDifference).Add(_center), NewVector2f(offsetX, -offsetY-lineDifference).Add(_center), _color)
	_sp.DrawLine(NewVector2f(-offsetX, offsetY+lineDifference).Add(_center), NewVector2f(-offsetX, -offsetY-lineDifference).Add(_center), _color)
}

func (_sp *ShapeBatch) DrawTriangle(_center, _dimensions Vector2f, _color RGBA8) {
	numOfVertices := 3

	pos := [3]Vector2f{}
	for i := 0; i < numOfVertices; i++ {
		angle := (float32(i) / float32(3.0)) * 2.0 * PI
		pos[i] = NewVector2f(_center.X+float32(math.Cos(float64(angle)))*_dimensions.Y, _center.Y+float32(math.Sin(float64(angle)))*_dimensions.X)
	}

	for i := 0; i < numOfVertices-1; i++ {
		_sp.DrawLine(pos[i], pos[i+1], _color)
	}
	_sp.DrawLine(pos[numOfVertices-1], pos[0], _color)
}

func (_sp *ShapeBatch) DrawTriangleRotated(_center, _dimensions Vector2f, _color RGBA8, rotation float32) {
	numOfVertices := 3

	pos := [3]Vector2f{}
	for i := 0; i < numOfVertices; i++ {
		angle := (float32(i) / float32(3.0)) * 2.0 * PI
		pos[i] = NewVector2f(_center.X+float32(math.Cos(float64(angle)))*_dimensions.Y, _center.Y+float32(math.Sin(float64(angle)))*_dimensions.X)
		pos[i] = pos[i].Rotate(rotation, _center)
	}

	for i := 0; i < numOfVertices-1; i++ {
		_sp.DrawLine(pos[i], pos[i+1], _color)
	}
	_sp.DrawLine(pos[numOfVertices-1], pos[0], _color)
}

func (_sp *ShapeBatch) DrawFillRect(_center, _dimensions Vector2f, _color RGBA8) {
	offset := _dimensions.Scale(0.5)

	vertsSize := len(_sp.Vertices)
	_sp.Vertices = append(_sp.Vertices, Vertex{Coordinates: _center.Subtract(offset), Color: _color})
	_sp.Vertices = append(_sp.Vertices, Vertex{Coordinates: _center.SubtractXY(offset.X, -offset.Y), Color: _color})
	_sp.Vertices = append(_sp.Vertices, Vertex{Coordinates: _center.AddXY(offset.X, -offset.Y), Color: _color})
	_sp.Vertices = append(_sp.Vertices, Vertex{Coordinates: _center.Add(offset), Color: _color})

	_sp.Indices = append(_sp.Indices, int32(vertsSize))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+1))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+2))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+2))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+1))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+3))
}

func (_sp *ShapeBatch) DrawFillRectRotated(_center, _dimensions Vector2f, _color RGBA8, _rotation float32) {
	offset := _dimensions.Scale(0.5)

	vertsSize := len(_sp.Vertices)
	_sp.Vertices = append(_sp.Vertices, Vertex{Coordinates: _center.Subtract(offset).Rotate(_rotation, _center), Color: _color})
	_sp.Vertices = append(_sp.Vertices, Vertex{Coordinates: _center.SubtractXY(offset.X, -offset.Y).Rotate(_rotation, _center), Color: _color})
	_sp.Vertices = append(_sp.Vertices, Vertex{Coordinates: _center.AddXY(offset.X, -offset.Y).Rotate(_rotation, _center), Color: _color})
	_sp.Vertices = append(_sp.Vertices, Vertex{Coordinates: _center.Add(offset).Rotate(_rotation, _center), Color: _color})

	_sp.Indices = append(_sp.Indices, int32(vertsSize))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+1))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+2))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+2))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+1))
	_sp.Indices = append(_sp.Indices, int32(vertsSize+3))
}

func (_sp *ShapeBatch) finalize() {

	glRef.BindVertexArray(_sp.VAO)

	jsVerts := js.Global().Get("Uint8Array").New(len(_sp.Vertices) * 12)
	var verticesBytes []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&verticesBytes))
	header.Cap = cap(_sp.Vertices) * 12
	header.Len = len(_sp.Vertices) * 12
	header.Data = uintptr(unsafe.Pointer(&_sp.Vertices[0]))

	js.CopyBytesToJS(jsVerts, verticesBytes)

	jsElem := js.Global().Get("Uint8Array").New(len(_sp.Indices) * 4)
	var elementsBytes []byte
	headerElem := (*reflect.SliceHeader)(unsafe.Pointer(&elementsBytes))
	headerElem.Cap = cap(_sp.Indices) * 4
	headerElem.Len = len(_sp.Indices) * 4
	headerElem.Data = uintptr(unsafe.Pointer(&_sp.Indices[0]))

	js.CopyBytesToJS(jsElem, elementsBytes)

	canvasContext.Call("bindBuffer", canvasContext.Get("ARRAY_BUFFER"), _sp.VBO.Value_JS)
	canvasContext.Call("bufferData", canvasContext.Get("ARRAY_BUFFER"), jsVerts, canvasContext.Get("STATIC_DRAW"))

	canvasContext.Call("bindBuffer", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), _sp.IBO.Value_JS)
	canvasContext.Call("bufferData", canvasContext.Get("ELEMENT_ARRAY_BUFFER"), jsElem, canvasContext.Get("STATIC_DRAW"))

	glRef.EnableVertexAttribArray(0)
	glRef.VertexAttribPointer(0, 2, webgl.FLOAT, false, 12, 0)
	glRef.EnableVertexAttribArray(1)
	glRef.VertexAttribPointer(1, 4, webgl.UNSIGNED_BYTE, true, 12, 8)

	glRef.BindVertexArray(&webgl2.VertexArrayObject{})
	glRef.BindBuffer(webgl2.ARRAY_BUFFER, &webgl.Buffer{})
	glRef.BindBuffer(webgl2.ELEMENT_ARRAY_BUFFER, &webgl.Buffer{})
	glRef.BindVertexArray(_sp.VAO)

	_sp.NumberOfElements = len(_sp.Indices)

	_sp.Vertices = _sp.Vertices[:0]
	_sp.Indices = _sp.Indices[:0]
}

func (_sp *ShapeBatch) Render(cam *Camera2D) {
	_sp.finalize()

	//glRef.DrawElements(webgl2.LINES, _sp.NumberOfElements, webgl2.UNSIGNED_INT, 0)

	UseShader(&_sp.Shader)

	matrixArray := cam.projectMatrix.Data()

	buffer := new(bytes.Buffer)
	{
		err := binary.Write(buffer, binary.LittleEndian, matrixArray[:])
		Assert(err == nil, "Error writing buffer b1")
	}
	byte_slice1 := buffer.Bytes()

	b1 := js.Global().Get("Uint8Array").New(len(byte_slice1))
	js.CopyBytesToJS(b1, byte_slice1)

	martixJS := js.Global().Get("Float32Array").New(b1.Get("buffer"), b1.Get("byteOffset"), b1.Get("byteLength").Int()/4)
	for i := 0; i < len(matrixArray); i++ {
		martixJS.SetIndex(i, js.ValueOf(matrixArray[i]))
	}

	projmatrix_loc := canvasContext.Call("getUniformLocation", _sp.Shader.ShaderProgramID.Value_JS, "projection_matrix")
	canvasContext.Call("uniformMatrix4fv", projmatrix_loc, false, martixJS)

	viewMatrix := cam.viewMatrix.Data()

	viewBuffer := new(bytes.Buffer)
	{
		err := binary.Write(viewBuffer, binary.LittleEndian, viewMatrix[:])
		Assert(err == nil, "Error writing viewBuffer")
	}

	byte_slice2 := buffer.Bytes()

	b2 := js.Global().Get("Uint8Array").New(len(viewMatrix) * 4)
	js.CopyBytesToJS(b2, byte_slice2)

	viewMatrixJS := js.Global().Get("Float32Array").New(b2.Get("buffer"), b2.Get("byteOffset"), b2.Get("byteLength").Int()/4)
	for i := 0; i < len(viewMatrix); i++ {
		viewMatrixJS.SetIndex(i, js.ValueOf(viewMatrix[i]))
	}

	viewmatrix_loc := canvasContext.Call("getUniformLocation", _sp.Shader.ShaderProgramID.Value_JS, "view_matrix")
	canvasContext.Call("uniformMatrix4fv", viewmatrix_loc, false, viewMatrixJS)
	//setupShaders(&glRef)

	canvasContext.Call("drawElements", canvasContext.Get("TRIANGLES"), _sp.NumberOfElements, canvasContext.Get("UNSIGNED_INT"), 0)
	UnuseShader()

	//glRef.BindVertexArray(nil)
}
