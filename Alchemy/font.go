package alchemy

import (
	"image"
	"image/color"
	"image/draw"
	"io"
	"net/http"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"golang.org/x/image/math/fixed"
)

type FontBatch struct {
	charSet      map[rune]CharGlyph
	sPatch       SpriteBatch
	fontSettings FontBatchSettings
}

type FontBatchAtlas struct {
	charAtlasSet map[rune]CharAtlasGlyph
	textureAtlas Texture2D
	sPatch       SpriteBatch
	fontSettings FontBatchSettings
}

type CharGlyph struct {
	textureId     Texture2D
	size, bearing Vector2f
	advance       float32
}

type CharAtlasGlyph struct {
	uv1           Vector2f
	uv2           Vector2f
	size, bearing Vector2f
	advance       float32
}

type FontBatchSettings struct {
	FontSize, DPI, CharDistance, LineHeight float32
}

func (self *FontBatch) Init() {
	self.charSet = make(map[rune]CharGlyph)
	self.sPatch.Init("")
}

func (self *FontBatchAtlas) Init() {
	self.charAtlasSet = make(map[rune]CharAtlasGlyph)
	self.sPatch.Init("")

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

const GLYPH_ATLAS_GAP = 5

func LoadFontToAtlas(_fontPath string, _fontSettings *FontBatchSettings) FontBatchAtlas {
	var tempFont FontBatchAtlas
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

	max_width, max_height := int(0), int(0)

	dot := fixed.Point26_6{fixed.Int26_6(face.Metrics().Ascent), 0}
	for i := 33; i < 127; i++ {
		char := rune(i)

		_, img, _, _, ok := face.Glyph(dot, char)

		if !ok {
			LogF("Failed to load rune: %v", char)
		}

		max_width += img.Bounds().Dx() + GLYPH_ATLAS_GAP
		if max_height < img.Bounds().Dy() {
			max_height = img.Bounds().Dy()
		}
	}

	atlas_img := image.NewRGBA(image.Rect(0, 0, max_width, max_height))
	draw.Draw(atlas_img, atlas_img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.ZP, draw.Src)

	x_offset := 0
	y_offset := 0

	for i := 32; i < 127; i++ {
		char := rune(i)
		_, img, _, ad, ok := face.Glyph(dot, char)
		if !ok {
			LogF("Failed To Draw Letter: %v", char)
		}
		y_offset = img.Bounds().Dy()

		if char == ' ' {
			tempFont.charAtlasSet[char] = CharAtlasGlyph{uv1: Vector2fZero, uv2: Vector2fZero, size: Vector2fZero, bearing: Vector2fOne, advance: float32(ad) / float32(1<<6)}
			continue
		}

		bounds, _, ok := face.GlyphBounds(char)
		if !ok {
			LogF("Failed to load bounds of rune: %v", char)
		}

		tempFont.charAtlasSet[char] = CharAtlasGlyph{
			NewVector2f(float32(x_offset)/float32(max_width), 0),
			NewVector2f(float32(x_offset+img.Bounds().Dx())/float32(max_width), float32(y_offset)/float32(max_height)),
			NewVector2f(float32(img.Bounds().Dx()), float32(img.Bounds().Dy())),
			NewVector2f(float32(bounds.Max.X)/64.0, float32(-bounds.Max.Y)/64.0),
			float32(ad) / 64.0,
		}
		draw.Draw(atlas_img, image.Rect(x_offset, 0, x_offset+img.Bounds().Dx(), y_offset), img, image.ZP, draw.Src)
		x_offset += img.Bounds().Dx() + GLYPH_ATLAS_GAP
	}

	tempFont.textureAtlas = LoadTextureFromImg(atlas_img)
	return tempFont
}

func (self *FontBatchAtlas) DrawString(_text string, _position Vector2f, _scale float32, _tint RGBA8) {
	originalPos := _position
	cameraRelatedScale := _scale

	for i := 0; i < len(_text); i++ {
		charglyph := self.charAtlasSet[rune(_text[i])]

		if rune(_text[i]) == ' ' {
			originalPos.X += charglyph.advance * cameraRelatedScale
			continue
		}
		loc_pos := originalPos
		//loc_pos.X += charglyph.bearing.X * cameraRelatedScale
		loc_pos.Y += (charglyph.bearing.Y) * cameraRelatedScale
		//Shapes.DrawRect(loc_pos, charglyph.size.Scale(cameraRelatedScale), NewRGBA8(255, 255, 255, 255))
		self.sPatch.DrawSpriteBottomLeft(loc_pos, charglyph.size.Scale(cameraRelatedScale), charglyph.uv1, charglyph.uv2, &self.textureAtlas, _tint)
		if i < len(_text)-1 {
			originalPos.X += charglyph.advance * cameraRelatedScale
		}
	}
}

func (self *FontBatchAtlas) Render() {
	self.sPatch.Render(&Cam)
}
