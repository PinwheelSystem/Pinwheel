function Spin()
	for i = 1000, 0, -1 do
		x = math.random(0, 239)
		y = math.random(0, 159)

		plot(x, y, math.random(0, 63))
	end
end
