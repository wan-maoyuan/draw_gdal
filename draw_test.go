package gdaldraw

import "testing"

func TestConvert4326To3857(t *testing.T) {
	x, y := convert4326To3857(90, -180)
	t.Logf("lat: 90 lon: -180 x: %f y: %f", x, y)

	x, y = convert4326To3857(90, 180)
	t.Logf("lat: 90 lon: 180 x: %f y: %f", x, y)

	x, y = convert4326To3857(-90, -180)
	t.Logf("lat: -90 lon: -180 x: %f y: %f", x, y)

	x, y = convert4326To3857(-90, 180)
	t.Logf("lat: -90 lon: 180 x: %f y: %f", x, y)
}
