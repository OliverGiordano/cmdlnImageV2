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
type Ccolor struct {
    cVal termColors.RGBColor
    aVal float64
}
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
    chunkSizeX := int(imgSize.X / 50)
    chunkSizeY := int(imgSize.Y / 50) //I want the image in the terminal to be 60x60
    colorList := []color.Color{}
    rowColorList := []Ccolor{}
    lowerRowColorList := []Ccolor{}
    //--get the background color
    bgColor := termenv.BackgroundColor().Sequence(true)
    colorBg := strings.Split(string(bgColor), ";")
    colorBg = colorBg[2:]
    bgR, _ := strconv.Atoi(colorBg[0])
    bgG, _ := strconv.Atoi(colorBg[1])
    bgB, _ := strconv.Atoi(colorBg[2])
    //to look at the image I scan side-length/60 by side-length/60 chunks, at a time and
    //average the colors in that chunk, that way the image will be scaled down to 60x60
    for y := 0; y*chunkSizeY < imgSize.Y; y += 1 {
        for x := 0; x*chunkSizeX < imgSize.X; x += 1 {
    	colorList = []color.Color{} //empty the color list so it doesnt grow each chunk we look at
	    for yC := y*chunkSizeY - chunkSizeY; yC < y*chunkSizeY; yC++ {
		for xC := x*chunkSizeX - chunkSizeX; xC < x*chunkSizeX; xC++ {
		    colorList = append(colorList, image.At(xC, yC))
		}
	    }
	    c, a := evalColor(colorList)
		// I use the lower half block ascii character and change the background color to
		// double the resolution, because of this I need to write these rows all at once, so I save the top half
		// and write it when I write the bottem half
		//}
	    if y%2 != 0 {
		lowerRowColorList = append(lowerRowColorList, Ccolor{c, a})
	    } else {
	        rowColorList = append(rowColorList, Ccolor{c, a})
	    }
	}
	if y%2 != 0 {
	    rowBlank := true
	    for i, _ := range(rowColorList){
		/*if lowerRowColorList[i].Basic().RGB()[0]+lowerRowColorList[i].Basic().RGB()[1]+lowerRowColorList[i].Basic().RGB()[2] == 0 {//if color is to dark set it to background color
		    lowerRowColorList[i] = termColors.RGB(uint8(bgR), uint8(bgG), uint8(bgB))
		} else {
		    rowBlank = false
		}
		if rowColorList[i].Basic().RGB()[0]+rowColorList[i].Basic().RGB()[1]+rowColorList[i].Basic().RGB()[2] == 0 {
		    rowColorList[i] = termColors.RGB(uint8(bgR), uint8(bgG), uint8(bgB))
		} else {
		    rowBlank = false
		}*/
		if lowerRowColorList[i].aVal == 0 {
		    lowerRowColorList[i] = Ccolor{termColors.RGB(uint8(bgR), uint8(bgG), uint8(bgB)), 1.0}
		} else {
		    rowBlank = false
		}
		if rowColorList[i].aVal == 0 {
		    rowColorList[i] = Ccolor{termColors.RGB(uint8(bgR), uint8(bgG), uint8(bgB)), 1.0}
		} else {
		    rowBlank = false
		}
	    }
	    if (rowBlank == true){
		lowerRowColorList = nil
		rowColorList = nil
		continue
	    }
	    for _, _ = range(rowColorList){
		fullStyle := termColors.NewRGBStyle(lowerRowColorList[0].cVal, rowColorList[0].cVal)
		fullStyle.Printf("\u2584") // this is the ascii character for the lower half block
		lowerRowColorList = lowerRowColorList[1:]
		rowColorList = rowColorList[1:]
	    }
	    fmt.Println()
        }
    }
}

func printAscii(image image.Image) {
	return // must implement
}

func evalColor(colorList []color.Color) (termColors.RGBColor, float64) {
	var tR uint32 = 0
	var tG uint32 = 0
	var tB uint32 = 0
	var tA float64 = 0.0
	// I average colors using the method outlined in this article & video, it is interesting watch&read
	//if you get a chance https://sighack.com/post/averaging-rgb-colors-the-right-way
	for _, color := range colorList {
	    r, g, b, a := color.RGBA()
	    r = r / 256
	    g = g / 256
	    b = b / 256
	    tR += (r * r)
	    tG += (g * g)
	    tB += (b * b)
	    tA += float64(a/256)
	}
	tR = uint32(math.Sqrt(float64(tR) / float64(len(colorList))))
	tG = uint32(math.Sqrt(float64(tG) / float64(len(colorList))))
	tB = uint32(math.Sqrt(float64(tB) / float64(len(colorList))))
	c := termColors.RGB(uint8(tR), uint8(tG), uint8(tB))
	return c, tA
}
