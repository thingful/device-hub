// Copyright © 2017 thingful
//
// Geolocation minimal js object mimic
// https://dev.w3.org/geo/api/spec-source.html#position

package engine

import (
	"time"
)

// GeoLocation holds the geographic location
var GeoLocation position

type coordinates struct {
	Latitude         float64 // specified in decimal degrees
	Longitude        float64 // specified in decimal degrees
	Accuracy         float64 // specified in meters
	Altitude         float64 // specified in meters
	AltitudeAccuracy float64 // specified in meters
	Heading          float64 // specified in degrees
	Speed            float64 // specified in m/s
}

type position struct {
	Coords    coordinates
	Timestamp int64 // represents the time in ms when the position was acquired
}

func (p *position) SetGeoLocation(lng, lat, acc float64) {
	p.Timestamp = time.Now().UnixNano() / 1000000
	p.Coords.Latitude = lat
	p.Coords.Longitude = lng
	p.Coords.Accuracy = acc
}
