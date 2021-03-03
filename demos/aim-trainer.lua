xx = math.random(9, 229)
yy = math.random(9, 149)
lw = pchar("O", xx, yy, 12)

events:On("mouseClick", function(x, y)
	if xx < x / 4 and xx + lw > x / 4 and yy < y / 4 and yy + 8 > y / 4 then
		xx = math.random(9, 229)
		yy = math.random(9, 149)
		lw = pchar("O", xx, yy, 12)
	else
		-- do nothing
	end
end)
events:On("mouseMove", function(x, y)
	clear()
	pchar("O", xx, yy, 12)

	print(string.format("(%03i, %03i)", x / 4, y / 4), 1, 1, 22)
end)

function Spin()end
