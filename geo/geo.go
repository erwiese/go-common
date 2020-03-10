package geo

import "math"

const (
	//degreesToRadian = math.Pi / 180
	rho = 180 / math.Pi

	// WGS-84
	ae = 6378137.000
	be = 6356752.3142
)

var (
	ef = (math.Pow(ae, 2) - math.Pow(be, 2)) / math.Pow(ae, 2)
	es = (math.Pow(ae, 2) - math.Pow(be, 2)) / math.Pow(be, 2)
)

// XYZ2ell converts cartesian coordinates to geographic coordinates in degree.
func XYZ2ell(x, y, z float64) (lat float64, lon float64, h float64) {
	// Special locations (geocenter, pole)
	if x == 0 && y == 0 && z == 0 {
		return 0, 0, 0
	}

	if x == 0 && y == 0 && z > 0 {
		return 90, 0, z - be
	}
	if x == 0 && y == 0 && z < 0 {
		return -90, 0, -z - be
	}

	// Auxiliary quantities
	p := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
	t := math.Atan2(z*ae, p*be)

	// Longitude and latitude
	lon = math.Atan2(y, x)
	lat = math.Atan2((z + es*be*math.Pow(math.Sin(t), 3)), (p - ef*ae*math.Pow(math.Cos(t), 3)))

	// Auxiliary quantity
	n := math.Pow(ae, 2) / math.Sqrt(math.Pow(ae, 2)*math.Pow(math.Cos(lat), 2)+math.Pow(be, 2)*math.Pow(math.Sin(lat), 2))

	// Height
	h = p/math.Cos(lat) - n

	lat *= rho
	lon *= rho
	return
}
