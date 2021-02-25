package main

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/yuin/gopher-lua"
	"github.com/PinwheelSystem/bitmap"
	_ "github.com/PinwheelSystem/PaletteNom"
	"github.com/chuckpreslar/emission"
	"layeh.com/gopher-luar"
)

type Pinwheel struct {
	version string
	screen *PinwheelScreen
	vram *[]uint8 // Why isn't this in the PinwheelScreen struct?
				  // To troll you :)
	ram *[]uint8
	l *lua.LState
	program string
	events *emission.Emitter
}

type PinwheelScreen struct {
	Width int
	Height int
}

const VERSION string = "0.1.0"
var SCREEN PinwheelScreen = PinwheelScreen{240, 160}

func Init(data map[string]interface{}) Pinwheel {
	//palettelib := palettenom.New()
	//colors, _ := palettelib.Load(data["palettedir"].(string) + data["palettefile"].(string))

	bm := bitmap.New()
	font = bm.Load(data["fontfile"].(string))
	ram_ := make([]uint8, 65536)
	vram_ := make([]uint8, SCREEN.Width * SCREEN.Height * 4) // * 4 because the format of a SDL pixel array is A, R, G, B (with our used pixel format),
															 // so we have to account for that
	return Pinwheel{
		version: VERSION,
		screen: &SCREEN,
		vram: &vram_,
		ram: &ram_, // 64K RAM (64 kibibytes)
		l: lua.NewState(),
		program: data["program"].(string),
		events: emission.NewEmitter(),
	}
}

func (p *Pinwheel) Start() {
	// Load Lua standard library
	p.l.OpenLibs()
	// Register our functions
	p.l.SetGlobal("vpoke", p.l.NewFunction(PWvPoke))
	p.l.SetGlobal("plot", p.l.NewFunction(PWplot))
	p.l.SetGlobal("termprint", p.l.NewFunction(PWtermPrint))
	p.l.SetGlobal("time", p.l.NewFunction(PWtime))
	p.l.SetGlobal("pchar", p.l.NewFunction(PWpchar))
	p.l.SetGlobal("vertline", p.l.NewFunction(PWvertline))
	p.l.SetGlobal("horizline", p.l.NewFunction(PWhorizline))
	p.l.SetGlobal("clear", p.l.NewFunction(PWclear))
	p.l.SetGlobal("print", p.l.NewFunction(PWprint))

	// Add event emitter
    p.l.SetGlobal("events", luar.New(p.l, p.events))

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	wintitle := "Pinwheel " + version

	window, err := sdl.CreateWindow(wintitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(p.screen.Width * 4), int32(p.screen.Height * 4), sdl.WINDOW_SHOWN)
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

	screen, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, int32(p.screen.Width), int32(p.screen.Height))
	if err != nil {
		panic(err)
	}
	defer screen.Destroy()

	if err := p.l.DoFile(p.program); err != nil {
		panic(err)
	}
	defer p.l.Close()

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
					if e.Type == sdl.MOUSEBUTTONDOWN { p.events.Emit("mouseClick", e.X, e.Y) }
				case *sdl.MouseMotionEvent:
					if _, _, state := sdl.GetMouseState(); state & sdl.Button(sdl.BUTTON_LEFT) == 1 {
						p.events.Emit("mouseDrag", e.X, e.Y)
					}
					p.events.Emit("mouseMove", e.X, e.Y)
			}
		}

		// Call the Spin function from Lua
		if err := p.l.CallByParam(lua.P{
			Fn: p.l.GetGlobal("Spin"),
			NRet: 0,
			Protect: true,
		}); err != nil {
			panic(err)
		}

		// Update the screen with our pixel backbuffert
		screen.Update(nil, []byte(*p.vram), p.screen.Width * 4)
		renderer.Copy(screen, nil, nil)

		// Flush screen
		renderer.Present()

		time.Sleep(16 * time.Millisecond)
	}
}