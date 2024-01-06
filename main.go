package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strconv"
	"strings"

	termColors "github.com/gookit/color"
	"github.com/muesli/termenv"
)

func main() {
	args := os.Args[1:]
	f, err := os.Open(args[0])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	image, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	if !(len(args) > 1) {
		printBlocks(image)
		return
	}
	if args[1] == "--ascii" {
		printAscii(image)
	} else if args[1] == "--blocks" {
		printBlocks(image)
	}
}

func printBlocks(image image.Image) {
	imgSize := image.Bounds().Size()
	chunkSizeX := int(imgSize.X / 60)
	chunkSizeY := int(imgSize.Y / 60)
	colorList := []color.Color{}
	rowColorList := []termColors.RGBColor{}
	bgColor := termenv.BackgroundColor().Sequence(true)
	colorBg := strings.Split(string(bgColor), ";")
	colorBg = colorBg[2:]
	bgR, _ := strconv.Atoi(colorBg[0])
	bgG, _ := strconv.Atoi(colorBg[1])
	bgB, _ := strconv.Atoi(colorBg[2])
	for y := 0; y*chunkSizeY < imgSize.Y; y += 1 {
		for x := 0; x*chunkSizeX < imgSize.X; x += 1 {
			colorList = []color.Color{}
			for yC := y*chunkSizeY - chunkSizeY; yC < y*chunkSizeY; yC++ {
				for xC := x*chunkSizeX - chunkSizeX; xC < x*chunkSizeX; xC++ {
					colorList = append(colorList, image.At(xC, yC))
				}
			}
			c, cb := evalColor(colorList)
			if y%2 != 0 {
				if c.Basic().RGB()[0]+c.Basic().RGB()[1]+c.Basic().RGB()[2] < 20 {
					c = termColors.RGB(uint8(bgR), uint8(bgG), uint8(bgB))
				}
				if rowColorList[0].Basic().RGB()[0]+rowColorList[0].Basic().RGB()[1]+rowColorList[0].Basic().RGB()[2] < 20 {
					rowColorList[0] = termColors.RGB(uint8(bgR), uint8(bgG), uint8(bgB))
				}
				fullStyle := termColors.NewRGBStyle(c, rowColorList[0])
				fullStyle.Printf("\u2584")
				rowColorList = rowColorList[1:]
			} else {
				rowColorList = append(rowColorList, cb)
			}

		}
		if y%2 != 0 {
			fmt.Println()
		}

	}
}

func printAscii(image image.Image) {
	return
}

func evalColor(colorList []color.Color) (termColors.RGBColor, termColors.RGBColor) {
	var tR uint32 = 0
	var tG uint32 = 0
	var tB uint32 = 0
	for _, color := range colorList {
		r, g, b, _ := color.RGBA()
		r = r / 256
		g = g / 256
		b = b / 256

		tR += (r * r)
		tG += (g * g)
		tB += (b * b)

	}
	tR = uint32(math.Sqrt(float64(tR) / float64(len(colorList))))
	tG = uint32(math.Sqrt(float64(tG) / float64(len(colorList))))
	tB = uint32(math.Sqrt(float64(tB) / float64(len(colorList))))
	//fmt.Println(tR, tG, tB)
	//fmt.Println(tR, tG, tB)
	c := termColors.RGB(uint8(tR), uint8(tG), uint8(tB))
	cb := termColors.RGB(uint8(tR), uint8(tG), uint8(tB))
	return c, cb
}
