package geo

import "math"

const (
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
func XYZ2ell(x, y, z float64) (lat, lon, height float64) {
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

	p := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
	t := math.Atan2(z*ae, p*be)

	lon = math.Atan2(y, x)
	lat = math.Atan2((z + es*be*math.Pow(math.Sin(t), 3)), (p - ef*ae*math.Pow(math.Cos(t), 3)))

	n := math.Pow(ae, 2) / math.Sqrt(math.Pow(ae, 2)*math.Pow(math.Cos(lat), 2)+math.Pow(be, 2)*math.Pow(math.Sin(lat), 2))

	height = p/math.Cos(lat) - n

	// radians to degrees
	lat *= rho
	lon *= rho
	return
}

// Ell2XYZ converts geographic coordinates north, east, up in degrees to cartesian coordinates.
func Ell2XYZ(n, e, up float64) (x, y, z float64) {
	//$typ = "d" unless $typ;

	// Convert degrees to radians
	//(n,e) = map { $_/rho } (n,e) if $typ eq 'd';
	n, e = n/rho, e/rho

	sp := math.Sin(n)
	cp := math.Cos(n)
	N := ae / math.Sqrt(1-ef*sp*sp)

	x = (N + up) * cp * math.Cos(e)
	y = (N + up) * cp * math.Sin(e)
	z = (N - N*ef + up) * sp
	return
}

/*
// ========================================================================
// xyzRot
// ========================================================================
//
// xyzRot(x,y,z,ax,a,typ)
//
// ------------------------------------------------------------------------
// Arguments  Description                                           Default
// ------------------------------------------------------------------------
// x          Geocentric cartesian X coordinate
// Y          Geocentric cartesian Y coordinate
// Z          Geocentric cartesian Z coordinate
// ax         Rotation axis (x, y, z)
// a          Rotation angle
// typ        Radians ('r') or degree ('d')                         d
// ------------------------------------------------------------------------
//
// Purpose:   Rotate geocentric coordinates
//
// Comment:   ---
//
// Changes:   23-05-2010 Created
//
// ========================================================================
func xyzRot {
   my @old = splice(@_,0,3);
   my ($ax,$a,$typ) = @_;
   $typ = "d" unless $typ;
   $typ = lc(substr($typ,0,1));

// Conversion to radians
   $a = $a/rho if $typ eq 'd';

// Buffer sine/cosine
   my $sa = math.Sin($a);
   my $ca = math.Cos($a);

// Rotation matrix
   my @R = ();
   if      (lc($ax) eq "x") {
      @R = ( [1, 0,   0  ],
             [0, $ca,-$sa],
             [0, $sa, $ca] );
   } elsif (lc($ax) eq "y") {
      @R = ( [ $ca, 0, $sa ],
             [  0,  1,  0  ],
             [-$sa, 0, $ca ] );
   } elsif (lc($ax) eq "z") {
      @R = ( [$ca,-$sa, 0 ],
             [$sa, $ca, 0 ],
             [ 0,   0,  1 ] );
   }

// Rotate point
   my @new = undef;
   for (my $r=0;$r<=2;$r++) {
      for (my $c=0;$c<=2;$c++) {
         $new[$r] += $R[$r][$c]*$old[$c];
      }
   }

// Return
   return (@new);
}


// ========================================================================
// eccEll
// ========================================================================
//
// (dn,de,du) = eccEll(n,e,u,dx,dy,dz)
//
// ------------------------------------------------------------------------
// Arguments  Description                                           Default
// ------------------------------------------------------------------------
// n          Geographic North coordinate (rad)
// e          Geographic East  coordinate (rad)
// u          Geographic Up    coordinate (m)
// dx         Eccentricity in X (m)
// dy         Eccentricity in Y (m)
// dz         Eccentricity in Z (m)
// ------------------------------------------------------------------------
//
// Purpose:   Convert cartesian eccentricities (dx, dy, dz) to local
//             geographic eccentricities (in m)
//
// Comment:   ---
//
// Changes:   23-05-2010 Created
//
// ========================================================================
func eccEll {
   my (n,e,up, @ecc) = @_;

// Buffer sine/cosine
   my $sn = math.Sin(n);
   my $cn = math.Cos(n);
   my $se = math.Sin(e);
   my $ce = math.Cos(e);

// Rotation matrix
   my @R = ();
   $R[0][0] = -$sn*$ce;
   $R[0][1] = -$sn*$se;
   $R[0][2] =  $cn    ;
   $R[1][0] = -$se    ;
   $R[1][1] =  $ce    ;
   $R[1][2] =  0      ;
   $R[2][0] =  $cn*$ce;
   $R[2][1] =  $cn*$se;
   $R[2][2] =  $sn    ;

// Compute eccentricities
   my @neu = (0,0,0);
   for (my $ii=0; $ii<=2; $ii++) {
      for (my $jj=0; $jj<=2; $jj++) {
         $neu[$ii] += $R[$ii][$jj]*$ecc[$jj];
      }
   }

// Return
   return(@neu);
}


// ========================================================================
// degGms
// ========================================================================
//
// degGms(a,typ)
//
// ------------------------------------------------------------------------
// Arguments  Description                                           Default
// ------------------------------------------------------------------------
// a          Angle
// typ        Radians ('r') or degree ('d')                         d
// ------------------------------------------------------------------------
//
// Purpose:   Convert angle to deg/min/sec format
//
// Comment:   ---
//
// Changes:   23-10-2010 Created
//
// ========================================================================
func degGms {
   my ($a,$typ) = @_;
   $typ = "d" unless $typ;
   $typ = lc(substr($typ,0,1));

// Conversion to degree
   $a = $a*rho if $typ eq 'r';

// Convert to deg/min/sec format
   my $d = int($a);
   my $m = int(($a-$d)*60);
   my $s = ($a-$d-$m/60)*3600;
   $m = -$m if ($m<0);
   $s = -$s if ($s<0);

// Return
   return ($d,$m,$s);
}

// How to convert degrees,minutes,seconds to decimal degrees
// http://www.rapidtables.com/convert/number/degrees-minutes-seconds-to-degrees.htm
// dd = d + m/60 + s/3600 */
