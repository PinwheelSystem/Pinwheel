package main

import (
	"math/rand"
	"fmt"
	"time"
	"os"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/yuin/gopher-lua"
	"github.com/PinwheelSystem/bitmap"
	"github.com/PinwheelSystem/PaletteNom"
	"github.com/akamensky/argparse"
)

const version = "v0.1.0"
const res int = 128 // Resolution of the *screen* ("internal") . Might change later in development. (res x res, 128 x 128)
const scale = 4 // Resolution scale (contributes to the size of the *window*)
var palette [][]uint8 = make([][]uint8, 64) // Array of array of RGB values ([[R, G, B], [R, G, B], ...])
var pixelbuf []byte = make([]byte, res * res * 4) // Pixel backbuffer (basically our VRAM)
var start time.Time
var font map[string]bitmap.Glyph
var renderer *sdl.Renderer
var pinwheel Pinwheel

func main() {
	parser := argparse.NewParser("pinwheel", "An awesome little fantasy computer designed to be simple.")

	palettedir := parser.String("P", "palettedir", &argparse.Options{
	 	Required: false,
		Help: "The palette folder path",
		Default: "Palettes/",
	})

	palettefile := parser.String("m", "palette", &argparse.Options{
	 	Required: false,
		Help: "Name of the palette image file",
		Default: "AAP-64.png",
	})

	fontfile := parser.String("f", "font", &argparse.Options{
	 	Required: false,
		Help: "Name of the bitmap font file",
		Default: "m5x7.png",
	})

	program := parser.String("p", "program", &argparse.Options{
	 	Required: false,
		Help: "Path to a Pinwheel program",
		Default: "program.lua",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	data := map[string]interface{}{
		"palettedir": *palettedir,
		"palettefile": *palettefile,
		"fontfile": *fontfile,
		"program": *program,
	}

	palettelib := palettenom.New()
	colors, err := palettelib.Load(*palettedir + "AAP-64.png")
	bm := bitmap.New()
	font = bm.Load("m5x7.png")

	for i := uint8(0); i < 64; i++ {
		r, g, b, _ := colors[i].RGBA()		
		palette[i] = []uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
	}

	pinwheel = Init(data)
	pinwheel.Start()	
}

func randf(r int) int {
	return rand.Intn(r)
}

// vpoke(addr, val)
func PWvPoke(L *lua.LState) int {
	addr := L.ToInt(1)
	val := L.ToInt(2)
	pixelbuf := *pinwheel.vram

	pixelbuf[addr] = byte(val)

	return 1
}

// plot(x, y, color)
func PWplot(L *lua.LState) int {
	x := L.ToInt(1)
	y := L.ToInt(2)
	color := L.ToInt(3)

	c := palette[color]
	setpixel(x, y, int(c[0]), int(c[1]), int(c[2]))

	return 1
}

// termprint(text)
func PWtermPrint(L *lua.LState) int {
	text := L.ToString(1)

	fmt.Println(text)

	return 1
}

// time() -> float64
func PWtime(L *lua.LState) int {
	current := time.Now()
	duration := float64(current.Sub(start)) / 1000000 / 1000

	L.Push(lua.LNumber(duration))

	return 1
}

// pchar(single_char, x, y, color)
func PWpchar(L *lua.LState) int {
	char := L.ToString(1)
	x := L.ToInt(2)
	y := L.ToInt(3)
	color := L.ToInt(4)

	c := palette[color]
	xx := x
	yy := y

	for i := 0; i < 8; i++ {
	 	bin := font[char].Data[i]
		binarr := strings.Split(bin, "")

		for _, pix := range binarr {
			if pix == "1" { setpixel(xx, yy, int(c[0]), int(c[1]), int(c[2])) }
		 	xx += 1
		}
		yy += 1
		xx = x
	}

	L.Push(lua.LNumber(font[char].Width))

	return 1
}

// vertline(x, color)
func PWvertline(L *lua.LState) int {
	x := L.ToInt(1)
	color := L.ToInt(2)
	
	for i := 0; i < res; i++ {
		c := palette[color]
		setpixel(x, i, int(c[0]), int(c[1]), int(c[2]))
	}

	return 1
}

// horizline(y, color)
func PWhorizline(L *lua.LState) int {
	y := L.ToInt(1)
	color := L.ToInt(2)
	
	for i := 0; i < res; i++ {
		c := palette[color]
		setpixel(i, y, int(c[0]), int(c[1]), int(c[2]))
	}

	return 1
}

func PWclear(L *lua.LState) int {
	c := palette[0]

	for y := 0; y < res; y++ {
		for x := 0; x < res; x++ {
			setpixel(x, y, int(c[0]), int(c[1]), int(c[2]))
		}
	}

	return 1
}

// print(text, x, y, color)
func PWprint(L *lua.LState) int {
	text := L.ToString(1)
	x := L.ToInt(2)
	y := L.ToInt(3)
	color := L.ToInt(4)

	c := palette[color]
	xx := x
	sx := x
	yy := y

	for _, ch := range text {
		char := font[string(ch)]
		for i := 0; i < char.Height; i++ {
		 	bin := char.Data[i]

			for _, pix := range bin {
				if string(pix) == "1" { setpixel(xx + sx, char.Y + yy, int(c[0]), int(c[1]), int(c[2])) }
			 	xx += 1
			}
			yy += 1
			xx = x
		}
		sx += char.Width
		yy = y
	}

	return 1
}

func setpixel(x, y, r, g, b int) {
	offset := ( pinwheel.screen.Width * 4 * y ) + x * 4;
	pixelbuf := *pinwheel.vram
	pixelbuf[offset + 0] = byte(b)
	pixelbuf[offset + 1] = byte(g)
	pixelbuf[offset + 2] = byte(r)
	pixelbuf[offset + 3] = sdl.ALPHA_OPAQUE;
}