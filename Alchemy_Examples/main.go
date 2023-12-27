package main

import (
	"math"

	alchemy "github.com/mhamedGd/Alchemy"
)

func main() {
	var rotation float32 = 90.0
	var midPoint alchemy.Vector2f = alchemy.Vector2fZero
	var midScreen alchemy.Vector2f

	var velocity alchemy.Vector2f = alchemy.Vector2fZero
	var direction alchemy.Vector2f
	speed := float32(0.3)

	var inputAxis alchemy.Vector2f = alchemy.Vector2fZero
	var dpad_modifier alchemy.Vector2f = alchemy.Vector2fZero

	var bgl_texture alchemy.Texture2D
	//var tile alchemy.Texture2D
	//var mapp alchemy.Texture2D
	var font alchemy.FontBatch
	var fontAtlas alchemy.FontBatchAtlas
	game := alchemy.App{
		Width:  1920,
		Height: 1080,
		Title:  "Test",
		OnStart: func() {
			//fmt.Printf("On Start\n")
			alchemy.LogF("STARTED\n")

			midScreen = alchemy.NewVector2f(400.0, 300.0)

			alchemy.BindInput("Up", alchemy.KEY_W)
			alchemy.BindInput("Down", alchemy.KEY_S)
			alchemy.BindInput("Right", alchemy.KEY_D)
			alchemy.BindInput("Left", alchemy.KEY_A)
			alchemy.BindInput("Zoom In", alchemy.KEY_E)
			alchemy.BindInput("Zoom Out", alchemy.KEY_Q)

			alchemy.Shapes.LineWidth = .5

			bgl_texture = alchemy.LoadPng("Assets/tile_0004.png")
			//tile = alchemy.LoadPng("Assets/Asset 2.png")
			//mapp = alchemy.LoadPng("Assets/Map.png")

			alchemy.MainButton_Pressed.AddListener(func(i ...int) {
				dpad_modifier.Y += 1.0
			})
			alchemy.MainButton_Released.AddListener(func(i ...int) {
				dpad_modifier.Y -= 1.0
			})

			alchemy.SideButton_Pressed.AddListener(func(i ...int) {
				dpad_modifier.Y -= 1.0
				alchemy.LogF(dpad_modifier.ToString())

			})

			alchemy.SideButton_Released.AddListener(func(i ...int) {
				dpad_modifier.Y += 1.0
				alchemy.LogF(dpad_modifier.ToString())
			})

			alchemy.DPadLeft_Pressed.AddListener(func(i ...int) {
				dpad_modifier.X -= 1.0
			})
			alchemy.DPadLeft_Released.AddListener(func(i ...int) {
				dpad_modifier.X += 1.0
			})

			alchemy.DPadRight_Pressed.AddListener(func(i ...int) {
				dpad_modifier.X += 1.0
			})
			alchemy.DPadRight_Released.AddListener(func(i ...int) {
				dpad_modifier.X -= 1.0
			})

			font_settings := alchemy.FontBatchSettings{
				FontSize: 48, DPI: 124, CharDistance: 4, LineHeight: 16,
			}

			font = alchemy.LoadFont("Assets/m5x7.ttf", &font_settings)
			fontAtlas = alchemy.LoadFontToAtlas("Assets/m5x7.ttf", &font_settings)
			alchemy.ScaleView(4)
		},
		OnUpdate: func(dt float64) {

			zoomAxis := 500.0 * float32(dt) * (alchemy.GetActionStrength("Zoom In") - alchemy.GetActionStrength("Zoom Out"))
			alchemy.IncreaseScaleU(zoomAxis)
			inputAxis.Y = alchemy.GetActionStrength("Up") - (alchemy.GetActionStrength("Down"))
			inputAxis.X = alchemy.GetActionStrength("Right") - (alchemy.GetActionStrength("Left"))

			velocity.X = alchemy.LerpFloat32(velocity.X, (inputAxis.X+dpad_modifier.X)*speed, float32(dt)*2.5)
			velocity.Y = alchemy.LerpFloat32(velocity.Y, (inputAxis.Y+dpad_modifier.Y)*speed, float32(dt)*2.5)
			//alchemy.ScrollView(alchemy.Vector2fRight.Scale(velocity.X * 3.0))

			rotation -= float32(dt*600.0) * velocity.X
			rotation = float32(math.Mod(float64(rotation), 360))
			direction = alchemy.Vector2fRight.Rotate(rotation, alchemy.Vector2fZero)

			midPoint.Y += velocity.Y * direction.Y
			midPoint.X += velocity.Y * direction.X
			alchemy.ScrollTo(midPoint)

			alchemy.LogF("FPS: %v", 1.0/float32(dt))
		},
		OnDraw: func() {
			alchemy.Shapes.DrawFillRectRotated(midScreen, alchemy.Vector2fOne.Scale(50.0), alchemy.NewRGBA8(255, 100, 230, 255), rotation)
			//alchemy.Shapes.DrawRect(midPoint, alchemy.NewVector2f(10, 10), alchemy.NewRGBA8(255, 255, 0, 255))
			alchemy.Shapes.DrawTriangleRotated(midPoint, alchemy.NewVector2f(2.0, 4.0), alchemy.NewRGBA8(255, 0, 0, 255), rotation)
			alchemy.Sprites.DrawSpriteOrigin(alchemy.NewVector2f(2, 0.0), alchemy.Vector2fZero, alchemy.Vector2fOne, &bgl_texture, alchemy.NewRGBA8(255, 255, 255, 255))

			//alchemy.Sprites.DrawSpriteOrigin(alchemy.NewVector2f(2, -50.0), alchemy.Vector2fZero, alchemy.Vector2fOne, &mapp, alchemy.NewRGBA8(255, 255, 255, 255))
			//alchemy.Sprites.DrawSpriteOrigin(alchemy.NewVector2f(2, 50.0), alchemy.Vector2fZero, alchemy.Vector2fOne, &tile, alchemy.NewRGBA8(255, 255, 255, 255))
			for i := 0; i < 1; i++ {
				//font.DrawString("Baghdad Game Lab -", alchemy.Vector2fOne.Scale(float32(i)), 0.5, alchemy.NewRGBA8(255, 255, 255, 255))
				fontAtlas.DrawString("Baghdad Game Lab -", alchemy.Vector2fOne.Scale(float32(i)), 0.5, alchemy.NewRGBA8(255, 255, 255, 255))
			}
			font.Render()
			fontAtlas.Render()
		},
		OnEvent: func(ae *alchemy.AppEvent) {

		},
	}

	alchemy.Run(&game)

}
