events:On("mouseMove", function(x, y)
	clear()
	horizline(x / 4, 12) -- divided by window scale 
	vertline(y / 4, 12)
end)

function Spin()end