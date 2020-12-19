events:On("mouseMove", function(x, y)
	clear()

	horizline(x / 4, 12) -- divided by window scale 
	vertline(y / 4, 12)

	print(string.format("(%03i, %03i)", x / 4, y / 4), 1, 1, 22)
end)

function Spin()end