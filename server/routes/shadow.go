package routes

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	png "image/png"
	"io"
	"math"
	"net/http"

	"github.com/crazy3lf/colorconv"
)

func ShadowRoute(writer http.ResponseWriter, request *http.Request) {
	res, _ := http.Get("https://static.wikia.nocookie.net/feheroes_gamepedia_en/images/e/e9/Alfonse_Heir_to_Openness_Face.webp/revision/latest/scale-to-width-down/500?cb=20240816044233")
	data, _ := io.ReadAll(res.Body)
	reader := bytes.NewReader(data)
	baseImage, e := png.Decode(reader)
	if e != nil {
		fmt.Println(e)
	}
	var bounds = baseImage.Bounds()
	var width = bounds.Max.X
	var height = bounds.Max.Y
	blackenedImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgbColor := baseImage.At(x, y)
			_, _, _, initialAlpha := rgbColor.RGBA()
			var h, s, l = colorconv.ColorToHSL(rgbColor)
			l = math.Min(math.Max(l*0.1, 0), 1)
			var backToRgb, _ = colorconv.HSLToColor(h, s, l)
			var newR, newG, newB, _ = backToRgb.RGBA()
			var colorWithAlpha = color.RGBA{
				R: uint8(newR),
				G: uint8(newG),
				B: uint8(newB),
				A: uint8(initialAlpha),
			}
			blackenedImage.Set(x, y, colorWithAlpha)
		}
	}

	png.Encode(writer, blackenedImage)
}
