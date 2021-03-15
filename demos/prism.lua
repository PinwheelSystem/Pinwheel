pi8 = math.pi / 8
pi2 = math.pi * 2

t=0
function Spin()
	clear()

	for i=t%8,135,8 do
		line(i, 0, 0, 135 - i, 8)
		line(i, 135, 135, 135 - i, 6)
		t=t+0.01
	end

	for i = (t / 16) % pi8, pi2, pi8 do
		x = 68 + 32 * math.cos(i)
		y = 68 + 32 * math.cos(i)
		line(135, 0, x, y, 22)
		line(0, 135, x, y, 22)
	end

	line(0, 0, 135, 0, 8)
	line(0, 0, 0, 135, 8)
	line(135, 0, 135, 135, 6)
	line(0, 135, 135, 135, 6)
end
