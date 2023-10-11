package alchemy

import (
	"fmt"
	"image/png"
	"net/http"
	"reflect"
	"syscall/js"
	"unsafe"

	"github.com/gowebapi/webapi/graphics/webgl"
)

type Texture2D struct {
	Width, Height, bpp int
	textureId          *webgl.Texture
}

type Pixel struct {
	RGBA RGBA8
}

func New(r, g, b, a uint8) Pixel {
	var pixel Pixel
	pixel.RGBA = RGBA8{r, g, b, a}
	return pixel
}

func LoadPng(_filePath string) {

	var tempTexture Texture2D

	resp, err := http.Get(app_url + "/" + _filePath)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	img, err := png.Decode(resp.Body)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	fmt.Printf("%v\n", img.Bounds().Dx())

	resp.Body.Close()

	tempTexture.Width = img.Bounds().Dx()
	LogF("%v", img.Bounds().Dx())
	tempTexture.Height = img.Bounds().Dy()
	pixels := make([]Pixel, tempTexture.Height*tempTexture.Width)

	for y := 0; y < tempTexture.Height; y++ {
		for x := 0; x < tempTexture.Width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[y*tempTexture.Width+x] = New(uint8(r), uint8(g), uint8(b), uint8(a))
		}
	}

	tempTexture.textureId = glRef.CreateTexture()
	glRef.ActiveTexture(webgl.TEXTURE0)
	glRef.BindTexture(webgl.TEXTURE_2D, tempTexture.textureId)

	glRef.TexParameteri(webgl.TEXTURE_2D, webgl.TEXTURE_MIN_FILTER, int(webgl.NEAREST))
	glRef.TexParameteri(webgl.TEXTURE_2D, webgl.TEXTURE_MAG_FILTER, int(webgl.NEAREST))
	glRef.TexParameteri(webgl.TEXTURE_2D, webgl.TEXTURE_WRAP_S, int(webgl.CLAMP_TO_EDGE))
	glRef.TexParameteri(webgl.TEXTURE_2D, webgl.TEXTURE_WRAP_T, int(webgl.CLAMP_TO_EDGE))

	jsPixels := js.Global().Get("Uint8Array").New(len(pixels) * 4)
	var pixelsBytes []byte
	headerPixels := (*reflect.SliceHeader)(unsafe.Pointer(&pixelsBytes))
	headerPixels.Cap = cap(pixels) * 4
	headerPixels.Len = len(pixels) * 4
	headerPixels.Data = uintptr(unsafe.Pointer(&pixels[0]))

	js.CopyBytesToJS(jsPixels, pixelsBytes)

	canvasContext.Call("texImage2D", canvasContext.Get("TEXTURE_2D"), 0, canvasContext.Get("RGBA8"), tempTexture.Width, tempTexture.Height, 0, canvasContext.Get("RGBA"), canvasContext.Get("UNSIGNED_BYTE"), jsPixels)
	//glRef.TexImage2D(webgl.TEXTURE_2D, 0, int(webgl2.RGBA8), tempTexture.Width, tempTexture.Height, 0, webgl2.RGBA, webgl2.UNSIGNED_BYTE, pixels)
}