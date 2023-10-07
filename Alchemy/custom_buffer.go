package alchemy

import (
	"reflect"
	"unsafe"
)

func VertexSliceAsBytes(_verts []Vertex) []byte {
	n := int(vertexByteSize) * len(_verts)

	up := unsafe.Pointer(&(_verts[0]))
	pi := (*[1]byte)(up)
	buf := (*pi)[:]
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	sh.Len = n
	sh.Cap = n

	return buf
}
