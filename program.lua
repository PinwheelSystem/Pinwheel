function Spin()
	--[[for i = 1000, 0, -1 do
		x = math.random(127)
		y = math.random(127)

		offset = ( 128 * 4 * y ) + x * 4
		vpoke(offset + 0, math.random(255))
		vpoke(offset + 1, math.random(255))
		vpoke(offset + 2, math.random(255))
		vpoke(offset + 3, 255)
	end]]--
	-- Ignore above, pixel noise for initial testing

	for i = 128, 0, -1 do
		plot(i, 64, i % 64)
	end
end