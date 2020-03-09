package geo

import "math"

const (
	degreesToRadian = math.Pi / 180
	rho             = 180 / math.Pi

	// WGS-84
	ae = 6378137.000
	be = 6356752.3142
)

var (
	ef = (math.Pow(ae, 2) - math.Pow(be, 2)) / math.Pow(ae, 2)
	es = (math.Pow(ae, 2) - math.Pow(be, 2)) / math.Pow(be, 2)
)

// XYZ2ell converts cartesian coordinates to geographic coordinates in degree.
func XYZ2ell(x, y, z float64) (phi float64, lam float64, h float64) {
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

	return

}

/*
sub xyzEll {


# Auxiliary quantities
  my $p = sqrt(x**2+y**2);
  my $t = atan2(z*ae,$p*be);

# Latitude and longitude
  my $lam = atan2(y,x);
  my $phi = atan2((z+$es*be*sin($t)**3),($p-$ef*ae*cos($t)**3));

# Auxiliary quantity
  my $n = ae**2/sqrt(ae**2*cos($phi)**2+be**2*sin($phi)**2);

# Height
  my $h = $p/cos($phi)-$n;

# Return
  return $typ eq 'r' ? ($phi,$lam,$h) : ($phi*$rho,$lam*$rho,$h);
} */
