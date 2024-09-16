package shadows

import (
	"bytes"
	"fehdle/routes/common"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
	"regexp"
	"strconv"

	"github.com/crazy3lf/colorconv"
	cron "github.com/robfig/cron/v3"
)

var cachedImage image.Image
var shadowUnitId string = ""

func Route(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	var brightnessModifier = 0
	if len(request.Form["b"]) > 0 {
		parsedMod, _ := strconv.Atoi(request.Form["b"][0])
		brightnessModifier = parsedMod
	}

	if cachedImage == nil {
		var foundHero, _ = common.FindHero(shadowUnitId)
		data := common.GetImage(foundHero.Title.WikiName, "portrait")
		reader := bytes.NewReader(data)
		var storedImage, e = png.Decode(reader)
		cachedImage = storedImage
		if e != nil {
			fmt.Println(e)
		}
	}

	var bounds = cachedImage.Bounds()
	var width = bounds.Max.X
	var height = bounds.Max.Y
	blackenedImage := image.NewRGBA(image.Rect(0, 0, width, height))
	var appliedBrightness = 0 + (0.01 * float64(brightnessModifier))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgbColor := cachedImage.At(x, y)
			_, _, _, initialAlpha := rgbColor.RGBA()
			if initialAlpha == 0 {
				continue
			}
			var h, s, l = colorconv.ColorToHSL(rgbColor)
			l = math.Min(l*appliedBrightness, 1)
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

func UpdateGoroutine() {
	updateCron := cron.New()
	updateCron.AddFunc("* * * * *", func() {
		var unit = common.UpdateMainUnit().CargoQuery[0].Title
		shadowUnitId = unit.IntID
		cachedImage = nil
	})

	updateCron.Start()

	select {}
}

func GuessShadow(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	if len(request.Form["intId"]) == 0 {
		writer.Write([]byte("Invalid payload format: expected a valid IntID, received nothing"))
		writer.WriteHeader(400)
	}

	byteIntId := []byte(request.Form["intId"][0])
	match, e := regexp.Match("[0-9]{1-4}", byteIntId)
	if !match || e != nil {
		writer.Write([]byte("Invalid payload format: expected a valid IntID, received " + string(byteIntId)))
		writer.WriteHeader(400)
	}

	var intIdString = string(byteIntId)
	_, heroFindError := common.FindHero(intIdString)
	if heroFindError != nil {
		writer.WriteHeader(404)
	} else {
		if shadowUnitId == intIdString {
			writer.Write([]byte("1"))
		} else {
			writer.Write([]byte("0"))
		}
	}
}
