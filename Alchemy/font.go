package alchemy

import (
	"io"
	"net/http"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func LoadFont(_fontPath string) Texture2D {
	resp, err := http.Get(app_url + "/" + _fontPath)
	if err != nil {
		LogF(err.Error())
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		LogF(err.Error())
	}

	f, err := opentype.Parse(data)

	// glyphIndex, err := f.GlyphIndex(buff, 'C')
	// if err != nil {
	// 	LogF(err.Error())
	// }
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    32,
		DPI:     512,
		Hinting: font.HintingFull,
	})

	if err != nil {
		LogF(err.Error())
	}

	dot := fixed.Point26_6{fixed.Int26_6(face.Metrics().Ascent), 0}
	_, ima, _, _, ok := face.Glyph(dot, 'A')
	if !ok {
		LogF("Failed to load Glyph")
	}

	//img := image.NewAlpha(dr)

	tempT := LoadImageFromImg(ima)
	return tempT
}
