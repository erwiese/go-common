// Convert cartesian coordinates xyz or lat/lon to a geohash.
//
// See: https://en.wikipedia.org/wiki/Geohash
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/StefanSchroeder/Golang-Ellipsoid/ellipsoid"
	"github.com/mmcloughlin/geohash"
)

func usage() {
	fmt.Fprintf(os.Stderr, `geohasher - convert coordinates to a geohash

Usage: 
    geohasher X Y Z
    geohasher lat lon
		   
Examples:
    $ Convert cartesioan coordinates XYZ to geohash
    geohasher 4075580.3453 931854.0052 4801568.2446
		
    $ Convert geographic coordinates lat/lon to geohash
    geohasher -4.4966 48.3805
	
Author: E. Wiesensarter, 2020
`)
	os.Exit(1)
}

func encodeXYZ(x, y, z float64) string {
	// Create Ellipsoid object with WGS84-ellipsoid, angle units are degrees, distance units are meter.
	elli := ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Meter, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)

	// Convert ECEF to Lat-Lon-Alt.
	lat, lon, _ := elli.ToLLA(x, y, z)
	//fmt.Printf("lat = %v lon = %v alt = %v\n", lat, lon, alt)

	return geohash.Encode(lat, lon)
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		usage()
	}

	if len(os.Args) == 3 {
		lat, err := strconv.ParseFloat(os.Args[1], 64)
		if err != nil {
			panic(err)
		}

		lon, err := strconv.ParseFloat(os.Args[2], 64)
		if err != nil {
			panic(err)
		}

		geoh := geohash.Encode(lat, lon)
		fmt.Printf("%s\n", geoh)
	} else if len(os.Args) == 4 {
		x, err := strconv.ParseFloat(os.Args[1], 64)
		if err != nil {
			panic(err)
		}

		y, err := strconv.ParseFloat(os.Args[2], 64)
		if err != nil {
			panic(err)
		}

		z, err := strconv.ParseFloat(os.Args[3], 64)
		if err != nil {
			panic(err)
		}

		geoh := encodeXYZ(x, y, z)
		fmt.Printf("%s\n", geoh)
	} else {
		usage()
	}

}
