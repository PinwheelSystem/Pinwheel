# Pinwheel
> 🍭 An awesome little fantasy computer designed to be simple.

Pinwheel is an all new fantasy computer, developed in Go and designed for ease of use and simplicity,
and highly customizable where applicable (and wanted..!).  
It is currently in the VERY early prototype-ish stage, and isn't meant for normal use yet.  
Stay tuned though, I'm spending my entire holiday to work on this little thing as much I can.

# Planned Features 
- Customizable Palette of 64 COLORS (by program and/or user)
- Easily scriptable 
- Simplicity 
- Joypad, keyboard and mouse input

# Documentation
## Functions
**Expected to change A LOT during development.**  

- `Spin()` - Called every CPU cycle  

- `termprint(text)` - Print text to the terminal
- `vpoke(address, value)` - Write `value` to VRAM `address`
- `plot(x, y, color)` - Place pixel at `x, y` with palette color number `color`
- `time()` -> `number` - Gets the amount of time since boot.
- `pchar(char, x, y)` - Place `char`acter at `x, y`