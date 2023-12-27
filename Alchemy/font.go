package alchemy

import (
	"image"
	"image/color"
	"image/draw"
	"io"
	"net/http"

	"github.com/gowebapi/webapi/graphics/webgl"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"golang.org/x/image/math/fixed"
)

type FontBatch struct {
	charSet      map[rune]CharGlyph
	sPatch       SpriteBatch
	fontSettings FontBatchSettings
}

type CharGlyph struct {
	textureId     Texture2D
	size, bearing Vector2f
	advance       float32
}

type FontBatchSettings struct {
	FontSize, DPI, CharDistance, LineHeight float32
}

func (self *FontBatch) Init() {
	self.charSet = make(map[rune]CharGlyph)
	self.sPatch.Init("")
	glRef.PixelStorei(webgl.UNPACK_ALIGNMENT, 1)
}

func LoadFont(_fontPath string, _fontSettings *FontBatchSettings) FontBatch {
	var tempFont FontBatch
	tempFont.Init()
	tempFont.fontSettings = *_fontSettings

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

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(_fontSettings.FontSize),
		DPI:     float64(_fontSettings.DPI),
		Hinting: font.HintingFull,
	})
	if err != nil {
		LogF(err.Error())
	}

	dot := fixed.Point26_6{fixed.Int26_6(face.Metrics().Ascent), 0}
	for i := 32; i < 127; i++ {
		char := rune(i)

		_, img, _, ad, ok := face.Glyph(dot, char)

		if !ok {
			LogF("Failed to load rune: %v", char)
		}

		if char == ' ' {
			tempFont.charSet[char] = CharGlyph{textureId: Texture2D{}, size: Vector2fZero, bearing: Vector2fOne, advance: float32(ad) / float32(1<<6)}
			continue
		}

		bounds, _, ok := face.GlyphBounds(char)
		if !ok {
			LogF("Failed to load bounds of rune: %v", char)
		}

		texture := LoadTextureFromImg(img)

		charGlyph := CharGlyph{
			texture,
			NewVector2f(float32(texture.Width), float32(texture.Height)),
			NewVector2f(float32(bounds.Max.X)/64.0, float32(-bounds.Max.Y)/64.0),
			float32(ad) / 64.0,
		}

		tempFont.charSet[char] = charGlyph
	}

	return tempFont
}

func (self *FontBatch) DrawString(_text string, _position Vector2f, _scale float32, _tint RGBA8) {
	originalPos := _position
	cameraRelatedScale := _scale

	for i := 0; i < len(_text); i++ {
		charglyph := self.charSet[rune(_text[i])]

		if rune(_text[i]) == ' ' {
			originalPos.X += charglyph.advance * cameraRelatedScale
			continue
		}
		loc_pos := originalPos
		//loc_pos.X += charglyph.bearing.X * cameraRelatedScale
		loc_pos.Y += (charglyph.bearing.Y) * cameraRelatedScale
		//Shapes.DrawRect(loc_pos, charglyph.size.Scale(cameraRelatedScale), NewRGBA8(255, 255, 255, 255))
		self.sPatch.DrawSpriteBottomLeft(loc_pos, charglyph.size.Scale(cameraRelatedScale), Vector2fZero, Vector2fOne, &charglyph.textureId, _tint)
		if i < len(_text)-1 {
			originalPos.X += charglyph.advance * cameraRelatedScale
		}
	}
}

func (self *FontBatch) Render() {
	self.sPatch.Render(&Cam)
}

func CreateFontAtlas(_fontPath string, _fontSettings *FontBatchSettings) Texture2D {
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

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(_fontSettings.FontSize),
		DPI:     float64(_fontSettings.DPI),
		Hinting: font.HintingFull,
	})
	if err != nil {
		LogF(err.Error())
	}

	max_width, max_height := int(0), int(0)

	dot := fixed.Point26_6{fixed.Int26_6(face.Metrics().Ascent), 0}
	for i := 33; i < 127; i++ {
		char := rune(i)

		_, img, _, _, ok := face.Glyph(dot, char)

		if !ok {
			LogF("Failed to load rune: %v", char)
		}

		max_width += img.Bounds().Dx()
		if max_height < img.Bounds().Dy() {
			max_height = img.Bounds().Dy()
		}
	}

	atlas_img := image.NewRGBA(image.Rect(0, 0, max_width, max_height))
	draw.Draw(atlas_img, atlas_img.Bounds(), &image.Uniform{color.RGBA{200, 100, 0, 255}}, image.ZP, draw.Src)

	x_offset := 0
	y_offset := 0
	for i := 33; i < 127; i++ {
		char := rune(i)
		_, img, _, _, ok := face.Glyph(dot, char)
		y_offset = img.Bounds().Dy()
		if !ok {
			LogF("Failed To Draw Letter: %v", char)
		}
		draw.Draw(atlas_img, image.Rect(x_offset, 0, x_offset+img.Bounds().Dx(), y_offset), img, image.ZP, draw.Src)
		x_offset += img.Bounds().Dx()

	}

	var tTexture = LoadTextureFromImg(atlas_img)
	return tTexture
}
