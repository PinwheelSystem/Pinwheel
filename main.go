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
	"github.com/chuckpreslar/emission"
	"layeh.com/gopher-luar"
)

const version = "v0.1.0"
const res int = 128 // Resolution of the *screen* ("internal") . Might change later in development. (res x res, 128 x 128)
const scale = 4 // Resolution scale (contributes to the size of the *window*)
var palette [][]uint8 = make([][]uint8, 64) // Array of array of RGB values ([[R, G, B], [R, G, B], ...])
var pixelbuf []byte = make([]byte, res * res * 4) // Pixel backbuffer (basically our VRAM)
var start time.Time
var font map[string]bitmap.Glyph
var renderer *sdl.Renderer

func main() {
	palettelib := palettenom.New()
	colors/*, _*/ := palettelib.Load("aap-64.png")
	bm := bitmap.New()
	font = bm.Load("m5x7.png")
	emitter := emission.NewEmitter()
	events := emitter

	for i := uint8(0); i < 64; i++ {
		r, g, b, _ := colors[i].RGBA()
		palette[i] = []uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
	}

	var program string = "program.lua"
	if len(os.Args) != 1{
		program = os.Args[1] 
	}
	
	L := lua.NewState()

	// Load Lua standard library
	L.OpenLibs()
	// Register our functions
	L.SetGlobal("vpoke", L.NewFunction(PWvPoke))
	L.SetGlobal("plot", L.NewFunction(PWplot))
	L.SetGlobal("termprint", L.NewFunction(PWtermPrint))
	L.SetGlobal("time", L.NewFunction(PWtime))
	L.SetGlobal("pchar", L.NewFunction(PWpchar))
	L.SetGlobal("vertline", L.NewFunction(PWvertline))
	L.SetGlobal("horizline", L.NewFunction(PWhorizline))
	L.SetGlobal("clear", L.NewFunction(PWclear))
	L.SetGlobal("print", L.NewFunction(PWprint))

	// Add event emitter
    L.SetGlobal("events", luar.New(L, events))

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	wintitle := "Pinwheel " + version

	window, err := sdl.CreateWindow(wintitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(res * scale), int32(res * scale), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	//cursortexture, _ := sdl.LoadBMP("cursor.bmp")
	//cursor := sdl.CreateColorCursor(cursortexture, 0, 0)

	screen, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, int32(res), int32(res))
	if err != nil {
		panic(err)
	}
	defer screen.Destroy()

	if err := L.DoFile(program); err != nil {
		panic(err)
	}
	defer L.Close()

	// "CPU Cycle," our main loop
	running := true
	start = time.Now()
	c := palette[0]
	renderer.SetDrawColor(c[0], c[1], c[2], sdl.ALPHA_OPAQUE)
	//sdl.SetCursor(cursor)
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
				case *sdl.QuitEvent:
					running = false
					break
				case *sdl.MouseButtonEvent:
					if e.Type == sdl.MOUSEBUTTONDOWN { emitter.Emit("mouseClick", e.X, e.Y) }
				case *sdl.MouseMotionEvent:
					if _, _, state := sdl.GetMouseState(); state & sdl.Button(sdl.BUTTON_LEFT) == 1 {
						emitter.Emit("mouseDrag", e.X, e.Y)
					}
					emitter.Emit("mouseMove", e.X, e.Y)
			}
		}

		// Call the Spin function from Lua
		if err := L.CallByParam(lua.P{
			Fn: L.GetGlobal("Spin"),
			NRet: 0,
			Protect: true,
		}); err != nil {
			panic(err)
		}

		// Update the screen with our pixel backbuffer
		screen.Update(nil, pixelbuf, res * 4)
		renderer.Copy(screen, nil, nil)

		// Flush screen
		renderer.Present()

		time.Sleep(16 * time.Millisecond)
	}
}

func randf(r int) int {
	return rand.Intn(r)
}

// vpoke(addr, val)
func PWvPoke(L *lua.LState) int {
	addr := L.ToInt(1)
	val := L.ToInt(2)

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
	offset := ( res * 4 * y ) + x * 4;
	pixelbuf[offset + 0] = byte(b)
	pixelbuf[offset + 1] = byte(g)
	pixelbuf[offset + 2] = byte(r)
	pixelbuf[offset + 3] = sdl.ALPHA_OPAQUE;
}