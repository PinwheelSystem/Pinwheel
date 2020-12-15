package main

import (
	"math/rand"
	"fmt"
	"time"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/yuin/gopher-lua"
	"github.com/muesli/gamut"
)

const version = "v0.0.1"
const res int = 128 // Resolution of the *screen* ("internal") . Might change later in development. (res x res, 128 x 128)
const scale = 4 // Resolution scale (contributes to the size of the *window*)
var palette [][]uint8 = make([][]uint8, 64) // Array of array of RGB values ([[R, G, B], [R, G, B], ...])
var pixelbuf []byte = make([]byte, res * res * 4) // Pixel backbuffer (basically our VRAM)
var start time.Time

func main() {
	colors, _ := gamut.Generate(64, gamut.PastelGenerator{})
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

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

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
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
				case *sdl.QuitEvent:
					running = false
					break
			}
		}

		renderer.SetDrawColor(0, 0, 0, 0)

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

func setpixel(x, y, r, g, b int) {
	offset := ( res * 4 * y ) + x * 4;
	pixelbuf[offset + 0] = byte(b)
	pixelbuf[offset + 1] = byte(g)
	pixelbuf[offset + 2] = byte(r)
	pixelbuf[offset + 3] = sdl.ALPHA_OPAQUE;
}