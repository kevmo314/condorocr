package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

//go:embed simsun.ttf
var simsunttf []byte

func isWhite(c color.Color) bool {
	r, g, b, a := c.RGBA()
	return r >= 0xdddd && g >= 0xdddd && b >= 0xdddd && a == 0xffff
}

func isPresent(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a != 0
}

func score(c *freetype.Context, s string, dx int, img image.Image) float64 {
	rgba := image.NewRGBA(img.Bounds())

	c.SetDst(rgba)

	if _, err := c.DrawString(s, freetype.Pt(38+24*dx, 75)); err != nil {
		panic(err)
	}

	// iterate over all pixels in rgba, count the number of pixels where rgba is white and img is white
	count := 0
	total := 0
	for x := 38 + 24*dx; x < 52+24*dx; x++ {
		for y := 0; y < 100; y++ {
			rgbaColor := rgba.At(x, y)
			imgColor := img.At(x, y)
			if isPresent(rgbaColor) && isWhite(imgColor) {
				count++
			}
			if isPresent(rgbaColor) {
				total++
			}
		}
	}
	return float64(count) / float64(total)
}

func main() {
	img, err := png.Decode(os.Stdin)
	if err != nil {
		panic(err)
	}

	// load
	f, err := freetype.ParseFont(simsunttf)
	if err != nil {
		panic(err)
	}

	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(f)
	ctx.SetFontSize(48)
	ctx.SetClip(img.Bounds())
	ctx.SetSrc(image.White)
	ctx.SetHinting(font.HintingFull)

	// search
	out := ""
	for c := 0; c < 14; c++ {
		if c == 2 || c == 4 {
			out += "/"
		}
		if c == 8 {
			out += " "
		}
		if c == 10 || c == 12 {
			out += ":"
		}
		best := ""
		bestScore := -1.0
		for i := 0; i < 10; i++ {
			score := score(ctx, fmt.Sprintf("%d", i), len(out), img)
			if score > bestScore {
				bestScore = score
				best = fmt.Sprintf("%d", i)
			}
		}
		out += best
	}

	fmt.Print(out)
}
