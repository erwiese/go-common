package geo

import (
	"math"
	"testing"
)

func TestXYZ2ell(t *testing.T) {
	type xyz struct {
		x float64
		y float64
		z float64
	}
	tests := []struct {
		name    string
		xyz     xyz
		wantLat float64
		wantLon float64
		wantH   float64
	}{
		{name: "LEIJ", xyz: xyz{3898736.5150, 855345.1250, 4958372.3700}, wantLat: 51.35398, wantLon: 12.37410, wantH: 178.389},
		{name: "WTZR", xyz: xyz{4075580.38441, 931853.97899, 4801568.24545}, wantLat: 49.14420, wantLon: 12.87891, wantH: 666.030},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLat, gotLon, gotH := XYZ2ell(tt.xyz.x, tt.xyz.y, tt.xyz.z)
			if math.Round(gotLat*100000)/100000 != tt.wantLat {
				t.Errorf("XYZ2ell() gotLat = %v, want %v", gotLat, tt.wantLat)
			}
			if math.Round(gotLon*100000)/100000 != tt.wantLon {
				t.Errorf("XYZ2ell() gotLon = %v, want %v", gotLon, tt.wantLon)
			}
			if math.Round(gotH*1000)/1000 != tt.wantH {
				t.Errorf("XYZ2ell() gotH = %v, want %v", gotH, tt.wantH)
			}
		})
	}
}

func TestEll2XYZ(t *testing.T) {
	type ell struct {
		n  float64
		e  float64
		up float64
	}
	tests := []struct {
		name  string
		ell   ell
		wantX float64
		wantY float64
		wantZ float64
	}{

		{name: "LEIJ", ell: ell{51.3539791666985, 12.374095866237646, 178.3889}, wantX: 3898736.515, wantY: 855345.125, wantZ: 4958372.370},
		{name: "WTZR", ell: ell{49.144199946431826, 12.878912704434896, 666.0298}, wantX: 4075580.384, wantY: 931853.979, wantZ: 4801568.245},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY, gotZ := Ell2XYZ(tt.ell.n, tt.ell.e, tt.ell.up)
			if math.Round(gotX*1000)/1000 != tt.wantX {
				t.Errorf("Ell2XYZ() gotX = %v, want %v", gotX, tt.wantX)
			}
			if math.Round(gotY*1000)/1000 != tt.wantY {
				t.Errorf("Ell2XYZ() gotY = %v, want %v", gotY, tt.wantY)
			}
			if math.Round(gotZ*1000)/1000 != tt.wantZ {
				t.Errorf("Ell2XYZ() gotZ = %v, want %v", gotZ, tt.wantZ)
			}
		})
	}
}
