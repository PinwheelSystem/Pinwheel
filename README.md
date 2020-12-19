![](https://modeus.is-inside.me/WcZYvhEk.png)  
<div align="center">
	<p>
		<a href="https://www.codeshelter.co/">
			<img alt="Code Shelter" src="https://www.codeshelter.co/static/badges/badge-flat.svg">
		</a>
	</p>
	<h1>Pinwheel</h1>
</div>  

> 🍭 An awesome little fantasy computer designed to be simple.

Pinwheel is an all new fantasy computer, developed in Go and designed for ease of use and simplicity,
and highly customizable where applicable (and wanted..!).  
It is currently in the VERY early alpha-ish stage, and isn't meant for normal, ordinary use yet.  
Stay tuned though, I'm spending my entire holiday to work on this little thing as much I can.

# Planned Features 
- Customizable Palette of 64 COLORS (by program and/or user)
- Easily scriptable 
- Simplicity 
- Joypad, keyboard and mouse input

# Install
## Precompiled Binaries 
soon:tm:

## Compiling
### Prerequisites 
1. [Go](https://go.dev/)  
2. Git (optional)
3. SDL2 (follow [these steps](https://github.com/veandco/go-sdl2#requirements))  
  - Currently only SDL2 itself is needed

### Steps
1. Download the repository
  - From GitHub: https://github.com/PinwheelSystem/Pinwheel/archive/master.zip
  - Or with Git: `git clone https://github.com/PinwheelSystem/Pinwheel`
2. Open a terminal, change to the downloaded repo and run: `go build -o Pinwheel .`
  - Windows users will have to add the `.exe` to `Pinwheel`
3. Assuming there are no errors, you should be able to run `./Pinwheel`.  
The first argument is the Lua source to run, for example to run a program: `./Pinwheel game.lua`

# Documentation
## Functions
**Expected to change A LOT during development.**  

- `Spin()` - Called every CPU cycle  

- `termprint(text)` - Print text to the terminal
- `vpoke(address, value)` - Write `value` to VRAM `address`
- `plot(x, y, color)` - Place pixel at `x, y` with palette color number `color`
- `time()` -> `number` - Gets the amount of time since boot
- `pchar(char, x, y)` - Place `char`acter at `x, y`
- `vertline(x, color)` - Draw a vertical line on the x axis
- `horizline(y, color)` - Draw a horizontal line on the y axis
- `clear()` - Clears the screen
- `print(text, x, y, color)` - Print text starting at `x, y` in the color `color`

## Events
Pinwheel, unlike most FCs, is event driven.  
To handle any event, you use `events`:  

```lua
events:On("eventName", function(arg1, arg2)
	-- Inside here is our callback function
	-- `arg1` and `arg2` are our callback arguments
end)
```

Events should be outside the `Spin` function.
#### mouseClick
Callback arguments:
  - `x` X coordinate of the location of the click
  - `y` Y coordinate of the location of the click
Emitted when a user clicks the mouse.

> ⚠ The `x` and `y` values must be divided by 4 to apply to Pinwheel's internal resolution!

#### mouseDrag
Callback arguments:
  - `x` X coordinate of the current cursor location
  - `y` Y coordinate of the current cursor location 
Emitted when a user holds a mouse button and moves the mouse.

> ⚠ The `x` and `y` values must be divided by 4 to apply to Pinwheel's internal resolution!

#### mouseMove
Callback arguments:
  - `x` X coordinate of the location of the cursor
  - `y` Y coordinate of the location of the cursor
Emitted when a user moves the mouse.

> ⚠ The `x` and `y` values must be divided by 4 to apply to Pinwheel's internal resolution!

# Thanks
- pixelbath#8145 from the [Fantasy Console Discord](https://discord.gg/BYbjDEP) for a Lua snippet of loading a Bitmap font.

# License
[BSD 3-Clause](LICENSE)