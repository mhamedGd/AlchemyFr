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
	speed := float32(3)

	var inputAxis alchemy.Vector2f = alchemy.Vector2fZero
	var dpad_modifier float32 = 0.0

	game := alchemy.App{
		Width:  800,
		Height: 600,
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

			alchemy.Shapes.LineWidth = 1.

			alchemy.DPadUp_Pressed.AddListener(func(i ...int) {
				dpad_modifier += 1.0
			})
			alchemy.DPadUp_Released.AddListener(func(i ...int) {
				dpad_modifier -= 1.0
			})

		},
		OnUpdate: func(dt float64) {
			zoomAxis := 500.0 * float32(dt) * (alchemy.GetActionStrength("Zoom In") - alchemy.GetActionStrength("Zoom Out"))
			alchemy.IncreaseScaleU(zoomAxis)

			inputAxis.Y = alchemy.GetActionStrength("Up") - (alchemy.GetActionStrength("Down"))

			alchemy.LogF("%v", dpad_modifier)
			velocity.X = alchemy.LerpFloat32(velocity.X, (alchemy.GetActionStrength("Right")-alchemy.GetActionStrength("Left"))*speed, float32(dt)*2.5)
			velocity.Y = alchemy.LerpFloat32(velocity.Y, (inputAxis.Y+dpad_modifier)*speed, float32(dt)*2.5)
			//alchemy.ScrollView(alchemy.Vector2fRight.Scale(velocity.X * 3.0))

			rotation -= float32(dt*100.0) * velocity.X
			rotation = float32(math.Mod(float64(rotation), 360))
			direction = alchemy.Vector2fRight.Rotate(rotation, alchemy.Vector2fZero)

			midPoint.Y += velocity.Y * direction.Y
			midPoint.X += velocity.Y * direction.X
		},
		OnDraw: func() {
			alchemy.Shapes.DrawFillRectRotated(midScreen, alchemy.Vector2fOne.Scale(50.0), alchemy.NewRGBA8(255, 100, 230, 255), rotation)
			//alchemy.Shapes.DrawRect(midPoint, alchemy.NewVector2f(10, 10), alchemy.NewRGBA8(255, 255, 0, 255))
			alchemy.Shapes.DrawTriangleRotated(midPoint, alchemy.NewVector2f(10.0, 20.0), alchemy.NewRGBA8(255, 0, 0, 255), rotation)
		},
		OnEvent: func(ae *alchemy.AppEvent) {

		},
	}

	alchemy.Run(&game)
}
